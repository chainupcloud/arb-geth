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

// Package ethconfig contains the configuration of the ETH and LES protocols.
package ethconfig

import (
	"errors"
	"time"

	"github.com/chainupcloud/arb-geth/common"
	"github.com/chainupcloud/arb-geth/consensus"
	"github.com/chainupcloud/arb-geth/consensus/beacon"
	"github.com/chainupcloud/arb-geth/consensus/clique"
	"github.com/chainupcloud/arb-geth/consensus/ethash"
	"github.com/chainupcloud/arb-geth/core"
	"github.com/chainupcloud/arb-geth/core/txpool"
	"github.com/chainupcloud/arb-geth/eth/downloader"
	"github.com/chainupcloud/arb-geth/eth/gasprice"
	"github.com/chainupcloud/arb-geth/ethdb"
	"github.com/chainupcloud/arb-geth/miner"
	"github.com/chainupcloud/arb-geth/params"
)

// FullNodeGPO contains default gasprice oracle settings for full node.
var FullNodeGPO = gasprice.Config{
	Blocks:           20,
	Percentile:       60,
	MaxHeaderHistory: 1024,
	MaxBlockHistory:  1024,
	MaxPrice:         gasprice.DefaultMaxPrice,
	IgnorePrice:      gasprice.DefaultIgnorePrice,
}

// LightClientGPO contains default gasprice oracle settings for light client.
var LightClientGPO = gasprice.Config{
	Blocks:           2,
	Percentile:       60,
	MaxHeaderHistory: 300,
	MaxBlockHistory:  5,
	MaxPrice:         gasprice.DefaultMaxPrice,
	IgnorePrice:      gasprice.DefaultIgnorePrice,
}

// Defaults contains default settings for use on the Ethereum main net.
var Defaults = Config{
	SyncMode:                downloader.SnapSync,
	NetworkId:               1,
	TxLookupLimit:           2350000,
	LightPeers:              100,
	UltraLightFraction:      75,
	DatabaseCache:           512,
	TrieCleanCache:          154,
	TrieCleanCacheJournal:   "triecache",
	TrieCleanCacheRejournal: 60 * time.Minute,
	TrieDirtyCache:          256,
	TrieTimeout:             60 * time.Minute,
	SnapshotCache:           102,
	FilterLogCacheSize:      32,
	Miner:                   miner.DefaultConfig,
	TxPool:                  txpool.DefaultConfig,
	RPCGasCap:               50000000,
	RPCEVMTimeout:           5 * time.Second,
	GPO:                     FullNodeGPO,
	RPCTxFeeCap:             1, // 1 ether
}

//go:generate go run github.com/fjl/gencodec -type Config -formats toml -out gen_config.go

// Config contains configuration options for of the ETH and LES protocols.
type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the Ethereum main net block is used.
	Genesis *core.Genesis `toml:",omitempty"`

	// Protocol options
	NetworkId uint64 // Network ID to use for selecting peers to connect to
	SyncMode  downloader.SyncMode

	// This can be set to list of enrtree:// URLs which will be queried for
	// for nodes to connect to.
	EthDiscoveryURLs  []string
	SnapDiscoveryURLs []string

	NoPruning  bool // Whether to disable pruning and flush everything to disk
	NoPrefetch bool // Whether to disable prefetching and only load state on demand

	TxLookupLimit uint64 `toml:",omitempty"` // The maximum number of blocks from head whose tx indices are reserved.

	// RequiredBlocks is a set of block number -> hash mappings which must be in the
	// canonical chain of all remote peers. Setting the option makes geth verify the
	// presence of these blocks for every new peer connection.
	RequiredBlocks map[uint64]common.Hash `toml:"-"`

	// Light client options
	LightServ        int  `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	LightIngress     int  `toml:",omitempty"` // Incoming bandwidth limit for light servers
	LightEgress      int  `toml:",omitempty"` // Outgoing bandwidth limit for light servers
	LightPeers       int  `toml:",omitempty"` // Maximum number of LES client peers
	LightNoPrune     bool `toml:",omitempty"` // Whether to disable light chain pruning
	LightNoSyncServe bool `toml:",omitempty"` // Whether to serve light clients before syncing

	// Ultra Light client options
	UltraLightServers      []string `toml:",omitempty"` // List of trusted ultra light servers
	UltraLightFraction     int      `toml:",omitempty"` // Percentage of trusted servers to accept an announcement
	UltraLightOnlyAnnounce bool     `toml:",omitempty"` // Whether to only announce headers, or also serve them

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	DatabaseFreezer    string

	TrieCleanCache          int
	TrieCleanCacheJournal   string        `toml:",omitempty"` // Disk journal directory for trie cache to survive node restarts
	TrieCleanCacheRejournal time.Duration `toml:",omitempty"` // Time interval to regenerate the journal for clean cache
	TrieDirtyCache          int
	TrieTimeout             time.Duration
	SnapshotCache           int
	Preimages               bool

	// This is the number of blocks for which logs will be cached in the filter system.
	FilterLogCacheSize int

	// Mining options
	Miner miner.Config

	// Transaction pool options
	TxPool txpool.Config

	// Gas Price Oracle options
	GPO gasprice.Config

	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool

	// Miscellaneous options
	DocRoot string `toml:"-"`

	// RPCGasCap is the global gas cap for eth-call variants.
	RPCGasCap uint64

	// RPCEVMTimeout is the global timeout for eth-call.
	RPCEVMTimeout time.Duration

	// RPCTxFeeCap is the global transaction fee(price * gaslimit) cap for
	// send-transaction variants. The unit is ether.
	RPCTxFeeCap float64

	// OverrideCancun (TODO: remove after the fork)
	OverrideCancun *uint64 `toml:",omitempty"`
}

// CreateConsensusEngine creates a consensus engine for the given chain config.
// Clique is allowed for now to live standalone, but ethash is forbidden and can
// only exist on already merged networks.
func CreateConsensusEngine(config *params.ChainConfig, db ethdb.Database) (consensus.Engine, error) {
	// If proof-of-authority is requested, set it up
	if config.Clique != nil {
		return beacon.New(clique.New(config.Clique, db)), nil
	}
	// If defaulting to proof-of-work, enforce an already merged network since
	// we cannot run PoW algorithms and more, so we cannot even follow a chain
	// not coordinated by a beacon node.
	if !config.TerminalTotalDifficultyPassed {
		return nil, errors.New("ethash is only supported as a historical component of already merged networks")
	}
	return beacon.New(ethash.NewFaker()), nil
}
