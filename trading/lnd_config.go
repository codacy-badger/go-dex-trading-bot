package trading

import (
	"fmt"
	"os"
)

type Connection struct {
	Host string
	Port int
}

func (c *Connection) Format() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type LNDConfig struct {
	TLSPaths    map[string]string     // map[LTC]-> /path/to/lnd_ltc.tls.cert
	Certs       map[string]string     // map[LTC]-> -----BEGIN CERTIFICATE-----...
	Connections map[string]Connection // map[LTC]{Host: localhost, Port: 10001}
}

func NewConfig() LNDConfig {
	cfg := LNDConfig{}
	cfg.TLSPaths = make(map[string]string, 0)
	cfg.Certs = make(map[string]string, 0)
	cfg.Connections = make(map[string]Connection, 0)
	return cfg
}

func (c *LNDConfig) IsEmtpy() bool {
	return len(c.TLSPaths) == 0
}

func (c *LNDConfig) Add(currency string, certPath string, host string, port int) error {
	if !ValidateCurrency(currency) {
		return fmt.Errorf("currency '%s' not allowed", currency)
	}
	if !FileExists(certPath) {
		return fmt.Errorf("certPath '%s' file does not exist", certPath)
	}
	if host == "" {
		return fmt.Errorf("host is empty")
	}
	if port == 0 {
		return fmt.Errorf("port is empty")
	}

	c.TLSPaths[currency] = certPath
	c.Connections[currency] = Connection{Host: host, Port: port}
	return nil
}

func ValidateCurrency(currency string) bool {
	switch currency {
	case CurrencyLTC:
		return true
	case CurrencyBTC:
		return true
	case CurrencyXSN:
		return true
	}
	return false
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
