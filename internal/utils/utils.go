package utils

import (
	"blockchain-simulator/pkg/transaction"
	"blockchain-simulator/pkg/types"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

func SerializePublicKey(pubKey ecdsa.PublicKey) string {
	return fmt.Sprintf("%x%x", pubKey.X, pubKey.Y)
}

func CalculateTransactionHash(tx transaction.Transaction) string {
	fromKey := SerializePublicKey(tx.From)
	toKey := SerializePublicKey(tx.To)

	data := fmt.Sprintf("%s:%s:%s:%d:%d", tx.ID, fromKey, toKey, tx.Amount, tx.Gas)

	hash := sha256.New()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func CalculateHash(b *types.Block) string {
	record := strconv.Itoa(b.Index) + b.Timestamp + b.Data + b.PreviousHash
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}
