package blockchain

import (
	"blockchain-simulator/types"
	"blockchain-simulator/utils"
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

// 블록에 대한 유효성을 검사
func (bc *Blockchain) isValidNewBlock(newBlock, previousBlock *types.Block) bool {
	// 블록의 인덱스 검증
	if previousBlock.Index+1 != newBlock.Index {
		return false
	}

	// 블록의 해시 검증
	if newBlock.PreviousHash != previousBlock.MerkleRoot {
		return false
	}

	// 블록의 해시가 올바르게 계산되었는지 검증 (이거 수정해야됨.)
	if newBlock.MerkleRoot != utils.CalculateHash(newBlock) {
		return false
	}

	return true
}