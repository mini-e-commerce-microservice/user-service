package user

import (
	"github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
)

type RegisterUserInput struct {
	BackgroundImage *primitive.PresignedFileUpload
	ImageProfile    *primitive.PresignedFileUpload
	Password        string
	Email           string
	FullName        string
}

type RegisterUserOutput struct {
	UserID                            int64
	BackgroundImagePresignedUrlUpload *primitive.PresignedFileUploadOutput
	ImageProfilePresignedUrlUpload    *primitive.PresignedFileUploadOutput
}

type ActivationEmailUserInput struct {
	Token string
}
