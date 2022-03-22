CONFIG_PATH = ${HOME}/.wal
TLS_CONFIG_PATH=internal/server/security/authentication/tls/cloudflare/config
CASBIN_CONFIG_PATH=internal/server/security/authorization/acl/casbin


.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}


gencert:
	cfssl gencert -initca=true ${TLS_CONFIG_PATH}/ca-csr.json | cfssljson -bare ca
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=${TLS_CONFIG_PATH}/ca-config.json -profile=server ${TLS_CONFIG_PATH}/server-csr.json | cfssljson -bare server
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=${TLS_CONFIG_PATH}/ca-config.json -profile=client -cn="root" ${TLS_CONFIG_PATH}/client-csr.json | cfssljson -bare root-client
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=${TLS_CONFIG_PATH}/ca-config.json -profile=client -cn="nobody" ${TLS_CONFIG_PATH}/client-csr.json | cfssljson -bare nobody-client
	mv *.pem *.csr ${CONFIG_PATH}/

clean:
	rm -rdf $(CONFIG_PATH)/*

cpsecurity_config:
	cp $(CASBIN_CONFIG_PATH)/*  $(CONFIG_PATH)/

test: cpsecurity_config
	go test -race ./test/...

protocgen:
	protoc --go_out=:api/v1 --proto_path=:api/v1 api/v1/proto/*msg.proto
	protoc --go-grpc_out=require_unimplemented_servers=false:api/v1 --proto_path=:api/v1/proto  api/v1/proto/*svc.proto


all: clean cpsecurity_config gencert
	echo 'all done!'