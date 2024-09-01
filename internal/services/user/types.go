package user

import (
	primitive2 "github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
)

type RegisterUserInput struct {
	BackgroundImage *primitive2.PresignedFileUpload
	ImageProfile    *primitive2.PresignedFileUpload
	Password        string
	Email           string
	FullName        string
	RegisterAs      primitive2.EnumRegisterAs
}

type RegisterUserOutput struct {
	UserID                            int64
	BackgroundImagePresignedUrlUpload *primitive2.PresignedFileUploadOutput
	ImageProfilePresignedUrlUpload    *primitive2.PresignedFileUploadOutput
}
