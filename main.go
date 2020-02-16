package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cwntr/go-stakenet/lnd"

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
	_ = simpleFlow(bot)

	http.HandleFunc("/lssd/orderbook", LSSDOrderbookFunc)
	http.HandleFunc("/xsn/balance", XSNBalanceFunc)
	//http.HandleFunc("/xsn/channels", XSNChannelsFunc)
	//	http.HandleFunc("/ltc/balance", LTCBalanceFunc)
	//http.HandleFunc("/ltc/channels", LTCChannelsFunc)
	logger.Fatal(http.ListenAndServe(":9999", nil))
}

func LSSDOrderbookFunc(w http.ResponseWriter, r *http.Request) {
	orders, err := bot.ListOrders(trading.PairXSNLTC, true)
	if err != nil {
		logger.Errorf("err while listing the orderbook %v", err)
		return
	}
	js, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func XSNBalanceFunc(w http.ResponseWriter, r *http.Request) {
	if cfg.Bot.LNCLIPath == "" || cfg.XSN.Directory == "" || cfg.XSN.Port == 0 {
		logger.Errorln("missing config for lncli or XSN lnd")
		return
	}
	lnPath := cfg.Bot.LNCLIPath
	lnOptions := []string{fmt.Sprintf("--lnddir=%s", cfg.XSN.Directory), "--no-macaroons", fmt.Sprintf("--rpcserver=%s:%d", cfg.XSN.Host, cfg.XSN.Port)}

	gi, err := lnd.GetGetInfo(lnPath, false, lnOptions...)
	if err != nil {
		logger.Errorf("err lnxsn get info: %v", err)
		return
	}
	js, err := json.Marshal(gi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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
		//Execute orders from config
		for _, order := range cfg.Orders {
			if order.PriceRangeStart == 0 || order.PriceRangeEnd == 0 || order.PriceRangeStepSize == 0 || order.FixedFunding == 0 {
				logger.Errorln("price range config: cannot have any '0' value")
				continue
			}

			if order.Side != "sell" && order.Side != "buy" {
				logger.Errorln("err: order wrong side - must be either `sell` or `buy`")
				continue
			}

			//Iterate over order price configs
			for _, price := range makeRange(order.PriceRangeStart, order.PriceRangeEnd, order.PriceRangeStepSize) {
				//resolve side
				var side lssdrpc.OrderSide
				if order.Side == "sell" {
					side = lssdrpc.OrderSide_sell
				} else if order.Side == "buy" {
					side = lssdrpc.OrderSide_buy
				}

				//Place the order
				res, err := bot.PlaceOrder(trading.PairXSNLTC, price, order.FixedFunding, side)
				if err != nil {
					logger.Errorf("err while placing an order %v", err)
				} else {
					logger.Infof("Added order, outcome: %v", res.Outcome)
				}
			}
		}

	*/
	//List Orders
	_, err = bot.ListOrders(trading.PairXSNLTC, true)
	if err != nil {
		logger.Errorf("err while listing the orderbook %v", err)
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
