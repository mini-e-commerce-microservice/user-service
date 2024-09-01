package otp

import (
	"context"
	"database/sql"
	"fmt"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/notification_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/model"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/otps"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *service) SendOtp(ctx context.Context, input SendOtpInput) (err error) {
	err = s.validateSendOtp(ctx, input)
	if err != nil {
		return tracer.Error(err)
	}

	err = s.dbTx.DoTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false},
		func(tx wsqlx.Rdbms) (err error) {
			err = s.otpRepository.DeleteOtp(ctx, otps.DeleteOtpInput{
				Tx:         tx,
				UserID:     null.IntFrom(input.UserID),
				Usecase:    null.StringFrom(string(input.Usecase)),
				Type:       null.StringFrom(string(input.Type)),
				ExpiredGTE: null.TimeFrom(time.Now().UTC()),
			})
			if err != nil {
				return tracer.Error(err)
			}

			expiredOtp := time.Now().UTC().Add(input.Usecase.GetTTL())
			otpCode := util.GenerateOTP()

			_, err = s.otpRepository.CreateOtp(ctx, otps.CreateOtpInput{
				Tx: tx,
				Payload: model.Otp{
					UserID:  input.UserID,
					Usecase: string(input.Usecase),
					Code:    fmt.Sprintf("%d", otpCode),
					Type:    string(input.Type),
					Counter: 0,
					Expired: expiredOtp,
				},
			})
			if err != nil {
				return tracer.Error(err)
			}

			err = s.rabbitmqRepository.Publish(ctx, rabbitmq.PublishInput{
				RoutingKey: rabbitmq.RoutingKeyEmailOTP,
				Exchange:   rabbitmq.ExchangeNameNotification,
				Payload: &notification_proto.Notification{
					Type: notification_proto.NotificationType_EMAIL_VERIFIED,
					Data: &notification_proto.Notification_EmailVerified{
						EmailVerified: &notification_proto.NotificationEmailVerifiedPayload{
							OtpCode:   fmt.Sprintf("%d", otpCode),
							ExpiredAt: timestamppb.New(expiredOtp),
						},
					},
				},
			})
			if err != nil {
				return tracer.Error(err)
			}

			return
		},
	)
	if err != nil {
		return tracer.Error(err)
	}

	return
}

func (s *service) validateSendOtp(ctx context.Context, input SendOtpInput) (err error) {
	if input.Type == primitive.OtpTypeEmail {
		exists, err := s.userRepository.CheckExistingUser(ctx, users.CheckExistingUserInput{
			Email: null.StringFrom(input.DestinationAddress),
		})
		if err != nil {
			return tracer.Error(err)
		}

		if !exists {
			return tracer.Error(fmt.Errorf("%w: %s", ErrDestinationAddressNotFound, input.DestinationAddress))
		}
	}

	return
}
