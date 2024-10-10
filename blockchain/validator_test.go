package blockchain

import "testing"

// 초기 블록체인을 설정
func setupBlockchain() *Blockchain {
	return &Blockchain{
		//
		Validators: map[string]*Validator{
			"validator1":  {Stake: 60, Count: 0},
			"validator2":  {Stake: 30, Count: 0},
			"validator3":  {Stake: 10, Count: 0},
			"validator4":  {Stake: 50, Count: 0},
			"validator5":  {Stake: 20, Count: 0},
			"validator6":  {Stake: 40, Count: 0},
			"validator7":  {Stake: 25, Count: 0},
			"validator8":  {Stake: 15, Count: 0},
			"validator9":  {Stake: 35, Count: 0},
			"validator10": {Stake: 5, Count: 0},
		},
	}
}

// 선택 횟수가 예상 범위 내에 있는지 확인하는 헬퍼 함수
func validateSelectionCount(t *testing.T, validator string, count int, expectedPercentage float64, totalSelections int) {
	expectedCount := expectedPercentage * float64(totalSelections)
	minCount := int(expectedCount * 0.9) // 예상 횟수의 90% ~ 110% 범위로 설정
	maxCount := int(expectedCount * 1.1)

	if count < minCount || count > maxCount {
		t.Errorf("%s selection count is out of expected range: %d (expected between %d and %d)", validator, count, minCount, maxCount)
	} else {
		t.Logf("%s selected %d times, which is within the expected range (%d - %d)", validator, count, minCount, maxCount)
	}
}

// go test -run TestSelectValidator -v : 로그까지 확인할 수 있는 터미널 명령어
func TestSelectValidator(t *testing.T) {
	bc := setupBlockchain()

	// 선택된 검증자의 횟수를 저장하는 맵
	selectionCount := map[string]int{
		"validator1":  0,
		"validator2":  0,
		"validator3":  0,
		"validator4":  0,
		"validator5":  0,
		"validator6":  0,
		"validator7":  0,
		"validator8":  0,
		"validator9":  0,
		"validator10": 0,
	}

	totalSelections := 10000
	for i := 0; i < totalSelections; i++ {
		selectionValidator := bc.SelectValidator()
		selectionCount[selectionValidator]++
	}

	totalStake := 0
	for _, validator := range bc.Validators {
		totalStake += validator.Stake
	}

	for validatorID, validator := range bc.Validators {
		expectedPercentage := float64(validator.Stake) / float64(totalStake)
		validateSelectionCount(t, validatorID, selectionCount[validatorID], expectedPercentage, totalSelections)
	}
}

// go test -run TestIsValidValidator -v : 검증자가 존재하는지 테스트하는 코드
func TestIsValidValidator(t *testing.T) {
	bc := setupBlockchain()

	validatorNames := [2]string{"validator1", "validator11"}

	for _, name := range validatorNames {
		if bc.isValidValidator(name) {
			t.Logf("%s가 validator에 존재합니다.", name)
		} else {
			t.Logf("%s는 validator에 존재하지 않습니다.", name)
		}
	}
}
