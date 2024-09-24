package chain

import (
	"blockchain-simulator/internal/utils"
	"blockchain-simulator/pkg/block"
	"blockchain-simulator/pkg/transaction"
	"blockchain-simulator/pkg/types"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

// chain에 대한 정보
type Blockchain struct {
	Blocks           []*types.Block
	Validators       map[string]int
	memPool          []transaction.Transaction
	BlockGasLimit    int                 // 블록 가스 한도
	TransactionGas   int                 // 트랜잭션당 가스 소비량
	MiningInterval   time.Duration       // 블록 생성 주기
	balances         map[string]*big.Int // 발신자 주소를 키로 사용하여 잔액을 저장
	ValidatorRewards map[string]float64
	RewardPerBlock   float64
	MaxBlockSize int
	MiniValidators int
}

// 새로운 체인을 만드는 함수
func NewBlockchain(blockGasLimit, transactionGas int, miningInterval time.Duration, maxBlockSize int, minValidators int) *Blockchain {
	bc := &Blockchain{
		Blocks:         []*types.Block{block.CreateGenesisBlock()},
		Validators:     make(map[string]int),
		BlockGasLimit:  blockGasLimit,
		TransactionGas: transactionGas,
		MiningInterval: miningInterval,
		MaxBlockSize: maxBlockSize,
		MiniValidators: minValidators,
	}

	go bc.startMining() // 백그라운드에서 블록 생성 시작
	return bc
}

// 올바른 서명이 되었는지를 확인하는 코드
func (bc *Blockchain) verifySignature(tx transaction.Transaction) bool {
	// 거래의 메시지를 추출
	msg := tx.Message()
	// 메시지의 해시값을 생성
	hash := sha256.Sum256(msg)

	pubKey := &tx.From
	r := new(big.Int).SetBytes(tx.Signature[:len(tx.Signature)/2])
	s := new(big.Int).SetBytes(tx.Signature[len(tx.Signature)/2:])

	return ecdsa.Verify(pubKey, hash[:], r, s)
}

// // 금액 설정
// func (bc *Blockchain) setBalance(address string, amount *big.Int) {
// 	bc.balances[address] = amount
// }
// 
// // 금액 가져오기
// func (bc *Blockchain) getBalance(address string) *big.Int {
// 	balance, exists := bc.balances[address]
// 	if !exists {
// 		return big.NewInt(0)
// 	}
// 	return balance
// }

// 잔액이 충분한지 확인하는 함수
func (bc *Blockchain) hasSufficientBalance(tx transaction.Transaction) bool {
	balance := bc.balances[tx.From.X.String()]

	amount := new(big.Int).SetInt64(int64(tx.Amount))
	gas := new(big.Int).SetInt64(int64(tx.Gas))

	total := new(big.Int).Add(amount, gas)

	return balance.Cmp(total) >= 0
}

// 중복 거래를 확인하는 함수
// 동일한 트랜젝션 아이디가 존재하면 true 반환
func (bc *Blockchain) isDuplicate(tx transaction.Transaction) bool {
	for _, existingTx := range bc.memPool {
		if existingTx.ID == tx.ID {
			return true
		}
	}

	return false
}

// 트랜젝션을 확인하는 함수
func (bc *Blockchain) validateTransaction(tx transaction.Transaction) error {
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

// 트랜젝션을 mempool에 저장
// 검증하는 로직이 추가적으로 필요할 것으로 판단이 됨.
func (bc *Blockchain) AddTransaction(tx transaction.Transaction, gas int) error {
	err := bc.validateTransaction(tx)
	if err != nil {
		return err
	}

	// 유효한 거래에 대해서 mempool에 저장
	bc.memPool = append(bc.memPool, tx)
	return nil
}

// PoS 알고리즘을 활용하여 Validator를 알아냄
func (bc *Blockchain) SelectValidator() string {
	totalStake := 0

	for _, stake := range bc.Validators {
		totalStake += stake
	}

	randValue := rand.Intn(totalStake)
	cumulaticeStake := 0

	for validator, stake := range bc.Validators {
		cumulaticeStake += stake
		if randValue < cumulaticeStake {
			return validator
		}
	}

	return ""
}

// 작업 시작
func (bc *Blockchain) startMining() {
	for {
		time.Sleep(bc.MiningInterval)

		validatorID := bc.SelectValidator()
		bc.CreateBlockFromTransactions(validatorID)
	}
}

// reward 분배 (구현)
// func (bc *Blockchain) distributeRewards(validatorID string) {
// 	if _, exists := bc.Validators[validatorID]; exists {
// 		if bc.ValidatorRewards == nil {
// 			bc.ValidatorRewards = make(map[string]float64)
// 		}
// 		bc.ValidatorRewards[validatorID] += bc.RewardPerBlock
// 	} else {
// 		fmt.Println("Invalid validator, reward not distributed")
// 	}
// }

// memPool에서 트랜젝션을 선택하여 Merkle Root를 기반으로 블록을 생성하는 함수
// 여기서 구현이 됨.
func (bc *Blockchain) CreateBlockFromTransactions(validatorID string) {
	if len(bc.memPool) == 0 {
		return
	}

	if !bc.isValidValidator(validatorID) {
		return
	}

	// 가스 한도에 맞게 트랜잭션 선택
	var selectedTxs []transaction.Transaction
	var totalGas int

	for _, tx := range bc.memPool {
		if totalGas+bc.TransactionGas > bc.BlockGasLimit {
			break
		}
		selectedTxs = append(selectedTxs, tx)
		totalGas += bc.TransactionGas
	}

	if len(selectedTxs) == 0 {
		return
	}

	data := "Transactions:"
	for _, tx := range selectedTxs {
		data += fmt.Sprintf(" %v;", tx)
	}

	_, err := bc.CreateBlockWithMerkleRoot(selectedTxs, validatorID)

	if err != nil {
		return
	}

	bc.memPool = bc.memPool[len(selectedTxs):] // 선택한 트랜잭션을 큐에서 제거
}

func (bc *Blockchain) CreateBlockWithMerkleRoot (transactions []transaction.Transaction, validatorID string) (*types.Block, error) {
	if !bc.isValidValidator(validatorID) {
		return nil, errors.New("유효하지 않은 검증자")
	}

	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	newIndex := previousBlock.Index + 1

	// 트랜젝션 해시 계산
	txHashes := make([]string, len(transactions))
	for i, tx := range transactions {
		txHashes[i] = utils.CalculateTransactionHash(tx)
	}

	merkleRoot := utils.CalculateMerkelRoot(txHashes)

	newBlock := block.NewBlock(newIndex, "", previousBlock.MerkleRoot, transactions, merkleRoot)

	isValidNewBlock := bc.isValidNewBlock(newBlock, previousBlock)
	if !isValidNewBlock {
		return nil, nil
	}
	bc.Blocks = append(bc.Blocks, newBlock)

	return newBlock, nil
}

// 블록에 대한 유효성을 검사
func (bc *Blockchain) isValidNewBlock(newBlock, previousBlock *types.Block) bool {
	// 블록의 인덱스 검증
	if previousBlock.Index+1 != newBlock.Index {
		return false
	}

	// 블록의 해시 검증
	if newBlock.PreviousHash != previousBlock.MerkleRoot {
		return false
	}

	// 블록의 해시가 올바르게 계산되었는지 검증 (이거 수정해야됨.)
	if newBlock.MerkleRoot != utils.CalculateHash(newBlock) {
		return false
	}

	return true
}

// 검증자가 유요한자 확인
func (bc *Blockchain) isValidValidator(validatorID string) bool {
	_, exists := bc.Validators[validatorID]
	return exists
}
