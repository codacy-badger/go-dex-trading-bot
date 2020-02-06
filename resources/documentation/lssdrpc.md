# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [lssd_app/lssdrpc.proto](#lssd_app/lssdrpc.proto)
    - [AddCurrencyRequest](#lssdrpc.AddCurrencyRequest)
    - [AddCurrencyResponse](#lssdrpc.AddCurrencyResponse)
    - [BigInteger](#lssdrpc.BigInteger)
    - [CancelOrderRequest](#lssdrpc.CancelOrderRequest)
    - [CancelOrderResponse](#lssdrpc.CancelOrderResponse)
    - [EnableTradingPairRequest](#lssdrpc.EnableTradingPairRequest)
    - [EnableTradingPairResponse](#lssdrpc.EnableTradingPairResponse)
    - [ListOrdersRequest](#lssdrpc.ListOrdersRequest)
    - [ListOrdersResponse](#lssdrpc.ListOrdersResponse)
    - [Order](#lssdrpc.Order)
    - [OrderUpdate](#lssdrpc.OrderUpdate)
    - [OrderbookFailure](#lssdrpc.OrderbookFailure)
    - [PlaceOrderFailure](#lssdrpc.PlaceOrderFailure)
    - [PlaceOrderRequest](#lssdrpc.PlaceOrderRequest)
    - [PlaceOrderResponse](#lssdrpc.PlaceOrderResponse)
    - [SubscribeOrdersRequest](#lssdrpc.SubscribeOrdersRequest)
    - [SubscribeSwapsRequest](#lssdrpc.SubscribeSwapsRequest)
    - [SwapFailure](#lssdrpc.SwapFailure)
    - [SwapResult](#lssdrpc.SwapResult)
    - [SwapSuccess](#lssdrpc.SwapSuccess)
  
    - [OrderSide](#lssdrpc.OrderSide)
    - [SwapSuccess.Role](#lssdrpc.SwapSuccess.Role)
  
  
    - [currencies](#lssdrpc.currencies)
    - [orders](#lssdrpc.orders)
    - [swaps](#lssdrpc.swaps)
    - [tradingPairs](#lssdrpc.tradingPairs)
  

- [Scalar Value Types](#scalar-value-types)



<a name="lssd_app/lssdrpc.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## lssd_app/lssdrpc.proto



<a name="lssdrpc.AddCurrencyRequest"></a>

### AddCurrencyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| currency | [string](#string) |  |  |
| lndChannel | [string](#string) |  |  |
| certPath | [string](#string) |  | Path to tlc certificate, must be used for grpc connection. |
| rawCert | [string](#string) |  |  |






<a name="lssdrpc.AddCurrencyResponse"></a>

### AddCurrencyResponse







<a name="lssdrpc.BigInteger"></a>

### BigInteger
A non-negative Big Integer represented as string, like &#34;100000000&#34;


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  |  |






<a name="lssdrpc.CancelOrderRequest"></a>

### CancelOrderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pairId | [string](#string) |  |  |
| orderId | [string](#string) |  |  |






<a name="lssdrpc.CancelOrderResponse"></a>

### CancelOrderResponse







<a name="lssdrpc.EnableTradingPairRequest"></a>

### EnableTradingPairRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pairId | [string](#string) |  | the trading pair to enable, like XSN_LTC |






<a name="lssdrpc.EnableTradingPairResponse"></a>

### EnableTradingPairResponse







<a name="lssdrpc.ListOrdersRequest"></a>

### ListOrdersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pairId | [string](#string) |  |  |
| includeOwnOrders | [bool](#bool) |  |  |






<a name="lssdrpc.ListOrdersResponse"></a>

### ListOrdersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| orders | [Order](#lssdrpc.Order) | repeated |  |






<a name="lssdrpc.Order"></a>

### Order



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pairId | [string](#string) |  |  |
| orderId | [string](#string) |  |  |
| price | [BigInteger](#lssdrpc.BigInteger) |  | The price of the order in satoshis. |
| funds | [BigInteger](#lssdrpc.BigInteger) |  | The funds of the order in satoshis. |
| createdAt | [uint64](#uint64) |  | The epoch time when this order was created. |
| side | [OrderSide](#lssdrpc.OrderSide) |  | Whether this order is a buy or sell |
| isOwnOrder | [bool](#bool) |  | Whether this order is a local own order or a remote peer order. |






<a name="lssdrpc.OrderUpdate"></a>

### OrderUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| order | [Order](#lssdrpc.Order) |  | An order that was added to the order book. |
| orderRemoval | [Order](#lssdrpc.Order) |  | An order that was removed from the order book. |






<a name="lssdrpc.OrderbookFailure"></a>

### OrderbookFailure



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pairId | [string](#string) |  |  |
| funds | [BigInteger](#lssdrpc.BigInteger) |  |  |
| failureReason | [string](#string) |  |  |
| requiredFee | [BigInteger](#lssdrpc.BigInteger) |  |  |






<a name="lssdrpc.PlaceOrderFailure"></a>

### PlaceOrderFailure



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| swapFailure | [SwapFailure](#lssdrpc.SwapFailure) |  |  |
| orderbookFalure | [OrderbookFailure](#lssdrpc.OrderbookFailure) |  |  |






<a name="lssdrpc.PlaceOrderRequest"></a>

### PlaceOrderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pairId | [string](#string) |  |  |
| side | [OrderSide](#lssdrpc.OrderSide) |  |  |
| funds | [BigInteger](#lssdrpc.BigInteger) |  |  |
| price | [BigInteger](#lssdrpc.BigInteger) |  | missing on market orders |






<a name="lssdrpc.PlaceOrderResponse"></a>

### PlaceOrderResponse
Outcome of place order, three possible situations
1. Order was placed
2. Order was placed and matched without going to orderbook
3. Place order or swap has failed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| swapSuccess | [SwapSuccess](#lssdrpc.SwapSuccess) |  |  |
| order | [Order](#lssdrpc.Order) |  |  |
| failure | [PlaceOrderFailure](#lssdrpc.PlaceOrderFailure) |  |  |






<a name="lssdrpc.SubscribeOrdersRequest"></a>

### SubscribeOrdersRequest







<a name="lssdrpc.SubscribeSwapsRequest"></a>

### SubscribeSwapsRequest







<a name="lssdrpc.SwapFailure"></a>

### SwapFailure



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| orderId | [string](#string) |  | Order id assigned by orderbook |
| pairId | [string](#string) |  | The trading pair that the swap is for. |
| funds | [BigInteger](#lssdrpc.BigInteger) |  | The order funds that was attempted to be swapped. |
| failureReason | [string](#string) |  | The reason why the swap failed. |






<a name="lssdrpc.SwapResult"></a>

### SwapResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [SwapSuccess](#lssdrpc.SwapSuccess) |  |  |
| failure | [SwapFailure](#lssdrpc.SwapFailure) |  |  |






<a name="lssdrpc.SwapSuccess"></a>

### SwapSuccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| orderId | [string](#string) |  | Order id assigned by orderbook |
| pairId | [string](#string) |  | The trading pair that the swap is for. |
| funds | [BigInteger](#lssdrpc.BigInteger) |  | The order funds that was swapped. |
| rHash | [string](#string) |  | The hex-encoded payment hash for the swap. |
| amountReceived | [BigInteger](#lssdrpc.BigInteger) |  | The amount received denominated in satoshis. |
| amountSent | [BigInteger](#lssdrpc.BigInteger) |  | The amount sent denominated in satoshis. |
| role | [SwapSuccess.Role](#lssdrpc.SwapSuccess.Role) |  | Our role in the swap, either MAKER or TAKER. |
| currencyReceived | [string](#string) |  | The ticker symbol of the currency received. |
| currencySent | [string](#string) |  | The ticker symbol of the currency sent. |
| rPreimage | [string](#string) |  | The hex-encoded preimage. |
| price | [BigInteger](#lssdrpc.BigInteger) |  | The price used for the swap. |





 


<a name="lssdrpc.OrderSide"></a>

### OrderSide


| Name | Number | Description |
| ---- | ------ | ----------- |
| buy | 0 |  |
| sell | 1 |  |



<a name="lssdrpc.SwapSuccess.Role"></a>

### SwapSuccess.Role


| Name | Number | Description |
| ---- | ------ | ----------- |
| TAKER | 0 |  |
| MAKER | 1 |  |


 

 


<a name="lssdrpc.currencies"></a>

### currencies
currencies

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddCurrency | [AddCurrencyRequest](#lssdrpc.AddCurrencyRequest) | [AddCurrencyResponse](#lssdrpc.AddCurrencyResponse) |  |


<a name="lssdrpc.orders"></a>

### orders
orders

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| PlaceOrder | [PlaceOrderRequest](#lssdrpc.PlaceOrderRequest) | [PlaceOrderResponse](#lssdrpc.PlaceOrderResponse) |  |
| CancelOrder | [CancelOrderRequest](#lssdrpc.CancelOrderRequest) | [CancelOrderResponse](#lssdrpc.CancelOrderResponse) |  |
| SubscribeOrders | [SubscribeOrdersRequest](#lssdrpc.SubscribeOrdersRequest) | [OrderUpdate](#lssdrpc.OrderUpdate) stream |  |
| ListOrders | [ListOrdersRequest](#lssdrpc.ListOrdersRequest) | [ListOrdersResponse](#lssdrpc.ListOrdersResponse) |  |


<a name="lssdrpc.swaps"></a>

### swaps
swaps

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SubscribeSwaps | [SubscribeSwapsRequest](#lssdrpc.SubscribeSwapsRequest) | [SwapResult](#lssdrpc.SwapResult) stream |  |


<a name="lssdrpc.tradingPairs"></a>

### tradingPairs
trading pairs

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| EnableTradingPair | [EnableTradingPairRequest](#lssdrpc.EnableTradingPairRequest) | [EnableTradingPairResponse](#lssdrpc.EnableTradingPairResponse) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

