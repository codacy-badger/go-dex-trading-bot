package common

type Config struct {
	LNCLIPath                string    `json:"lncliPath"`
	LNCLIOptions             []string  `json:"lncliOptions"`
	IsCoreWalletTestRequests bool      `json:"testRequests"` //core wallet
	IsLNWalletTestRequests   bool      `json:"testRequestsLightning"`
	Environment              string    `json:"env"`
	IPStackAccessKey         string    `json:"ipStackAccessKey"`
	MysqlDSL                 string    `json:"mysqlDSL"`
	CorePath                 string    `json:"corePath"`
	SupportedCoins           []string  `json:"coins"`
	CMCAccessKey             string    `json:"cmcAccessKey"`
	TreasuryAddresses        []string  `json:"treasuryAddresses"`
	IsAuthEnabled            bool      `json:"enableAuth"`
	AuthPairs                []Auth    `json:"authPairs"`
	Jobs                     JobConfig `json:"jobs"`
	DbSyncOnStart            bool      `json:"dbSyncOnStart"`
}

type Auth struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type JobConfig struct {
	CMCInterval     string `json:"cmcInterval"`
	MNCache         string `json:"mnCache"`
	GovCache        string `json:"govCache"`
	HistoryCache    string `json:"historyCache"`
	TPOSCache       string `json:"tposCache"`
	LNCache         string `json:"lightningCache"`
	TreasuryHistory string `json:"treasuryHistory"`
	TPOSHistory     string `json:"tposHistory"`
	MNHistory       string `json:"mnHistory"`
	LNHistory       string `json:"lightningHistory"`
}

func (c *Config) GetLNCLIPath() string {
	return c.LNCLIPath
}

func (c *Config) GetCoreWalletPath() string {
	return c.CorePath
}
