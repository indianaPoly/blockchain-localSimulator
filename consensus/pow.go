package consensus

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

type ProofOfWork struct {
    Difficulty int
}

func NewProofOfWork(difficulty int) *ProofOfWork {
    return &ProofOfWork{Difficulty: difficulty}
}

// 블록에 대한 작업 증명(해시 계산)
func (pow *ProofOfWork) Run(data string) (string, int) {
    var nonce int
    var hash [32]byte

    for {
        hash = sha256.Sum256([]byte(data + strconv.Itoa(nonce)))
        hashStr := fmt.Sprintf("%x", hash)
        if strings.HasPrefix(hashStr, strings.Repeat("0", pow.Difficulty)) {
            break
        }
        nonce++
    }

    return fmt.Sprintf("%x", hash), nonce
}
