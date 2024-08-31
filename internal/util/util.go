package util

import (
	"math/rand"
	"time"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func GenerateOTP() int {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return newRand.Intn(900000) + 100000
}
