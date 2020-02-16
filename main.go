package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cwntr/go-dex-trading-bot/trading"
)

var (
	cfg    Config
	logger *logrus.Entry
	bot    *trading.Bot
)

func main() {
	logger = logrus.WithFields(logrus.Fields{"context": "main"})

	//Read global config
	err := readConfig()
	if err != nil {
		logger.Errorf("error reading global config, err %v", err)
		return
	}
	cfg.LSSDConfig.Timeout, err = time.ParseDuration(cfg.LSSDConfig.TimeoutStr)
	if err != nil {
		logger.Errorf("unable parsing lssd timeout, err %v", err)
		return
	}
	logger.Infoln("global config loaded")

	//Initialize Clients
	tpClient, tpConn := createTradingPairClient()
	defer tpConn.Close()
	oClient, oConn := createOrdersClient()
	defer oConn.Close()
	cClient, cConn := createCurrencyClient()
	defer cConn.Close()
	sClient, sConn := createSwapClient()
	defer sConn.Close()
	logger.Infoln("clients initiated")

	//Initialize LNDConfig
	tradingCfg := trading.NewConfig()
	err = tradingCfg.Add(trading.CurrencyXSN, cfg.XSN.CertPath, cfg.XSN.Host, cfg.XSN.Port)
	if err != nil {
		logger.Errorf("error adding XSN to trading config, err %v", err)
		return
	}
	err = tradingCfg.Add(trading.CurrencyLTC, cfg.LTC.CertPath, cfg.LTC.Host, cfg.LTC.Port)
	if err != nil {
		logger.Errorf("error adding LTC to trading config, err %v", err)
		return
	}
	logger.Infoln("trading config loaded")

	//Initialize Bot
	bot, err = trading.NewBot(oClient, sClient, cClient, tpClient, tradingCfg, cfg.LSSDConfig.Timeout)
	if err != nil {
		logger.Errorf("error initializing trading bot, err %v", err)
		return
	}
	logger.Infoln("trading bot initialized")

	err = checkInfra(trading.CurrencyXSN)
	if err != nil {
		logger.Errorf("infra check, err %v", err)
		return
	}
	err = checkInfra(trading.CurrencyLTC)
	if err != nil {
		logger.Errorf("infra check, err %v", err)
		return
	}

	_ = simpleFlow(bot)

	//LSSD routes
	http.HandleFunc("/xsn_ltc/orderbook", OrderbookXSNLTC)
	http.HandleFunc("/xsn_btc/orderbook", OrderbookXSNBTC)
	http.HandleFunc("/orders", OrdersFunc)
	http.HandleFunc("/orders/cancel", OrdersCancelFunc)

	//LND routes
	http.HandleFunc("/xsn/balance", XSNBalanceFunc)
	http.HandleFunc("/ltc/balance", LTCBalanceFunc)
	http.HandleFunc("/btc/balance", BTCBalanceFunc)

	//Run HTTP server from bot config
	logger.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Bot.Host, cfg.Bot.Port), nil))
}

func simpleFlow(bot *trading.Bot) error {

	//Subscribe Swaps
	err := bot.SubscribeSwaps()
	if err != nil {
		logger.Errorf("error subscribing to swaps, err %v", err)
		return err
	}

	//Subscribe Orders
	err = bot.SubscribeOrders()
	if err != nil {
		logger.Errorf("error subscribing to orders, err %v", err)
		return err
	}

	//Add Currencies
	err = bot.AddCurrencies()
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return err
	}

	/*
		//List Orders
		_, err = bot.ListOrders(trading.PairXSNLTC, true, true)
		if err != nil {
			logger.Errorf("err while listing the orderbook %v", err)
			return err
		}
	*/
	return nil
}
