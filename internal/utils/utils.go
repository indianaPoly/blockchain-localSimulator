package utils

import (
	"blockchain-simulator/pkg/types"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func CalculateHash(b *types.Block) string {
	record := strconv.Itoa(b.Index) + b.Timestamp + b.Data + b.PreviousHash
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}