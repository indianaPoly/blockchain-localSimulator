package transaction

import "fmt"

// Transaction 인터페이스 정의
type Transaction interface {
    Process() error
}

// 계약 함수 실행 트랜잭션
type ContractExecutionTransaction struct {
    ContractID string
    Function   string
    Params     []interface{}
}

func (tx *ContractExecutionTransaction) Process() error {
    // 계약 함수 실행 로직
    fmt.Printf("Executing contract ID: %s, Function: %s with params: %v\n", tx.ContractID, tx.Function, tx.Params)
    // 예를 들어, 계약의 함수 호출 로직 구현
    return nil
}

// 송금 트랜잭션
type TransferTransaction struct {
    From   string
    To     string
    Amount float64
}

func (tx *TransferTransaction) Process() error {
    // 송금 로직
    fmt.Printf("Transferring %f from %s to %s\n", tx.Amount, tx.From, tx.To)
    // 예를 들어, 계좌의 잔액을 업데이트하는 로직 구현
    return nil
}
