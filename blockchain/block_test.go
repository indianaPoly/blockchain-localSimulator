package blockchain

import (
	"blockchain-simulator/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func generateSignature(privateKey *ecdsa.PrivateKey, msg []byte) ([]byte, error) {
	hash := sha256.Sum256(msg)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	// 서명 r, s 값을 바이트 배열로 변환하여 저장
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func TestNewBlock(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := privateKey.PublicKey

	// 임의의 메시지 생성
	msg := []byte("test message")

	// 서명 생성
	signature, err := generateSignature(privateKey, msg)
	if err != nil {
		t.Fatalf("Failed to generate signature: %v", err)
	}

	transactions := []types.Transaction{
		{
			ID:        "tx1",
			From:      publicKey,
			To:        publicKey,
			Amount:    10,
			Gas:       1,
			Signature: signature,
		},
	}

	block := NewBlock(1, "Test Data", "previousHash", transactions, "merkleRoot")

	// 블록의 인덱스가 정확이 들어갔는지 확인
	if block.Index != 1 {
		t.Errorf("Expected block index to be 1, got %d", block.Index)
	}

	// 블록에 올바른 데이터가 들어갔는지 확인
	if block.Data != "Test Data" {
		t.Errorf("Expected block data to be 'Test Data', got %s", block.Data)
	}

	// 트랜젝션이 정확하게 들어갔는지 확인
	if len(block.Transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(block.Transactions))
	}

	// 이전 해시가 올바르게 설정되었는지 확인
	if block.PreviousHash != "previousHash" {
		t.Errorf("Expected previous hash to be 'previousHash', got %s", block.PreviousHash)
	}

	// Merkle Root가 올바르게 설정되었는지 확인
	if block.MerkleRoot != "merkleRoot" {
		t.Errorf("Expected Merkle Root to be 'merkleRoot', got %s", block.MerkleRoot)
	}
}

func TestCreateGenesisBlock(t *testing.T) {
	genesisBlock := CreateGenesisBlock()

	// 제네시스 블록의 인덱스가 0인지 확인
	if genesisBlock.Index != 0 {
		t.Errorf("Expected genesis block index to be 0, got %d", genesisBlock.Index)
	}

	// 제네시스 블록의 데이터가 'Genesis Block'인지 확인
	if genesisBlock.Data != "Genesis Block" {
		t.Errorf("Expected genesis block data to be 'Genesis Block', got %s", genesisBlock.Data)
	}

	// 제네시스 블록에 트랜잭션이 비어 있는지 확인
	if len(genesisBlock.Transactions) != 0 {
		t.Errorf("Expected no transactions in genesis block, got %d", len(genesisBlock.Transactions))
	}
}
