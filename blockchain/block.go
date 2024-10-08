package blockchain

import (
	"blockchain-simulator/types"
	"time"
)

// * : 메모리 주소에 저장된 값에 접근
// & : 주소를 나탬
func NewBlock(index int, data string, previousHash string, transactions []types.Transaction, merkleRoot string) *types.Block {
	timestamp := time.Now().String()
	block := &types.Block{
		Index:        index,
		Timestamp:    timestamp,
		Data:         data,
		PreviousHash: previousHash,
		Transactions: transactions,
		MerkleRoot:   merkleRoot,
	}
	return block
}

func CreateGenesisBlock() *types.Block {
	return NewBlock(0, "Genesis Block", "", []types.Transaction{}, "")
}
