package blockchain

import "math/rand"

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

// 검증자가 유효한지 확인
func (bc *Blockchain) isValidValidator(validatorID string) bool {
	_, exists := bc.Validators[validatorID]
	return exists
}
