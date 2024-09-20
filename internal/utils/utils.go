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

// Merkle Root를 활용한 블록의 해시 생성
func hashPair(left, right string) string {
	h := sha256.New()
	h.Write([]byte(left + right))

	return hex.EncodeToString(h.Sum(nil))
}

// Merkel Root를 활용한 해시 생성
func CalculateMerkelRoot(txHashes []string) string {
	if len(txHashes) == 0 {
		return ""
	}

	for len(txHashes) > 1 {
		var newLevel []string

		for i := 0; i < len(txHashes); i += 2 {
			if i + 1 < len(txHashes) {
				newLevel = append(newLevel, hashPair(txHashes[i], txHashes[i+1]))
			} else {
				newLevel = append(newLevel, txHashes[i])
			}
		}

		txHashes = newLevel
	}

	return txHashes[0]
}