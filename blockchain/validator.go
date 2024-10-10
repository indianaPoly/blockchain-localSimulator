package blockchain

import (
	"math/rand"
	"time"
)

type Validator struct {
	Stake int
	Count int
}

// PoS 알고리즘을 활용하여 Validator를 알아냄
func (bc *Blockchain) SelectValidator() string {
	totalStake := 0

	// 전체 stake를 계산
	for _, validator := range bc.Validators {
		// validator들 중에서 stake가 조금이라도 있는 사람들한테 기본적인 가중치를 제공하여
		// validator가 될 수 있게끔
		if validator.Stake > 0 {
			totalStake += validator.Stake
		}
	}

	// 로컬 랜덤 생성
	rng := rand.New(rand.NewSource(time.Now().UnixMicro()))
	randomFloat := rng.Float64()

	// 랜덤 값에 따른 검증자 선택
	cumulativeProbability := 0.0

	for validatorID, validator := range bc.Validators {
		// 각 검증자의 확률 비율 계산
		probability := float64(validator.Stake) / float64(totalStake)
		cumulativeProbability += probability

		// 랜덤 값이 해당 검증자의 확률 구간에 속하면 해당 검증자를 선택
		if randomFloat <= cumulativeProbability {
			validator.Count++
			return validatorID
		}
	}

	return ""
}

// 검증자가 유효한지 확인
func (bc *Blockchain) isValidValidator(validatorID string) bool {
	_, exists := bc.Validators[validatorID]
	return exists
}
