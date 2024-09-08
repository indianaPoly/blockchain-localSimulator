package chain

import (
    "blockchain/pkg/block"
    "blockchain/pkg/transaction"
    "fmt"
)

type Blockchain struct {
    Blocks        []*block.Block
    Validators    map[string]int
    Transactions  []transaction.Transaction
}

func NewBlockchain() *Blockchain {
    return &Blockchain{
        Blocks:     []*block.Block{createGenesisBlock()},
        Validators: make(map[string]int),
    }
}

func createGenesisBlock() *block.Block {
    return block.NewBlock(0, "Genesis Block", "")
}

func (bc *Blockchain) AddTransaction(tx transaction.Transaction) {
    bc.Transactions = append(bc.Transactions, tx)
}

func (bc *Blockchain) CreateBlockFromTransactions(validatorID string) {
    if len(bc.Transactions) == 0 {
        fmt.Println("No transactions to create a block.")
        return
    }

    data := "Transactions:"
    for _, tx := range bc.Transactions {
        data += fmt.Sprintf(" %v;", tx)
    }
    
    bc.AddBlock(data, validatorID)
    bc.Transactions = nil // 트랜잭션 큐 비우기
}

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

func (bc *Blockchain) isValidNewBlock(newBlock, previousBlock *block.Block) bool {
    // 블록의 인덱스 검증
    if previousBlock.Index+1 != newBlock.Index {
        return false
    }

    // 블록의 해시 검증
    if newBlock.PreviousHash != previousBlock.Hash {
        return false
    }

    // 블록의 해시가 올바르게 계산되었는지 검증
    if newBlock.Hash != newBlock.calculateHash() {
        return false
    }

    return true
}

func (bc *Blockchain) isValidValidator(validatorID string) bool {
    _, exists := bc.Validators[validatorID]
    return exists
}

-> 위에서 업그레이드한 트렌젝션
package chain

import (
    "blockchain/pkg/block"
    "blockchain/pkg/transaction"
    "fmt"
    "time"
)

type Blockchain struct {
    Blocks         []*block.Block
    Validators     map[string]int
    Transactions   []transaction.Transaction
    BlockGasLimit  int // 블록 가스 한도
    TransactionGas int // 트랜잭션당 가스 소비량
    MiningInterval time.Duration // 블록 생성 주기
}

func NewBlockchain(blockGasLimit, transactionGas int, miningInterval time.Duration) *Blockchain {
    bc := &Blockchain{
        Blocks:         []*block.Block{createGenesisBlock()},
        Validators:     make(map[string]int),
        BlockGasLimit:  blockGasLimit,
        TransactionGas: transactionGas,
        MiningInterval: miningInterval,
    }
    
    go bc.startMining() // 백그라운드에서 블록 생성 시작
    return bc
}

func createGenesisBlock() *block.Block {
    return block.NewBlock(0, "Genesis Block", "")
}

func (bc *Blockchain) AddTransaction(tx transaction.Transaction, gas int) {
    bc.Transactions = append(bc.Transactions, tx)
}

func (bc *Blockchain) startMining() {
    for {
        time.Sleep(bc.MiningInterval)
        bc.CreateBlockFromTransactions("validator1")
    }
}

func (bc *Blockchain) CreateBlockFromTransactions(validatorID string) {
    if len(bc.Transactions) == 0 {
        return
    }

    // 가스 한도에 맞게 트랜잭션 선택
    var selectedTxs []transaction.Transaction
    var totalGas int
    for _, tx := range bc.Transactions {
        if totalGas + bc.TransactionGas > bc.BlockGasLimit {
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
    bc.Transactions = bc.Transactions[len(selectedTxs):] // 선택한 트랜잭션을 큐에서 제거
}

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

func (bc *Blockchain) isValidNewBlock(newBlock, previousBlock *block.Block) bool {
    // 블록의 인덱스 검증
    if previousBlock.Index+1 != newBlock.Index {
        return false
    }

    // 블록의 해시 검증
    if newBlock.PreviousHash != previousBlock.Hash {
        return false
    }

    // 블록의 해시가 올바르게 계산되었는지 검증
    if newBlock.Hash != newBlock.calculateHash() {
        return false
    }

    return true
}

func (bc *Blockchain) isValidValidator(validatorID string) bool {
    _, exists := bc.Validators[validatorID]
    return exists
}