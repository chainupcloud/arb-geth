// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package core implements the Ethereum consensus protocol.
package core

import (
	"errors"
	"time"

	"github.com/chainupcloud/arb-geth/core/state"
	"github.com/chainupcloud/arb-geth/core/types"
	"github.com/chainupcloud/arb-geth/log"
	"github.com/chainupcloud/arb-geth/rpc"
)

// WriteBlockAndSetHeadWithTime also counts processTime, which will cause intermittent TrieDirty cache writes
func (bc *BlockChain) WriteBlockAndSetHeadWithTime(block *types.Block, receipts []*types.Receipt, logs []*types.Log, state *state.StateDB, emitHeadEvent bool, processTime time.Duration) (status WriteStatus, err error) {
	if !bc.chainmu.TryLock() {
		return NonStatTy, errChainStopped
	}
	defer bc.chainmu.Unlock()
	bc.gcproc += processTime
	return bc.writeBlockAndSetHead(block, receipts, logs, state, emitHeadEvent)
}

func (bc *BlockChain) ReorgToOldBlock(newHead *types.Block) error {
	bc.wg.Add(1)
	defer bc.wg.Done()
	locked := bc.chainmu.TryLock()
	if !locked {
		return errors.New("couldn't catch lock to reorg")
	}
	defer bc.chainmu.Unlock()
	oldHead := bc.CurrentBlock()
	if oldHead.Hash() == newHead.Hash() {
		return nil
	}
	bc.writeHeadBlock(newHead)
	err := bc.reorg(oldHead, newHead)
	if err != nil {
		return err
	}
	bc.chainHeadFeed.Send(ChainHeadEvent{Block: newHead})
	return nil
}

func (bc *BlockChain) ClipToPostNitroGenesis(blockNum rpc.BlockNumber) (rpc.BlockNumber, rpc.BlockNumber) {
	currentBlock := rpc.BlockNumber(bc.CurrentBlock().Number.Uint64())
	nitroGenesis := rpc.BlockNumber(bc.Config().ArbitrumChainParams.GenesisBlockNum)
	if blockNum == rpc.LatestBlockNumber || blockNum == rpc.PendingBlockNumber {
		blockNum = currentBlock
	}
	if blockNum > currentBlock {
		blockNum = currentBlock
	}
	if blockNum < nitroGenesis {
		blockNum = nitroGenesis
	}
	return blockNum, currentBlock
}

func (bc *BlockChain) RecoverState(block *types.Block) error {
	if bc.HasState(block.Root()) {
		return nil
	}
	log.Warn("recovering block state", "num", block.Number(), "hash", block.Hash(), "root", block.Root())
	_, err := bc.recoverAncestors(block)
	return err
}
