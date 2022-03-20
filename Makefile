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

$(CONFIG_PATH)/model.conf:
	cp $(CASBIN_CONFIG_PATH)/model.conf $(CONFIG_PATH)/model.conf

$(CONFIG_PATH)/policy.csv:
	cp $(CASBIN_CONFIG_PATH)/policy.csv  $(CONFIG_PATH)/policy.csv


.PHONY: test
test: $(CONFIG_PATH)/model.conf $(CONFIG_PATH)/policy.csv
	go test -race ./internal/...

.PHONY: genprotobuf

genmessages:
	protoc --go_out=:internal/adapters/framework/left/grpc --proto_path=internal/adapters/framework/left/grpc/proto internal/adapters/framework/left/grpc/proto/*msg.proto

genservices:
	protoc --go-grpc_out=require_unimplemented_servers=false:internal/adapters/framework/left/grpc  \
	--proto_path=internal/adapters/framework/left/grpc/proto internal/adapters/framework/left/grpc/proto/*svc.proto

protogen: genmessages genservices
	echo 'Generated messages and services'

all: genmsg gensvc
	echo 'Generate GRPC messages and services'

genmsg:
	protoc --go_out=:internal/adapters/framework/left/grpc --proto_path=:internal/adapters/framework/left/grpc/proto \
	internal/adapters/framework/left/grpc/proto/*msg.proto

gensvc:
	protoc --go-grpc_out=require_unimplemented_servers=false:internal/adapters/framework/left/grpc \
	--proto_path=:internal/adapters/framework/left/grpc/proto internal/adapters/framework/left/grpc/proto/*svc.proto

	
