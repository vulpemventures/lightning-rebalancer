package util

import "github.com/btcsuite/btcd/chaincfg"

func getNetworkTypeByString(key string) *chaincfg.Params {

	if key == "mainnet" {
		return &chaincfg.MainNetParams
	}

	return &chaincfg.TestNet3Params
}
