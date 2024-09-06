package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Block struct {
	Index int
	Timestamp string
	Data string
	PreviousHash string
	Hash string
}

func NewBlock(index int, data string, previousHash string) *Block {
	timestamp := time.Now().String()
	block := &Block{
		Index: index,
		Timestamp: timestamp,
		Data: data,
		PreviousHash: previousHash,
		Hash: "",
	}
	block.Hash = block.calculateHash()
	return block
}

func (b *Block) calculateHash() string {
	record := strconv.Itoa(b.Index) + b.Timestamp + b.Data + b.PreviousHash
	hash := sha256.New()
	hash.Write([]byte(record))
	hased := hash.Sum(nil)
	return hex.EncodeToString(hased)
}