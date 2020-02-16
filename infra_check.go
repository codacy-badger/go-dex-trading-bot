package main

import (
	"strconv"

	"github.com/cwntr/go-dex-trading-bot/lncli"
	"github.com/cwntr/go-dex-trading-bot/trading"
)

func getMinimumSatoshis(currency string) int64 {
	switch currency {
	case trading.CurrencyXSN:
		return 60000
	case trading.CurrencyLTC:
		return 275000
	case trading.CurrencyBTC:
		return 20000
	}
	return 0
}

type Peer struct {
	RemoteKey string
	Address   string
}

func getHubPeers(currency string) []Peer {
	switch currency {
	case trading.CurrencyXSN:
		return []Peer{
			{"02a49dc96ebcfd889f2cc694e9135fc8a502f7df4aa42b9a6f88d57759ddee5385", "134.209.164.91:8384"},
			{"0207e0cc77fcb96a895e6935dee205e2b16a8df5ed49827f70977a268625237214", "134.209.164.91:7777"},
		}
	case trading.CurrencyLTC:
		return []Peer{
			{"032c6e03e7a316baa3fb64fb360ebd8520e90a21b5a3c4bfca7fd75689a1564ae3", "134.209.164.91:8002"},
			{"02b9068bdb2837afc380a166434dea1a3c798cb2b7cbc8da7c7b7ac7854be11664", "134.209.164.91:9002"},
		}
	case trading.CurrencyBTC:
		return []Peer{
			{"03f1b0a270f5e57292732c3b44b2582324b144df6b2981a3ba46b6a5bc435479bc", "134.209.164.91:8000"},
			{"023a0c53433904339b7bdbd725151aaed10922e2f9aca5ec979ed75a370615f24e", "134.209.164.91:9000"},
		}
	}
	return []Peer{}
}

func getChannelCapacities(currency string) (min int64, max int64) {
	switch currency {
	case trading.CurrencyXSN:
		min = 60000
		max = 10000000000
		return
	case trading.CurrencyLTC:
		min = 275000
		max = 100000000
		return
	case trading.CurrencyBTC:
		min = 20000
		max = 1600000
		return
	}
	return
}

func checkInfra(currency string) error {
	var lndCfg LNDConfig

	switch currency {
	case trading.CurrencyXSN:
		lndCfg = cfg.XSN
	case trading.CurrencyLTC:
		lndCfg = cfg.LTC
	case trading.CurrencyBTC:
		lndCfg = cfg.BTC
	}
	logger.Infoln("-------------------------------------")
	logger.Infof("Infra check: %s ", currency)
	balance, err := lncli.GetWalletBalance(cfg.Bot.LNCLIPath, lndCfg.Directory, lndCfg.Host, lndCfg.Port)
	if err != nil {
		logger.Errorf("unable to get wallet balance, err: %v", err)
		return err
	}
	if balance.ConfirmedSat > getMinimumSatoshis(currency) {
		logger.Infof("total balance %d (%.8f %s)", balance.TotalSat, balance.Total, currency)
		logger.Infof("confirmed balance %d > min-confirmed balance %d", balance.ConfirmedSat, getMinimumSatoshis(currency))
		logger.Infof("Infra check: %s - balance (1/3) OK ", currency)
	} else {
		logger.Errorf("Infra check: %s - balance (1/3) FAILED", currency)
		return err
	}

	lp, err := lncli.GetPeers(cfg.Bot.LNCLIPath, lndCfg.Directory, lndCfg.Host, lndCfg.Port)
	if err != nil {
		logger.Errorf("unable to list peers , err: %v", err)
		return err
	}

	peerCheck := false
	hubPeers := getHubPeers(currency)
	for _, p := range lp.Peers {
		for _, hp := range hubPeers {
			if p.PubKey == hp.RemoteKey && p.Address == hp.Address {
				logger.Infof("hub peer found: %s@%s", p.PubKey, p.Address)
				peerCheck = true
			}
		}
	}
	if peerCheck {
		logger.Infof("Infra check: %s - peers (2/3) OK ", currency)
	} else {
		logger.Errorf("Infra check: %s - peers (2/3) FAILED", currency)
		return err
	}

	lc, err := lncli.GetListChannels(cfg.Bot.LNCLIPath, lndCfg.Directory, lndCfg.Host, lndCfg.Port)
	if err != nil {
		logger.Errorf("unable to get listChannels, err: %v", err)
		return err
	}

	capacityMin, capacityMax := getChannelCapacities(currency)
	channelCheck := false
	for _, c := range lc.Channels {
		for _, hp := range hubPeers {
			if hp.RemoteKey == c.RemotePubkey {
				capacity, _ := strconv.ParseInt(c.Capacity, 10, 64)

				if capacity >= capacityMin && capacity <= capacityMax {
					localBalance, _ := strconv.ParseInt(c.LocalBalance, 10, 64)
					lbf := float64(localBalance)

					remoteBalance, _ := strconv.ParseInt(c.RemoteBalance, 10, 64)
					rbf := float64(remoteBalance)
					logger.Infof("channel capacity to hub is OK (chanID: %s) local_balance: %d (%.8f %s), remote_balance: %d (%.8f %s) ",
						c.ChanID,
						localBalance,
						lbf/1e8,
						currency,
						remoteBalance,
						rbf/1e8,
						currency,
					)
					channelCheck = true
				} else {
					logger.Errorf("channel capacity NOT OK -> your capacity %d must be between %d and %d ", capacity, capacityMin, capacityMax)
				}
			}
		}
	}
	if channelCheck {
		logger.Infof("Infra check: %s - channels (3/3) OK ", currency)
	} else {
		logger.Errorf("Infra check: %s - channels (3/3) FAILED", currency)
		return err
	}
	logger.Infof("Infra check: %s complete", currency)
	logger.Infoln("-------------------------------------")
	return nil
}
