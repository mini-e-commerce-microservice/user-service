package user

import (
	"context"
	"errors"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	s3wrapper "github.com/SyaibanAhmadRamadhan/go-s3-wrapper"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/profiles"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	"golang.org/x/sync/errgroup"
	"time"
)

func (s *service) GetUser(ctx context.Context, input GetUserInput) (output GetUserOutput, err error) {
	userOutput, err := s.userRepository.FindOneUser(ctx, users.FindOneUserInput{
		ID: null.IntFrom(input.UserID),
	})
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			err = ErrUserNotFound
		}
		return output, collection.Err(err)
	}

	profileOutput, err := s.profileRepository.FindOneProfile(ctx, profiles.FindOneProfileInput{
		UserID: null.IntFrom(input.UserID),
	})
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			err = ErrUserNotFound
		}
		return output, collection.Err(err)
	}

	output = GetUserOutput{
		ID:              userOutput.Data.ID,
		Email:           userOutput.Data.Email,
		FullName:        profileOutput.Data.FullName,
		IsEmailVerified: userOutput.Data.IsEmailVerified,
	}

	var eg errgroup.Group

	if profileOutput.Data.ImageProfile != nil {
		eg.Go(func() (err error) {
			imageProfile, err := s.s3.PresignedUrlGetObject(ctx, s3wrapper.PresignedUrlGetObjectInput{
				ObjectName: *profileOutput.Data.ImageProfile,
				BucketName: s.bucketName,
				Expired:    10 * time.Minute,
			})
			if err != nil {
				return collection.Err(err)
			}
			output.ImageProfile = null.StringFrom(imageProfile.URL)
			return
		})
	}
	if profileOutput.Data.BackgroundImage != nil {
		eg.Go(func() (err error) {
			backgroundImage, err := s.s3.PresignedUrlGetObject(ctx, s3wrapper.PresignedUrlGetObjectInput{
				ObjectName: *profileOutput.Data.BackgroundImage,
				BucketName: s.bucketName,
				Expired:    10 * time.Minute,
			})
			if err != nil {
				return collection.Err(err)
			}
			output.BackgroundImage = null.StringFrom(backgroundImage.URL)
			return
		})
	}

	if err = eg.Wait(); err != nil {
		return output, collection.Err(err)
	}
	return
}

type GetUserInput struct {
	UserID int64
}

type GetUserOutput struct {
	ID              int64
	Email           string
	FullName        string
	IsEmailVerified bool
	ImageProfile    null.String
	BackgroundImage null.String
}
