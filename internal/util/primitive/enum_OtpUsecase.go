package primitive

import "time"

type OtpUseCase string

const (
	OtpUseCaseVerifyEmail OtpUseCase = "verify-email"
)

func (o OtpUseCase) GetTTL() time.Duration {
	switch o {
	case OtpUseCaseVerifyEmail:
		return 5 * time.Minute
	}

	return 0
}

func (o OtpUseCase) GetLimitRetry() int16 {
	switch o {
	case OtpUseCaseVerifyEmail:
		return 5
	}

	return 0
}
