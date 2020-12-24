// Copyright (c) 2016 The thaibaoautonomous developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"github.com/incognitochain/go-incognito-sdk/common"
	"sync"
)

const (
	// estimateFeeDepth is the maximum number of blocks before a transaction
	// is confirmed that we want to track.
	estimateFeeDepth = 200
)

// CoinPerKilobyte is number with units of coins per kilobyte.
type CoinPerKilobyte uint64


// observedTransaction represents an observed transaction and some
// additional data required for the fee estimation algorithm.
type observedTransaction struct {
	// A transaction hash.
	hash common.Hash

	// The PRV fee per kilobyte of the transaction in coins.
	feeRate CoinPerKilobyte

	// The token fee per kilobyte of the transaction in coins.
	feeRateForToken map[common.Hash]CoinPerKilobyte

	// The block height when it was observed.
	observed uint64

	// The height of the block in which it was mined.
	// If the transaction has not yet been mined, it is zero.
	mined uint64
}

// registeredBlock has the hash of a block and the list of transactions
// it mined which had been previously observed by the feeEstimator. It
// is used if Rollback is called to reverse the effect of registering
// a block.
type registeredBlock struct {
	hash         common.Hash
	transactions []*observedTransaction
}

// feeEstimator manages the data necessary to create
// fee estimations. It is safe for concurrent access.
type FeeEstimator struct {
	maxRollback uint32
	binSize     int32

	// The maximum number of replacements that can be made in a single
	// bin per block. Default is estimateFeeMaxReplacements
	maxReplacements int32

	// The minimum number of blocks that can be registered with the fee
	// estimator before it will provide answers.
	minRegisteredBlocks uint32

	// The last known height.
	lastKnownHeight uint64

	// The number of blocks that have been registered.
	numBlocksRegistered uint32

	mtx      sync.RWMutex
	observed map[common.Hash]*observedTransaction
	bin      [estimateFeeDepth][]*observedTransaction

	// The cached estimates.
	cached []CoinPerKilobyte

	// Transactions that have been removed from the bins. This allows us to
	// revert in case of an orphaned block.
	dropped []*registeredBlock

	// min fee which be needed for payment on tx(per Kb data)
	limitFee uint64
}

// returns the limit fee of tokenID
// if there is no exchange rate between native token and privacy token, return limit fee of native token
func (ef FeeEstimator) GetLimitFeeForNativeToken() uint64 {
	limitFee := ef.limitFee
	//isFeePToken := false

	//if tokenID != nil {
	//	limitFeePToken, err := ConvertNativeTokenToPrivacyToken(ef.limitFee, tokenID)
	//	if err == nil {
	//		limitFee = limitFeePToken
	//		isFeePToken = true
	//	}
	//}

	return limitFee
}