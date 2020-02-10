package trading

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/cwntr/go-dex-trading-bot/lssdrpc"
)

const (
	PairXSNLTC = "XSN_LTC"
	PairXSNBTC = "XSN_BTC"

	CurrencyLTC = "LTC"
	CurrencyXSN = "XSN"
	CurrencyBTC = "BTC"
)

type Bot struct {
	OrderClient       lssdrpc.OrdersClient
	SwapClient        lssdrpc.SwapsClient
	CurrencyClient    lssdrpc.CurrenciesClient
	TradingPairClient lssdrpc.TradingPairsClient

	OrderConnection       *grpc.ClientConn
	SwapConnection        *grpc.ClientConn
	CurrencyConnection    *grpc.ClientConn
	TradingPairConnection *grpc.ClientConn

	LNDConfig LNDConfig
}

func (t *Bot) Init() error {
	t.LNDConfig.Certs = make(map[string]string, 0)
	for currency, path := range t.LNDConfig.TLSPaths {
		b, err := ReadFile(path)
		if err != nil {
			return err
		}
		t.LNDConfig.Certs[currency] = string(b)
	}
	return nil
}

func (t *Bot) AddCurrencies() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for currency, cert := range t.LNDConfig.Certs {
		c := t.LNDConfig.Connections
		x := c[currency]
		cr := &lssdrpc.AddCurrencyRequest{
			Currency:   currency,
			LndChannel: x.Format(),
		}
		cr.TlsCert = &lssdrpc.AddCurrencyRequest_RawCert{RawCert: cert}
		_, err := t.CurrencyClient.AddCurrency(ctx, cr)
		if err != nil {
			return err
		}
	}
	return nil
}

//Let the Bot subscribe to swaps
func (t *Bot) SubscribeSwaps() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := t.SwapClient.SubscribeSwaps(ctx, &lssdrpc.SubscribeSwapsRequest{})
	return err
}

//Let the Bot subscribe to orders
func (t *Bot) SubscribeOrders() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := t.OrderClient.SubscribeOrders(ctx, &lssdrpc.SubscribeOrdersRequest{})
	return err
}

//Retrieves orders from the XSN DEX Orderbook by enabling the TradingPair and request the all orders
func (t *Bot) ListOrders(tradingPair string, myOrders bool) ([]lssdrpc.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := t.TradingPairClient.EnableTradingPair(ctx, &lssdrpc.EnableTradingPairRequest{PairId: PairXSNLTC})
	if err != nil {
		log.Fatalf("err setting tradingPair: %v", err)
		return nil, err
	}
	var orders []lssdrpc.Order
	var processed uint32 = 0
	var limit uint32 = 200
	for {
		res, err := t.OrderClient.ListOrders(ctx, &lssdrpc.ListOrdersRequest{
			PairId:           PairXSNLTC,
			IncludeOwnOrders: myOrders,
			Skip:             processed,
			Limit:            limit,
		})
		if err != nil {
			log.Fatalf("err listing orders: %v", err)
			return nil, err
		}
		if len(res.Orders) == 0 {
			break
		}
		for _, o := range res.Orders {
			if o != nil {
				orders = append(orders, *o)
			}
		}
		processed += uint32(len(res.Orders))
	}
	for i, o := range orders {
		fmt.Printf("id: %d | pair: %s | side: %s | orderId: %s | price: %v | funds: %v | isMy: %v \n", i, o.PairId, o.Side, o.OrderId, o.Price, o.Funds, o.IsOwnOrder)
	}
	return orders, nil
}

func (t *Bot) PlaceOrder(tradingPair string, price int, amount int, side lssdrpc.OrderSide) (*lssdrpc.PlaceOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	p := strconv.Itoa(price)
	priceBigInt := &lssdrpc.BigInteger{Value: p}
	a := strconv.Itoa(amount)
	amountBigInt := &lssdrpc.BigInteger{Value: a}
	order := &lssdrpc.PlaceOrderRequest{
		PairId: tradingPair,
		Side:   side,
		Funds:  amountBigInt,
		Price:  priceBigInt,
	}
	res, err := t.OrderClient.PlaceOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewBot(
	o lssdrpc.OrdersClient,
	oc *grpc.ClientConn,
	s lssdrpc.SwapsClient,
	sc *grpc.ClientConn,
	c lssdrpc.CurrenciesClient,
	cc *grpc.ClientConn,
	t lssdrpc.TradingPairsClient,
	tc *grpc.ClientConn,
	lndConfig LNDConfig) (*Bot, error) {
	if lndConfig.IsEmtpy() {
		return nil, fmt.Errorf("lndConfig is empty")
	}
	b := &Bot{
		OrderClient:           o,
		SwapClient:            s,
		CurrencyClient:        c,
		TradingPairClient:     t,
		OrderConnection:       oc,
		SwapConnection:        sc,
		CurrencyConnection:    cc,
		TradingPairConnection: tc,
		LNDConfig:             lndConfig,
	}
	err := b.Init()
	return b, err
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
