package eth

import (
	"context"

	"github.com/chainupcloud/arb-geth/core"
	"github.com/chainupcloud/arb-geth/core/state"
	"github.com/chainupcloud/arb-geth/core/types"
	"github.com/chainupcloud/arb-geth/core/vm"
	"github.com/chainupcloud/arb-geth/eth/tracers"
	"github.com/chainupcloud/arb-geth/ethdb"
)

func NewArbEthereum(
	blockchain *core.BlockChain,
	chainDb ethdb.Database,
) *Ethereum {
	return &Ethereum{
		blockchain: blockchain,
		chainDb:    chainDb,
	}
}

func (eth *Ethereum) StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64) (*core.Message, vm.BlockContext, *state.StateDB, tracers.StateReleaseFunc, error) {
	return eth.stateAtTransaction(ctx, block, txIndex, reexec)
}
