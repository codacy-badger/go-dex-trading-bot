generate:
	protoc -I lssdrpc/ protobuf/lssdrpc.proto --go_out=plugins=grpc:lssdrpc

build_linux:
    GOOS="linux" GOARCH="amd64" go build -o bot

build_windows:
    GOOS="windows" GOARCH="386" go build -o bot_win

.PHONY: generate build_linux build_windows