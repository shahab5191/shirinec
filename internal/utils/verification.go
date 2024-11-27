package utils

import (
	"math/rand"
)

func GenerateVerificationCode() int {
    return rand.Intn(900000) + 100000
}
