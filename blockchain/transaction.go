package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"blockchain-simulator/types"

	"golang.org/x/crypto/sha3"
)

type TransactionWrppter struct {
	types.Transaction
}

// 서명된 거래의 메세지를 생성
func (tx *TransactionWrppter) Message() []byte {
	data := tx.ID + tx.From.X.String() + tx.From.Y.String() + tx.To.X.String() + tx.To.Y.String() + fmt.Sprint(tx.Amount) + fmt.Sprint(tx.Gas)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// 공개키를 기반으로 주소를 알아냄
func hashPublicKey(pubkey ecdsa.PublicKey) string {
	pubKeyBytes := append(pubkey.X.Bytes(), pubkey.Y.Bytes()...)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKeyBytes)
	addressBytes := hash.Sum(nil)

	return hex.EncodeToString(addressBytes[:])
}

// 발신자 주소를 반환하는 메소드 (공개키를 해시하여 주소를 생성 )
func (tx *TransactionWrppter) FromAddress() string {
	return hashPublicKey(tx.From)
}
