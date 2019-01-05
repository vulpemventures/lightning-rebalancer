package main

import "os"

var HTTPPort = "3000"
var LNDPort = "10009"
var LNDHost = "localhost"
var TLSCertPath = "/home/bitcoin/.lnd/tls.cert"
var MACPath = "/home/bitcoin/.lnd/data/chain/bitcoin/testnet/admin.macaroon"

func getConfigFromEnv() {
	if _, present := os.LookupEnv("HTTP_PORT"); present {
		HTTPPort = os.Getenv("HTTP_PORT")
	}

	if _, present := os.LookupEnv("LND_PORT"); present {
		LNDPort = os.Getenv("LND_PORT")
	}

	if _, present := os.LookupEnv("LND_HOST"); present {
		LNDHost = os.Getenv("LND_HOST")
	}

	if _, present := os.LookupEnv("TLS_PATH"); present {
		TLSCertPath = os.Getenv("TLS_PATH")
	}

	if _, present := os.LookupEnv("MAC_PATH"); present {
		MACPath = os.Getenv("MAC_PATH")
	}

}
