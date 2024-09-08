package block

import (
	"blockchain-simulator/internal/utils"
	"time"
)

type Block struct {
	Index        int
	Timestamp    string
	Data         string
	PreviousHash string
	Hash         string
}

// * : 메모리 주소에 저장된 값에 접근
// & : 주소를 나탬
func NewBlock(index int, data string, previousHash string) *Block {
	timestamp := time.Now().String()
	block := &Block{
		Index:        index,
		Timestamp:    timestamp,
		Data:         data,
		PreviousHash: previousHash,
		Hash:         "",
	}
	block.Hash = utils.CalculateHash()
	return block
}