// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package digest

import (
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/signer-service/crypto/hash"
)

func CalcTransactionMetaDataHash(metaData *protocol.TransactionsBlockMetadata) primitives.Sha256 {
	return hash.CalcSha256(metaData.Raw())
}

func CalcTransactionsBlockHash(transactionsBlock *protocol.TransactionsBlockContainer) primitives.Sha256 {
	if transactionsBlock == nil || transactionsBlock.Header == nil {
		return nil
	}
	return hash.CalcSha256(transactionsBlock.Header.Raw())
}

func CalcResultsBlockHash(resultsBlock *protocol.ResultsBlockContainer) primitives.Sha256 {
	if resultsBlock == nil || resultsBlock.Header == nil {
		return nil
	}
	return hash.CalcSha256(resultsBlock.Header.Raw())
}

func CalcBlockHash(transactionsBlock *protocol.TransactionsBlockContainer, resultsBlock *protocol.ResultsBlockContainer) primitives.Sha256 {
	if transactionsBlock == nil || resultsBlock == nil {
		return nil
	}
	transactionsBlockHash := CalcTransactionsBlockHash(transactionsBlock)
	resultsBlockHash := CalcResultsBlockHash(resultsBlock)
	return hash.CalcSha256(transactionsBlockHash, resultsBlockHash)
}

// TODO v1 Rewrite without Merkle tree and then rename the function
// See https://tree.taiga.io/project/orbs-network/us/651

func CalcStateDiffHash(stateDiffs []*protocol.ContractStateDiff) (primitives.Sha256, error) {
	stdHashValues := make([][]byte, len(stateDiffs))
	for i := 0; i < len(stateDiffs); i++ {
		stdHashValues[i] = stateDiffs[i].Raw()
	}
	return hash.CalcSha256(stdHashValues...), nil
}

func CalcNewBlockTimestamp(prevBlockTimestamp primitives.TimestampNano, now primitives.TimestampNano) primitives.TimestampNano {
	if now > prevBlockTimestamp {
		return now
	}
	return prevBlockTimestamp + 1
}
