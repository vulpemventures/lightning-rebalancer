package processor

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"

	"github.com/vulpemventures/lightning-rebalancer/rebalancer/util"
)

//Processor represents the current instance
type Processor struct {
	queque    []*btcutil.Address
	network   *chaincfg.Params
	confirmed map[string]*btcutil.Tx
}

//NewProcessor returns an instance that watches payments belongs to known addresses
func NewProcessor(network string) *Processor {
	return &Processor{
		[]*btcutil.Address{},
		util.getNetworkTypeByString(network),
	}
}
