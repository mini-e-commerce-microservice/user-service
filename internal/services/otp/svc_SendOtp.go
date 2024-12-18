package otp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/notification_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/model"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/otps"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *service) SendOtp(ctx context.Context, input SendOtpInput) (err error) {
	userOutput, err := s.validateExistingUser(ctx, input.Type, input.DestinationAddress)
	if err != nil {
		return collection.Err(err)
	}

	err = s.dbTx.DoTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false},
		func(tx wsqlx.Rdbms) (err error) {
			err = s.otpRepository.DeleteOtp(ctx, otps.DeleteOtpInput{
				Tx:         tx,
				UserID:     null.IntFrom(userOutput.Data.ID),
				Usecase:    null.StringFrom(string(input.Usecase)),
				Type:       null.StringFrom(string(input.Type)),
				TokenIsNil: true,
			})
			if err != nil {
				return collection.Err(err)
			}

			expiredOtp := time.Now().UTC().Add(input.Usecase.GetTTL())
			otpCode := util.GenerateOTP()

			_, err = s.otpRepository.CreateOtp(ctx, otps.CreateOtpInput{
				Tx: tx,
				Payload: model.Otp{
					UserID:  userOutput.Data.ID,
					Usecase: string(input.Usecase),
					Code:    fmt.Sprintf("%d", otpCode),
					Type:    string(input.Type),
					Counter: 0,
					Expired: expiredOtp,
				},
			})
			if err != nil {
				return collection.Err(err)
			}

			err = s.rabbitmqRepository.Publish(ctx, rabbitmq.PublishInput{
				RoutingKey: s.rabbitMQConf.Exchanges.NotificationExchange.NotificationTypeEmail.RoutingKey.Name,
				Exchange:   s.rabbitMQConf.Exchanges.NotificationExchange.Name,
				Payload: &notification_proto.Notification{
					Type: notification_proto.NotificationType_ACTIVATION_EMAIL,
					Data: &notification_proto.Notification_ActivationEmail{
						ActivationEmail: &notification_proto.NotificationActivationEmailPayload{
							OtpCode:        fmt.Sprintf("%d", otpCode),
							RecipientEmail: userOutput.Data.Email,
							RecipientName:  userOutput.Data.Email,
							ExpiredAt:      timestamppb.New(expiredOtp),
						},
					},
				},
			})
			if err != nil {
				return collection.Err(err)
			}

			return
		},
	)
	if err != nil {
		return collection.Err(err)
	}

	return
}

func (s *service) validateExistingUser(ctx context.Context, otpType primitive.OtpType, destinationAddr string) (userOutput users.FindOneUserOutput, err error) {
	if otpType == primitive.OtpTypeEmail {
		userOutput, err = s.userRepository.FindOneUser(ctx, users.FindOneUserInput{
			Email: null.StringFrom(destinationAddr),
		})
		if err != nil {
			if errors.Is(err, repositories.ErrRecordNotFound) {
				err = errors.Join(err, fmt.Errorf("%w: %s", ErrDestinationAddressNotFound, destinationAddr))
				return userOutput, collection.Err(err)
			}
			return userOutput, collection.Err(err)
		}
		if userOutput.Data.IsEmailVerified {
			return userOutput, collection.Err(ErrEmailUserIsVerified)
		}
	}

	return
}
