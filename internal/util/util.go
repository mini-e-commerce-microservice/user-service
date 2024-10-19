package util

import (
	"context"
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
	"math/rand"
	"time"
)

func GenerateOTP() int {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return newRand.Intn(900000) + 100000
}

func GetTraceParent(ctx context.Context) *string {
	traceParent := whttp.GetTraceParent(ctx)
	if traceParent != "" {
		return &traceParent
	}

	return nil
}
