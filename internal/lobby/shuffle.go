package lobby

import (
	"crypto/rand"
	"math/big"
)

// shuffleIndices returns a Fisher-Yates shuffle of [0, n).
func shuffleIndices(n int) ([]int, error) {
	if n == 0 {
		return nil, nil
	}

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

	for i := n - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return nil, err
		}
		j := int(jBig.Int64())
		order[i], order[j] = order[j], order[i]
	}

	return order, nil
}
