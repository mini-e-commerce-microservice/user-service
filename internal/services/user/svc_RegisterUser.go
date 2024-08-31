package user

import (
	"context"
	"database/sql"
	"fmt"
	s3wrapper "github.com/SyaibanAhmadRamadhan/go-s3-wrapper"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/notification_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/model"
	"github.com/mini-e-commerce-microservice/user-service/internal/primitive"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/profiles"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// RegisterUser
// return error ErrEmailAvailable
func (s *service) RegisterUser(ctx context.Context, input RegisterUserInput) (output RegisterUserOutput, err error) {
	err = s.validateRegisterUser(ctx, input)
	if err != nil {
		return output, tracer.Error(err)
	}

	passwordHash := make([]byte, 0)
	var backgroundImageFileName *string
	var imageProfileFileName *string

	var eg errgroup.Group

	eg.Go(func() (err error) {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return tracer.Error(err)
		}
		return
	})

	if input.BackgroundImage != nil {
		eg.Go(func() (err error) {
			presignedOutput, err := s.s3.PresignedUrlUploadObject(ctx, s3wrapper.PresignedUrlUploadObjectInput{
				BucketName: s.bucketName,
				Path:       input.BackgroundImage.GeneratedFileName,
				MimeType:   string(input.BackgroundImage.MimeType),
				Checksum:   input.BackgroundImage.ChecksumSHA256,
				Expired:    5 * time.Minute,
			})
			if err != nil {
				return tracer.Error(err)
			}

			output.BackgroundImagePresignedUrlUpload = &primitive.PresignedFileUploadOutput{
				Identifier:      input.BackgroundImage.Identifier,
				UploadURL:       presignedOutput.URL,
				UploadExpiredAt: presignedOutput.ExpiredAt,
				MinioFormData:   presignedOutput.MinioFormData,
			}
			backgroundImageFileName = &input.BackgroundImage.GeneratedFileName
			return
		})
	}

	if input.ImageProfile != nil {
		eg.Go(func() (err error) {
			presignedOutput, err := s.s3.PresignedUrlUploadObject(ctx, s3wrapper.PresignedUrlUploadObjectInput{
				BucketName: s.bucketName,
				Path:       input.ImageProfile.GeneratedFileName,
				MimeType:   string(input.ImageProfile.MimeType),
				Checksum:   input.ImageProfile.ChecksumSHA256,
				Expired:    5 * time.Minute,
			})
			if err != nil {
				return tracer.Error(err)
			}

			output.ImageProfilePresignedUrlUpload = &primitive.PresignedFileUploadOutput{
				Identifier:      input.ImageProfile.Identifier,
				UploadURL:       presignedOutput.URL,
				UploadExpiredAt: presignedOutput.ExpiredAt,
				MinioFormData:   presignedOutput.MinioFormData,
			}
			imageProfileFileName = &input.ImageProfile.GeneratedFileName
			return
		})
	}

	if err = eg.Wait(); err != nil {
		return output, tracer.Error(err)
	}

	err = s.dbTx.DoTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false},
		func(tx wsqlx.Rdbms) (err error) {
			userCreateOutput, err := s.userRepository.CreateUser(ctx, users.CreateUserInput{
				Tx: tx,
				Payload: model.User{
					Email:           input.Email,
					Password:        string(passwordHash),
					IsEmailVerified: false,
				},
			})
			if err != nil {
				return tracer.Error(err)
			}

			_, err = s.profileRepository.CreateProfile(ctx, profiles.CreateProfileInput{
				Tx: tx,
				Payload: model.Profile{
					UserID:          userCreateOutput.ID,
					FullName:        input.FullName,
					ImageProfile:    imageProfileFileName,
					BackgroundImage: backgroundImageFileName,
				},
			})
			if err != nil {
				return tracer.Error(err)
			}

			output.UserID = userCreateOutput.ID

			return
		},
	)
	if err != nil {
		return output, tracer.Error(err)
	}

	err = s.rabbitmqRepository.Publish(ctx, rabbitmq.PublishInput{
		RoutingKey: rabbitmq.RoutingKeyEmailOTP,
		Exchange:   rabbitmq.ExchangeNameNotification,
		Payload: &notification_proto.Notification{
			Type: notification_proto.NotificationType_EMAIL_VERIFIED,
			Data: &notification_proto.Notification_EmailVerified{
				EmailVerified: &notification_proto.NotificationEmailVerifiedPayload{
					OtpCode:   fmt.Sprintf("%d", util.GenerateOTP()),
					ExpiredAt: timestamppb.New(time.Now().UTC().Add(6 * time.Minute)),
				},
			},
		},
	})
	if err != nil {
		return output, tracer.Error(err)
	}

	return
}

func (s *service) validateRegisterUser(ctx context.Context, input RegisterUserInput) (err error) {
	exists, err := s.userRepository.CheckExistingUser(ctx, users.CheckExistingUserInput{
		ID:    null.Int{},
		Email: null.StringFrom(input.Email),
	})
	if err != nil {
		return tracer.Error(err)
	}

	if exists {
		return tracer.Error(ErrEmailAvailable)
	}
	return
}
