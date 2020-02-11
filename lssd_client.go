package main

import (
	"fmt"
	"github.com/cwntr/go-dex-trading-bot/lssdrpc"
	"google.golang.org/grpc"
	"log"
)

func createSwapClient() (lssdrpc.SwapsClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.LSSDConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewSwapsClient(conn), conn
}

func createCurrencyClient() (lssdrpc.CurrenciesClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.LSSDConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewCurrenciesClient(conn), conn
}

func createTradingPairClient() (lssdrpc.TradingPairsClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.LSSDConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewTradingPairsClient(conn), conn
}

func createOrdersClient() (lssdrpc.OrdersClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.LSSDConfig.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewOrdersClient(conn), conn
}
