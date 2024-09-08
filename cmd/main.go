package main

import (
	"fmt"

	"blockchain-simulator/pkg/chain"
)

func main() {
	// 새로운 블록체인 생성
	blockchain := chain.NewBlockchain()

	// 블록 추가
	blockchain.AddBlock("첫 번째 트랜잭션")
	blockchain.AddBlock("두 번째 트랜잭션")

	// 블록체인 출력
	for _, block := range blockchain.Blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Println()
	}
}
