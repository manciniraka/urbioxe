package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

func GenerateEmployeeNumber(
	departmentCode string,
	joinDate time.Time,
	sequence int,
) string {

	return fmt.Sprintf(
		"%s-%d-%04d",
		strings.ToUpper(departmentCode),
		joinDate.Year(),
		sequence,
	)
}

func GenerateTemporaryPassword(departmentCode string) (string, error) {
	randomNumber, err := rand.Int(
		rand.Reader,
		big.NewInt(10000),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s@%04d",
		strings.ToLower(departmentCode),
		randomNumber.Int64(),
	), nil
}