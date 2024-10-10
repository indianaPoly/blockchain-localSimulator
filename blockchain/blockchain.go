package blockchain

import (
	"blockchain-simulator/types"
	"blockchain-simulator/utils"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// chain에 대한 정보
type Blockchain struct {
	Blocks           []*types.Block
	Validators       map[string]int
	memPool          []types.Transaction
	BlockGasLimit    int                 // 블록 가스 한도
	TransactionGas   int                 // 트랜잭션당 가스 소비량
	MiningInterval   time.Duration       // 블록 생성 주기
	balances         map[string]*big.Int // 발신자 주소를 키로 사용하여 잔액을 저장
	ValidatorRewards map[string]float64
	RewardPerBlock   float64
	MaxBlockSize     int
	MiniValidators   int
}

// 새로운 체인을 만드는 함수
func NewBlockchain(blockGasLimit, transactionGas int, miningInterval time.Duration, maxBlockSize int, minValidators int) *Blockchain {
	bc := &Blockchain{
		Blocks:         []*types.Block{CreateGenesisBlock()},
		Validators:     make(map[string]int),
		BlockGasLimit:  blockGasLimit,
		TransactionGas: transactionGas,
		MiningInterval: miningInterval,
		MaxBlockSize:   maxBlockSize,
		MiniValidators: minValidators,
	}

	go bc.startMining() // 백그라운드에서 블록 생성 시작
	return bc
}

// 트랜젝션을 mempool에 저장
// 검증하는 로직이 추가적으로 필요할 것으로 판단이 됨.
func (bc *Blockchain) AddTransaction(tx types.Transaction, gas int) error {
	err := bc.validateTransaction(tx)
	if err != nil {
		return err
	}

	// 유효한 거래에 대해서 mempool에 저장
	bc.memPool = append(bc.memPool, tx)
	return nil
}

// 작업 시작
func (bc *Blockchain) startMining() {
	for {
		time.Sleep(bc.MiningInterval)

		validatorID := bc.SelectValidator()
		bc.CreateBlockFromTransactions(validatorID)
	}
}

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
	var selectedTxs []types.Transaction
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

func (bc *Blockchain) CreateBlockWithMerkleRoot(transactions []types.Transaction, validatorID string) (*types.Block, error) {
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

	newBlock := NewBlock(newIndex, "", previousBlock.MerkleRoot, transactions, merkleRoot)

	isValidNewBlock := bc.isValidNewBlock(newBlock, previousBlock)
	if !isValidNewBlock {
		return nil, nil
	}
	bc.Blocks = append(bc.Blocks, newBlock)

	return newBlock, nil
}