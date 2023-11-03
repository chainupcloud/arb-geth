package arbitrum

import (
	"context"

	"github.com/chainupcloud/arb-geth/arbitrum_types"
	"github.com/chainupcloud/arb-geth/core"
	"github.com/chainupcloud/arb-geth/core/types"
)

type ArbInterface interface {
	PublishTransaction(ctx context.Context, tx *types.Transaction, options *arbitrum_types.ConditionalOptions) error
	BlockChain() *core.BlockChain
	ArbNode() interface{}
}
