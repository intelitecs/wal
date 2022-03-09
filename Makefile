
compile:
	protoc api/v1/proto/*.proto --go_out=:./api/v1  --go-grpc_out=./api/v1  
	 

test:
	go test -race ./...