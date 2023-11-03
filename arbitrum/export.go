package arbitrum

import (
	"context"

	"github.com/chainupcloud/arb-geth/common/hexutil"
	"github.com/chainupcloud/arb-geth/core"
	"github.com/chainupcloud/arb-geth/internal/ethapi"
	"github.com/chainupcloud/arb-geth/rpc"
)

type TransactionArgs = ethapi.TransactionArgs

func EstimateGas(ctx context.Context, b ethapi.Backend, args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, gasCap uint64) (hexutil.Uint64, error) {
	return ethapi.DoEstimateGas(ctx, b, args, blockNrOrHash, gasCap)
}

func NewRevertReason(result *core.ExecutionResult) error {
	return ethapi.NewRevertError(result)
}
