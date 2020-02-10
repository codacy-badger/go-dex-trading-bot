package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cwntr/go-dex-trading-bot/lssdrpc"
	"github.com/cwntr/go-dex-trading-bot/trading"
	"google.golang.org/grpc"
)

var (
	cfg Config
)

func main() {
	//Read global config
	err := readConfig()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Println("global config loaded")
	//Initialize Clients
	tpClient, tpConn := createTradingPairClient()
	defer tpConn.Close()

	oClient, oConn := createOrdersClient()
	defer oConn.Close()

	cClient, cConn := createCurrencyClient()
	defer cConn.Close()

	sClient, sConn := createSwapClient()
	defer sConn.Close()

	//Initialize LNDConfig
	tradingCfg := trading.NewConfig()
	err = tradingCfg.Add(trading.CurrencyXSN, cfg.XSN.CertPath, cfg.XSN.Host, cfg.XSN.Port)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	err = tradingCfg.Add(trading.CurrencyLTC, cfg.LTC.CertPath, cfg.LTC.Host, cfg.LTC.Port)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	fmt.Println("trading config loaded")

	//Initialize Bot
	bot, err := trading.NewBot(
		oClient,
		oConn,
		sClient,
		sConn,
		cClient,
		cConn,
		tpClient,
		tpConn,
		tradingCfg,
	)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	//Subscribe Swaps
	err = bot.SubscribeSwaps()
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	fmt.Println("subscribed to swaps")

	//Subscribe Orders
	err = bot.SubscribeOrders()
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	fmt.Println("subscribed to orders")

	//Add Currencies
	err = bot.AddCurrencies()
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}
	fmt.Println("added currencies")

	//Execute orders from config
	for _, order := range cfg.Orders {
		if order.PriceRangeStart == 0 || order.PriceRangeEnd == 0 || order.PriceRangeStepSize == 0 || order.FixedFunding == 0 {
			fmt.Printf("err: price setting missing")
			continue
		}

		if order.Side != "sell" && order.Side != "buy" {
			fmt.Printf("err: wrong side - must be either `sell` or `buy` ")
			continue
		}
		//place couple of sell orders
		for _, price := range makeRange(order.PriceRangeStart, order.PriceRangeEnd, order.PriceRangeStepSize) {
			funds := order.FixedFunding

			var side lssdrpc.OrderSide
			if order.Side == "sell" {
				side = lssdrpc.OrderSide_sell
			} else if order.Side == "buy" {
				side = lssdrpc.OrderSide_buy
			}

			res, err := bot.PlaceOrder(trading.PairXSNLTC, price, funds, side)
			if err != nil {
				fmt.Printf("err: %v \n", err)
				return
			}
			fmt.Println("added order")
			fmt.Printf("outcome: %v \n", res.Outcome)
		}
	}

	//List Orders
	_, _ = bot.ListOrders(trading.PairXSNLTC, true)
}

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

func readConfig() error {
	file, err := os.Open("cfg.json")
	if err != nil {
		log.Fatal("can't open config file: ", err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
		return err
	}
	return nil
}

func makeRange(min int, max int, step int) []int {
	var rangeList []int
	for i := min; i < max; i += step {
		rangeList = append(rangeList, i)
	}
	return rangeList
}
