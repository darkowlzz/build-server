generate:
	protoc -I pkg/build/ pkg/build/build.proto --go_out=plugins=grpc:pkg/build

cli:
	go build -v -o out/build ./cmd/client/

server:
	go build -v -o out/build-server ./cmd/server
