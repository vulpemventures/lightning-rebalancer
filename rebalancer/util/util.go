package util

import (
	"encoding/json"
	"io"

	"github.com/btcsuite/btcd/chaincfg"
)

// GetNetworkTypeByString retunrs a network object given one of mainnet | testnet
func GetNetworkTypeByString(key string) *chaincfg.Params {

	if key == "mainnet" {
		return &chaincfg.MainNetParams
	}

	return &chaincfg.TestNet3Params
}

// GetRequestBody retuns the body as a map
func GetRequestBody(body io.ReadCloser) map[string]string {
	decoder := json.NewDecoder(body)
	var decodedBody map[string]string
	decoder.Decode(&decodedBody)

	return decodedBody
}
