package main

import (
	"fmt"
	"log"

	"github.com/cwntr/go-dex-trading-bot/trading"

	"github.com/cwntr/go-dex-trading-bot/lssdrpc"
	"google.golang.org/grpc"
)

type Config struct {
	Host string
	Port int
}

func GetDefaultConfig() Config {
	return Config{"localhost", 50051}
}

func main() {
	//Initialize Services
	tpClient, tpConn := createTradingPairClient()
	defer tpConn.Close()

	oClient, oConn := createOrdersClient()
	defer oConn.Close()

	cClient, cConn := createCurrencyClient()
	defer cConn.Close()

	sClient, sConn := createSwapClient()
	defer sConn.Close()

	//Init Bot
	bot := trading.NewBot(
		oClient,
		oConn,
		sClient,
		sConn,
		cClient,
		cConn,
		tpClient,
		tpConn,
	)

	//Subscribe Swaps
	err := bot.SubscribeSwaps()
	if err != nil {
		fmt.Printf("err: %v \n", err)
	}

	//Subscribe Orders
	err = bot.SubscribeOrders()
	if err != nil {
		fmt.Printf("err: %v \n", err)
	}

	//List Orders
	_, _ = bot.ListOrders(trading.PairXSNBTC, true)

	//Place Order
	bot.PlaceOrder(trading.PairXSNLTC)

}

func createSwapClient() (lssdrpc.SwapsClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	cfg := GetDefaultConfig()
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewSwapsClient(conn), conn
}

func createCurrencyClient() (lssdrpc.CurrenciesClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	cfg := GetDefaultConfig()
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewCurrenciesClient(conn), conn
}

func createTradingPairClient() (lssdrpc.TradingPairsClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	cfg := GetDefaultConfig()
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewTradingPairsClient(conn), conn
}

func createOrdersClient() (lssdrpc.OrdersClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	cfg := GetDefaultConfig()
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "", cfg.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return lssdrpc.NewOrdersClient(conn), conn
}
