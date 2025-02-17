package util

import (
	"math/rand"
	"strconv"
)

func GenerateAccountNumber() string {
	var accountNumber string

	for i := 0; i < 16; i++ {
		if i > 0 && i%4 == 0 {
			accountNumber += "-"
		}
		accountNumber += strconv.Itoa(rand.Intn(10))
	}

	return accountNumber
}
