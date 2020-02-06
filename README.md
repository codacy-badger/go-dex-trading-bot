# go-dex-trading-bot
A trading bot for Stakenet XSN Decentralized Exchange (DEX)


#protoc --version -> 3.6.1
#protoc --proto_path=lssdrpc --go_out=generated lssdrpc/lssdrpc.proto

documentation
tool: go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

protoc --doc_out=doc --doc_opt=html,index.html  lssdrpc/lssdrpc.proto
