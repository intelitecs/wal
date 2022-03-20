package server

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"testing"

	api "github.com/intelitecs/wal/api/v1/log"
	"github.com/intelitecs/wal/internal/log"

	"github.com/intelitecs/wal/internal/server/security/authentication/tls"
	"github.com/intelitecs/wal/internal/server/security/authorization/acl"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func TestServer(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T,
		rootClient api.LogClient,
		nobodyClient api.LogClient,
		config *Config,
	){
		"produce/consume a message to/from the log succeeds": testProduceConsume,
		"produce/consume stream succeeds":                    testProduceConsumeStream,
		"consume past boundary failed":                       testConsumePastBoundary,
		"authorization should fail":                          testUnAuthorized,
	} {
		t.Run(scenario, func(t *testing.T) {
			rootClient, nobodyClient, config, teardown := setupTest(t, nil)
			defer teardown()
			fn(t, rootClient, nobodyClient, config)
		})
	}
}

func setupTest(t *testing.T, fn func(*Config)) (
	rootClient api.LogClient,
	nobodyClient api.LogClient,
	cfg *Config,
	teradown func(),

) {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	newClient := func(certPath, keyPath string) (
		*grpc.ClientConn,
		api.LogClient,
		[]grpc.DialOption,
	) {
		tlsConfig, err := tls.SetupTLSConfig(tls.TLSConfig{
			CAFile:   tls.CAFile,
			CertFile: certPath,
			KeyFile:  keyPath,
			Server:   false,
		})
		require.NoError(t, err)
		tlsCreds := credentials.NewTLS(tlsConfig)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
		conn, err := grpc.Dial(l.Addr().String(), opts...)
		require.NoError(t, err)
		client := api.NewLogClient(conn)
		return conn, client, opts
	}

	var rootConn *grpc.ClientConn
	rootConn, rootClient, _ = newClient(tls.RootClientCertFile, tls.RootClientKeyFile)

	var nobodyConn *grpc.ClientConn
	nobodyConn, nobodyClient, _ = newClient(tls.NobodyClientCertFile, tls.NobodyClientKeyFile)

	serverTLSConfig, err := tls.SetupTLSConfig(tls.TLSConfig{
		CertFile:      tls.ServerCertFile,
		KeyFile:       tls.ServerKeyFile,
		CAFile:        tls.CAFile,
		ServerAddress: l.Addr().String(),
		Server:        true,
	})
	require.NoError(t, err)
	serverCredentials := credentials.NewTLS(serverTLSConfig)
	dir, err := ioutil.TempDir("", "server-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	clog, err := log.NewLog(dir, log.Config{})
	require.NoError(t, err)
	authorizer := acl.New(tls.ACLModelFile, tls.ACLPolicyFile)

	cfg = &Config{
		CommitLog:  clog,
		Authorizer: *authorizer,
	}
	if fn != nil {
		fn(cfg)
	}

	server, err := NewGRPCServer(cfg, grpc.Creds(serverCredentials))
	require.NoError(t, err)

	go func() {
		server.Serve(l)
	}()

	return rootClient, nobodyClient, cfg, func() {
		server.Stop()
		rootConn.Close()
		nobodyConn.Close()
		l.Close()
	}
}

func testProduceConsume(t *testing.T,
	client,
	_ api.LogClient,
	config *Config) {
	ctx := context.Background()
	want := &api.Record{
		Value: []byte("hello world"),
	}
	produce, err := client.Produce(ctx, &api.ProduceRequest{
		Record: want,
	},
	)

	require.NoError(t, err)
	consume, err := client.Consume(ctx, &api.ConsumeRequest{
		Offset: produce.Offset,
	})
	require.NoError(t, err)
	require.Equal(t, want.Value, consume.Record.Value)
	require.Equal(t, want.Offset, consume.Record.Offset)

}

func testProduceConsumeStream(
	t *testing.T,
	client,
	_ api.LogClient,
	config *Config,
) {
	ctx := context.Background()
	records := []*api.Record{{
		Value:  []byte("First message"),
		Offset: 0,
	}, {
		Value:  []byte("Second message"),
		Offset: 1,
	}}
	{
		stream, err := client.ProduceStream(ctx)
		require.NoError(t, err)
		for offset, record := range records {
			err = stream.Send(&api.ProduceRequest{
				Record: record,
			})
			require.NoError(t, err)
			res, err := stream.Recv()
			require.NoError(t, err)
			if res.Offset != uint64(offset) {
				t.Fatalf("got offset: %d, want: %d", res.Offset, offset)
			}
		}
	}
	{
		stream, err := client.ConsumeStream(ctx, &api.ConsumeRequest{
			Offset: 0,
		},
		)
		require.NoError(t, err)
		for i, record := range records {
			res, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, res.Record, &api.Record{
				Value:  record.Value,
				Offset: uint64(i),
			})
		}
	}
}

func testConsumePastBoundary(t *testing.T,
	client,
	_ api.LogClient,
	config *Config) {
	ctx := context.Background()
	want := &api.Record{
		Value: []byte("hello world"),
	}
	produce, err := client.Produce(ctx, &api.ProduceRequest{
		Record: want,
	})
	require.NoError(t, err)
	consume, err := client.Consume(ctx, &api.ConsumeRequest{
		Offset: produce.Offset + 1,
	})
	if consume != nil {
		t.Fatalf("consume not nil")
	}
	got := status.Code(err)
	wanted := status.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())

	if got != wanted {
		t.Fatalf("got err: %v want: %v", got, wanted)
	}
}

func testUnAuthorized(
	t *testing.T,
	_,
	client api.LogClient,
	config *Config,
) {
	ctx := context.Background()
	produce, err := client.Produce(ctx, &api.ProduceRequest{
		Record: &api.Record{
			Value: []byte("hello world"),
		},
	})
	if produce != nil {
		t.Fatalf("produce response should be nil")
	}
	gotCode, wantCode := status.Code(err), codes.PermissionDenied
	if gotCode != wantCode {
		t.Fatalf("got code: %d want: %d", gotCode, wantCode)
	}

	consume, err := client.Consume(ctx, &api.ConsumeRequest{
		Offset: 0,
	})
	if consume != nil {
		t.Fatalf("consume response should be nil")
	}
	gotCode, wantCode = status.Code(err), codes.PermissionDenied
	if gotCode != wantCode {
		t.Fatalf("got code: %d, want: %d", gotCode, wantCode)
	}
}
