package entities

import (
	"math/big"

	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type FetchedBlock struct {
	BlockNumber  uint64
	// Block        *ethtypes.Block
	Transactions []*ethcommon.Hash
	Jobs         []*entities.Job
}

func NewFetchedBlock(txs []*ethcommon.Hash, blockNumber uint64) *FetchedBlock {
	return &FetchedBlock{
		Transactions: txs,
		BlockNumber: blockNumber,
		Jobs: []*entities.Job{},
	}
}

func (b *FetchedBlock) AppendJobs(jobs []*entities.Job) *FetchedBlock {
	b.Jobs = append(b.Jobs, jobs...)
	return b
}

func (b *FetchedBlock) TransactionHashes() []string {
	hashes := []string{}
	for _, tx := range b.Transactions {
		hashes = append(hashes, tx.Hex())
	}
	
	return hashes
}

func (b *FetchedBlock) JobsUUIDs() []string {
	uuids := []string{}
	for _, job := range b.Jobs {
		uuids = append(uuids, job.UUID)
	}
	
	return uuids
}

func (b *FetchedBlock) Number() *big.Int {
	return new(big.Int).SetUint64(b.BlockNumber)
}


func (b *FetchedBlock) NumberU64() uint64 {
	return b.BlockNumber
}
