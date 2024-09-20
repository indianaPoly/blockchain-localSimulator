package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// Transaction 인터페이스 정의
type Transaction struct {
	ID        string
	From      ecdsa.PublicKey
	To        ecdsa.PublicKey
	Amount    int
	Gas       int
	Signature []byte
}

// 서명된 거래의 메세지를 생성
func (tx *Transaction) Message() []byte {
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
func (tx *Transaction) FromAddress() string {
	return hashPublicKey(tx.From)
}
