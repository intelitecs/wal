package log_test

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	api "github.com/intelitecs/wal/api/v1/log"
	wal_log "github.com/intelitecs/wal/internal/log"
	"github.com/stretchr/testify/require"
)

func TestSegment(t *testing.T) {
	dir, _ := ioutil.TempDir("", "segment_test")
	defer os.RemoveAll(dir)
	want := &api.Record{Value: []byte("hello world")}
	c := wal_log.Config{}
	c.Segment.MaxStoreBytes = 1024
	c.Segment.MaxIndexBytes = wal_log.EntWidth * 3

	s, err := wal_log.NewSegment(dir, 16, c)
	require.NoError(t, err)
	require.Equal(t, uint64(16), s.NextOffset, s.NextOffset)
	require.False(t, s.IsMaxed())

	for i := uint64(0); i < 3; i++ {
		off, err := s.Append(want)
		require.NoError(t, err)
		require.Equal(t, 16+i, off)

		got, err := s.Read(off)
		require.NoError(t, err)
		require.Equal(t, want.Value, got.Value)
	}

	_, err = s.Append(want)
	require.Equal(t, io.EOF, err)

	require.True(t, s.IsMaxed())
	c.Segment.MaxStoreBytes = uint64(len(want.Value) * 3)
	c.Segment.MaxIndexBytes = 1024

	s, err = wal_log.NewSegment(dir, 16, c)
	require.NoError(t, err)
	require.True(t, s.IsMaxed())
	err = s.Remove()
	require.NoError(t, err)
	s, err = wal_log.NewSegment(dir, 16, c)
	require.NoError(t, err)
	require.False(t, s.IsMaxed())
}
