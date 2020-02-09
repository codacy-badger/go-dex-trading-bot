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

type LNDConfig struct {
	TLSPaths map[string]string // map[LTC]-> /path/to/lnd_ltc.tls.cert
}

type Bot struct {
	OrderClient       lssdrpc.OrdersClient
	SwapClient        lssdrpc.SwapsClient
	CurrencyClient    lssdrpc.CurrenciesClient
	TradingPairClient lssdrpc.TradingPairsClient

	OrderConnection       *grpc.ClientConn
	SwapConnection        *grpc.ClientConn
	CurrencyConnection    *grpc.ClientConn
	TradingPairConnection *grpc.ClientConn

	LNDConfig

}

func (t *Bot) Init() {

}

/*
func (t *Bot) AddCurrency() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()


	ReadFile(t.)
		val tlsCertRawLTC = readFile("ltc.cert")
		val tlsCertRawXSN = readFile("xsn.cert")
		println("exchangeA - LTC tls")
		println(tlsCertRawLTC)
		println("exchangeA -XSN tls")
		println(tlsCertRawXSN)
		Exchange.apply(
			lssdConfig = LssdRpcConfig("localhost", 50051),
			lndLtcConfig =
		LndRpcConfig(host = "localhost", port = 10001, tlsCert = tlsCertRawLTC),
		lndXsnConfig =
			LndRpcConfig(host = "localhost", port = 10003, tlsCert = tlsCertRawXSN)
	)
	}

	private def readFile(name: String): String = {
		val url = getClass.getResource("/" + name)
		val source = scala.io.Source.fromURL(url)
		val data = source.getLines().mkString("\n")
		source.close()
		data
	}









	cr := &lssdrpc.AddCurrencyRequest{
		Currency: CurrencyLTC,
		LndChannel:
	}
	Currency   string `protobuf:"bytes,1,opt,name=currency,proto3" json:"currency,omitempty"`
	LndChannel string `protobuf:"bytes,2,opt,name=lndChannel,proto3" json:"lndChannel,omitempty"`
	// Types that are valid to be assigned to TlsCert:
	//	*AddCurrencyRequest_CertPath
	//	*AddCurrencyRequest_RawCert
	TlsCert              isAddCurrencyRequest_TlsCert `protobuf_oneof:"tlsCert"`
	t.CurrencyClient.AddCurrency(ctx, )
}

*/
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
		fmt.Printf("id: %d | pair: %s | side: %s | orderId: %s | price: %v | isMy: %v \n", i, o.PairId, o.Side, o.OrderId, o.Price, o.IsOwnOrder)
	}
	return orders, nil
}

func (t *Bot) PlaceOrder(tradingPair string, price int, side lssdrpc.OrderSide) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p:= strconv.Itoa(price)

	priceBigInt := &lssdrpc.BigInteger{Value:p}

	order := &lssdrpc.PlaceOrderRequest{
		PairId: tradingPair,
		Side: side,
		Funds: priceBigInt,
		Price: priceBigInt,
	}

	_, err := t.OrderClient.PlaceOrder(ctx, order)
	return err
}

func NewBot(
	o lssdrpc.OrdersClient,
	oc *grpc.ClientConn,
	s lssdrpc.SwapsClient,
	sc *grpc.ClientConn,
	c lssdrpc.CurrenciesClient,
	cc *grpc.ClientConn,
	t lssdrpc.TradingPairsClient,
	tc *grpc.ClientConn) *Bot {
	return &Bot{
		OrderClient:           o,
		SwapClient:            s,
		CurrencyClient:        c,
		TradingPairClient:     t,
		OrderConnection:       oc,
		SwapConnection:        sc,
		CurrencyConnection:    cc,
		TradingPairConnection: tc,
	}
}

func ReadFile (path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
