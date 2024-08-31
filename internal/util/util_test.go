package util_test

import (
	"fmt"
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	fmt.Println(util.GenerateOTP())
}
