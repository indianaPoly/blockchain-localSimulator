package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

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

// 잔액이 충분한지 확인하는 함수
func (bc *Blockchain) hasSufficientBalance(tx types.Transaction) bool {
	balance := bc.balances[tx.From.X.String()]
	total := new(big.Int).Add(big.NewInt(int64(tx.Amount)), big.NewInt(int64(tx.Gas)))
	return balance.Cmp(total) >= 0
}

// 중복 거래를 확인하는 함수
// 동일한 트랜젝션 아이디가 존재하면 true 반환
func (bc *Blockchain) isDuplicate(tx types.Transaction) bool {
	for _, existingTx := range bc.memPool {
		if existingTx.ID == tx.ID {
			return true
		}
	}

	return false
}

// 올바른 서명이 되었는지를 확인하는 코드
func (bc *Blockchain) verifySignature(tx types.Transaction) bool {
	hash := sha256.Sum256([]byte(tx.ID)) // 트랜잭션 ID로 해시 생성
	pubKey := &tx.From
	r := new(big.Int).SetBytes(tx.Signature[:len(tx.Signature)/2])
	s := new(big.Int).SetBytes(tx.Signature[len(tx.Signature)/2:])
	return ecdsa.Verify(pubKey, hash[:], r, s)
}

// 트랜젝션을 확인하는 함수
func (bc *Blockchain) validateTransaction(tx types.Transaction) error {
	if tx.Amount <= 0 {
		return errors.New("보내는 금액이 0보다 작습니다")
	}

	if tx.Gas < 1 {
		return errors.New("가스값이 최소 단위보다 작습니다.")
	}

	if !bc.verifySignature(tx) {
		return errors.New("서명이 올바르지 않습니다.")
	}

	if !bc.hasSufficientBalance(tx) {
		return errors.New("insufficient balacnce for tx")
	}

	if bc.isDuplicate(tx) {
		return errors.New("tx is a duplicate")
	}

	return nil
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
