package main

import (
	"encoding/json"
	"os"
	"time"
)

type LNDConfig struct {
	Directory string `json:"lndDir"`
	CertPath  string `json:"certPath"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
}

type LSSDConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	TimeoutStr string `json:"timeout"`
	Timeout    time.Duration
}

type BotConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	LNCLIPath string `json:"lnCLIPath"`
}

type Config struct {
	Bot        BotConfig  `json:"botCfg"`
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

func readConfig() error {
	file, err := os.Open("cfg.json")
	if err != nil {
		logger.Fatalf("can't open config file: %v", err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Fatalf("can't decode config JSON: %v", err)
		return err
	}
	return nil
}
