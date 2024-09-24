package types

import "blockchain-simulator/pkg/transaction"

type Block struct {
	Index        int
	Timestamp    string
	Transactions [] transaction.Transaction
	Data         string
	PreviousHash string
	MerkleRoot   string
}
