package main

import (
	"encoding/json"
	"net/http"

	"github.com/cwntr/go-dex-trading-bot/lncli"
	"github.com/cwntr/go-dex-trading-bot/lssdrpc"
	"github.com/cwntr/go-dex-trading-bot/trading"
)

func OrderbookXSNLTC(w http.ResponseWriter, r *http.Request) {
	orders, err := bot.ListOrders(trading.PairXSNLTC, true, true)
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

func OrderbookXSNBTC(w http.ResponseWriter, r *http.Request) {
	orders, err := bot.ListOrders(trading.PairXSNBTC, true, true)
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
	b, err := lncli.GetWalletBalance(cfg.Bot.LNCLIPath, cfg.XSN.Directory, cfg.XSN.Host, cfg.XSN.Port)
	if err != nil {
		logger.Errorf("err while GetWalletBalance %v", err)
		return
	}
	js, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func LTCBalanceFunc(w http.ResponseWriter, r *http.Request) {
	b, err := lncli.GetWalletBalance(cfg.Bot.LNCLIPath, cfg.LTC.Directory, cfg.LTC.Host, cfg.LTC.Port)
	if err != nil {
		logger.Errorf("err while GetWalletBalance %v", err)
		return
	}
	js, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func BTCBalanceFunc(w http.ResponseWriter, r *http.Request) {
	b, err := lncli.GetWalletBalance(cfg.Bot.LNCLIPath, cfg.BTC.Directory, cfg.BTC.Host, cfg.BTC.Port)
	if err != nil {
		logger.Errorf("err while GetWalletBalance %v", err)
		return
	}
	js, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func OrdersFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "only POST method allowed", http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var orders []Order
	err := decoder.Decode(&orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, order := range orders {
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
			res, err := bot.PlaceOrder(order.TradingPair, price, order.FixedFunding, side)
			if err != nil {
				logger.Errorf("err while placing an order %v", err)
			} else {
				logger.Infof("Added order, outcome: %v", res.Outcome)
			}
		}
	}
}

func OrdersCancelFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "only POST method allowed", http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var orderCancel OrderCancel
	err := decoder.Decode(&orderCancel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, cancel := range orderCancel.TradingPair {
		for _, id := range cancel.OrderIDs {
			//Cancel the order @TODO times out for some reason, but still works -> quick fix run it in parallel
			go bot.CancelOrder(cancel.TradingPair, id)
		}
		if cancel.DeleteAll {
			orders, err := bot.ListOrders(cancel.TradingPair, true, false)
			if err != nil {
				logger.Errorf("err while fetching all orders from orderbook, err %v", err)
				return
			}
			for _, o := range orders {
				if o.IsOwnOrder {
					//Cancel the order @TODO times out for some reason, but still works -> quick fix run it in parallel
					go bot.CancelOrder(cancel.TradingPair, o.OrderId)
					//if err != nil {
					//		logger.Errorf("err while canceling order: %s, err %v", o.OrderId, err)
					//	} else {
					//		logger.Infof("canceled orderId: %s", o.OrderId)
					//	}
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte{})
}

func makeRange(min int, max int, step int) []int {
	var rangeList []int
	for i := min; i <= max; i += step {
		rangeList = append(rangeList, i)
	}
	return rangeList
}
