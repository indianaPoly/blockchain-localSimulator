package main

import (
	"fmt"
	"strings"
)

type Contract struct {
	ID    string
	State string
	Code  string // 계약의 코드, 예: "합계: %d", 10
}

// 계약을 생성하는 함수
func NewContract(id, state, code string) *Contract {
	return &Contract{
		ID:    id,
		State: state,
		Code:  code,
	}
}

// 계약을 배포하는 함수
func (bc *Blockchain) DeployContract(contract *Contract) {
	bc.AddBlock(fmt.Sprintf("Depoly Contract: %s", contract.ID))
}

// 계약을 실행하는 메소드
func (bc *Blockchain) ExecuteContract(id string) {
	for _, block := range bc.Blocks {
		if strings.Contains(block.Data, id) {
			// 계약 실행 로직을 구현
			fmt.Printf("Executing contract ID: %s\n", id)
			// 예를 들어, 계약의 코드 실행
		}
	}
}

// 계약의 상태 관리
func (contract *Contract) UpdateState(newState string) {
	contract.State = newState
}

// 계약을 찾는 함수
func (bc *Blockchain) GetContract(id string) *Contract {
	for _, block := range bc.Blocks {
		if strings.Contains(block.Data, id) {
			// 계약을 찾으면 반환
			return &Contract{ID: id, Code: "Example Code"}
		}
	}
	return nil
}

// 계약 코드 실행
func ExecuteCode(code string) string {
	// 코드 실행 로직 (예: 간단한 문자열 처리)
	result := fmt.Sprintf("Executing code: %s", code)
	return result
}
