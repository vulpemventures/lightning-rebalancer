package util

import "os"

var (
	// HTTPPort is the default
	HTTPPort = "3000"
	// LNDPort is the default
	LNDPort = "10009"
	// LNDHost is the default
	LNDHost = "localhost"
	// TLSCertPath is the default
	TLSCertPath = "/home/bitcoin/.lnd/tls.cert"
	// MACPath is the default
	MACPath = "/home/bitcoin/.lnd/data/chain/bitcoin/testnet/admin.macaroon"
)

// Config represents the configuration instance
type Config struct {
	HTTPPort    string
	LNDHost     string
	LNDPort     string
	MACPath     string
	TLSCertPath string
}

// GetConfigFromEnv override default vars using if present the ones from enviroment
func GetConfigFromEnv() *Config {
	config := &Config{HTTPPort, LNDHost, LNDPort, MACPath, TLSCertPath}

	if _, present := os.LookupEnv("HTTP_PORT"); present {
		config.HTTPPort = os.Getenv("HTTP_PORT")
	}

	if _, present := os.LookupEnv("LND_PORT"); present {
		config.LNDPort = os.Getenv("LND_PORT")
	}

	if _, present := os.LookupEnv("LND_HOST"); present {
		config.LNDHost = os.Getenv("LND_HOST")
	}

	if _, present := os.LookupEnv("TLS_PATH"); present {
		config.TLSCertPath = os.Getenv("TLS_PATH")
	}

	if _, present := os.LookupEnv("MAC_PATH"); present {
		config.MACPath = os.Getenv("MAC_PATH")
	}

	return config

}
