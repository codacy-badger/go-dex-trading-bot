# go-dex-trading-bot
A trading bot for Stakenet XSN Decentralized Exchange (DEX) written in golang.

## Create stub via protoc
Use the following link to install the prerequisites (https://grpc.io/docs/quickstart/go/):
1. install protoc compiler (3.6.1+) 
2. `protoc-gen-go` compiler plugin 

#### Generate a stub by using the lssdrpc.proto file
The latest lssdrpc.proto file can be found on: https://github.com/X9Developers/DexAPI/releases

Go to the project root, execute the following command to generate a go client 

`protoc -I. --go_out=plugin=grpc=:lssdrpc  protobuf/*.protoc`
`mv lssdrpc/protobuf/lssdrpc.pb.go lssdrpc/ && rm -R lssdrpc/protobuf`
which will output a lssdrpc.rb.go that has client and server connectors automatically generated.


