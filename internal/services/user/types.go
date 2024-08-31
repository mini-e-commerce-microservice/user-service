package user

import (
	"github.com/mini-e-commerce-microservice/user-service/internal/primitive"
)

type RegisterUserInput struct {
	BackgroundImage *primitive.PresignedFileUpload
	ImageProfile    *primitive.PresignedFileUpload
	Password        string
	Email           string
	FullName        string
	RegisterAs      primitive.EnumRegisterAs
}

type RegisterUserOutput struct {
	UserID                            int64
	BackgroundImagePresignedUrlUpload *primitive.PresignedFileUploadOutput
	ImageProfilePresignedUrlUpload    *primitive.PresignedFileUploadOutput
}
