package block

import (
	"blockchain-simulator/internal/utils"
	"blockchain-simulator/pkg/types"
	"time"
)

// * : 메모리 주소에 저장된 값에 접근
// & : 주소를 나탬
func NewBlock(index int, data string, previousHash string) *types.Block {
	timestamp := time.Now().String()
	block := &types.Block{
		Index:        index,
		Timestamp:    timestamp,
		Data:         data,
		PreviousHash: previousHash,
		Hash:         "",
	}
	block.Hash = utils.CalculateHash(block)
	return block
}

func CreateGenesisBlock() *types.Block {
	return NewBlock(0, "Genesis Block", "")
}
