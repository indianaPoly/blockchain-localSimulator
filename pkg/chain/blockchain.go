package chain

import (
	"blockchain-simulator/internal/utils"
	"blockchain-simulator/pkg/block"
	"blockchain-simulator/pkg/transaction"
	"blockchain-simulator/pkg/types"
	"fmt"
	"math/rand"
	"time"
)

// chain에 대한 정보
type Blockchain struct {
	Blocks         []*types.Block
	Validators     map[string]int
	memPool        []transaction.Transaction
	BlockGasLimit  int           // 블록 가스 한도
	TransactionGas int           // 트랜잭션당 가스 소비량
	MiningInterval time.Duration // 블록 생성 주기
}

// 새로운 체인을 만드는 함수
func NewBlockchain(blockGasLimit, transactionGas int, miningInterval time.Duration) *Blockchain {
	bc := &Blockchain{
		Blocks:         []*types.Block{block.CreateGenesisBlock()},
		Validators:     make(map[string]int),
		BlockGasLimit:  blockGasLimit,
		TransactionGas: transactionGas,
		MiningInterval: miningInterval,
	}

	go bc.startMining() // 백그라운드에서 블록 생성 시작
	return bc
}

// 트랜젝션을 mempool에 저장
func (bc *Blockchain) AddTransaction(tx transaction.Transaction, gas int) {
	bc.memPool = append(bc.memPool, tx)
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

// memPool에서 트랜젝션을 선택하여 블록을 생성하는 함수
func (bc *Blockchain) CreateBlockFromTransactions(validatorID string) {
	if len(bc.memPool) == 0 {
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

	bc.AddBlock(data, validatorID)
	bc.memPool = bc.memPool[len(selectedTxs):] // 선택한 트랜잭션을 큐에서 제거
}

// 블록을 추가하는 함수
func (bc *Blockchain) AddBlock(data string, validatorID string) {
	if !bc.isValidValidator(validatorID) {
		fmt.Println("Invalid validator.")
		return
	}

	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	newIndex := previousBlock.Index + 1
	newBlock := block.NewBlock(newIndex, data, previousBlock.Hash)

	// 블록의 해시를 검증
	if !bc.isValidNewBlock(newBlock, previousBlock) {
		fmt.Println("Invalid block, not adding to blockchain.")
		return
	}

	bc.Blocks = append(bc.Blocks, newBlock)
}

// 블록에 대한 유효성을 검사
func (bc *Blockchain) isValidNewBlock(newBlock, previousBlock *types.Block) bool {
	// 블록의 인덱스 검증
	if previousBlock.Index+1 != newBlock.Index {
		return false
	}

	// 블록의 해시 검증
	if newBlock.PreviousHash != previousBlock.Hash {
		return false
	}

	// 블록의 해시가 올바르게 계산되었는지 검증
	if newBlock.Hash != utils.CalculateHash(newBlock) {
		return false
	}

	return true
}

// 검증자가 유요한자 확인
func (bc *Blockchain) isValidValidator(validatorID string) bool {
	_, exists := bc.Validators[validatorID]
	return exists
}
