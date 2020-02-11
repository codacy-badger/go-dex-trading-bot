package main

type LNDConfig struct {
	CertPath string `json:"certPath"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type LSSDConfig struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Timeout string `json:"timeout"`
}

type Config struct {
	LSSDConfig LSSDConfig `json:"lssdConfig"`
	XSN        LNDConfig  `json:"xsnLNDConfig"`
	LTC        LNDConfig  `json:"ltcLNDConfig"`
	BTC        LNDConfig  `json:"btcLNDConfig"`
	Orders     []Order    `json:"orders"`
}

type Order struct {
	Side               string `json:"side"`
	TradingPair        string `json:"tradingPair"`
	PriceRangeStart    int    `json:"priceRangeStart"`
	PriceRangeEnd      int    `json:"priceRangeEnd"`
	PriceRangeStepSize int    `json:"priceRangeStepSize"`
	FixedFunding       int    `json:"fixedFunding"`
}
