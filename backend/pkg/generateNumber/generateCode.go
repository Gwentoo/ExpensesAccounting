package generatenumber

import (
	"crypto/rand"
	"math/big"
)

func GenerateFiveDigitNumber() (int, error) {
	num := big.NewInt(90000)

	n, err := rand.Int(rand.Reader, num)
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + 10000, nil
}
