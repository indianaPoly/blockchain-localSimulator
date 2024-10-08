package types

import (
	"crypto/ecdsa"
)

type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Data         string
	PreviousHash string
	MerkleRoot   string
}

// Transaction 인터페이스 정의
type Transaction struct {
	ID        string
	From      ecdsa.PublicKey
	To        ecdsa.PublicKey
	Amount    int
	Gas       int
	Signature []byte
}
