#!/usr/bin/env bash

#XSN_LTC place SELL order with range (
curl -X POST \
  http://localhost:9999/orders \
  -H 'content-type: application/json' \
  -d '[{
    "side": "sell",
    "tradingPair": "XSN_LTC",
    "priceRangeStart": 95000,
    "priceRangeEnd": 95005,
    "priceRangeStepSize": 1,
    "fixedFunding": 100000
}]'

#XSN_BTC place BUY order with range
curl -X POST \
  http://localhost:9999/orders \
  -H 'content-type: application/json' \
  -d '[{
    "side": "buy",
    "tradingPair": "XSN_LTC",
    "priceRangeStart": 95000,
    "priceRangeEnd": 95005,
    "priceRangeStepSize": 1,
    "fixedFunding": 100000
}]'


#cancel orders by orderId for XSN_LTC
curl -X POST http://localhost:9999/orders/cancel \
  -d '{
	"cancelTradingPairs": [
		{
			"tradingPair": "XSN_LTC",
			"orderIds": ["e4e6fefb-5c44-4acb-a208-5e699f9d9a20"]
		}
	]
}'

#cancel all orders for XSN_LTC
curl -X POST \
  http://localhost:9999/orders/cancel \
  -d '{
	"cancelTradingPairs": [
		{
			"tradingPair": "XSN_LTC",
			"deleteAll": true
		}
	]
}'