package handler

import (
	"errors"
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"github.com/mini-e-commerce-microservice/user-service/internal/primitive"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/user"
	"net/http"
)

func (h *handler) V1RegisterPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1RegisterPostRequestBody{}

	if !h.bodyRequestBindToStruct(w, r, &req) {
		return
	}

	registerUserInput := user.RegisterUserInput{
		BackgroundImage: nil,
		ImageProfile:    nil,
		Password:        req.Password,
		Email:           req.Email,
		FullName:        req.FullName,
		RegisterAs:      primitive.EnumRegisterAs(req.RegisterAs),
	}

	if req.BackgroundImage != nil {
		backgroundImageFileUpload, ok := h.bindUploadFileRequest(w, r, *req.BackgroundImage)
		if !ok {
			return
		}

		registerUserInput.BackgroundImage = &backgroundImageFileUpload
	}
	if req.ImageProfile != nil {
		imageProfileFileUpload, ok := h.bindUploadFileRequest(w, r, *req.ImageProfile)
		if !ok {
			return
		}

		registerUserInput.ImageProfile = &imageProfileFileUpload
	}

	registerOutput, err := h.userService.RegisterUser(r.Context(), registerUserInput)
	if err != nil {
		if errors.Is(err, user.ErrEmailAvailable) {
			Error(w, r, http.StatusBadRequest, err, user.ErrEmailAvailable.Error())
		} else {
			Error(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	resp := api.V1RegisterPost200Response{
		Id: registerOutput.UserID,
	}
	if registerOutput.BackgroundImagePresignedUrlUpload != nil {
		resp.UploadBackgroundImage = &api.FileUploadResponse{
			Identifier:      registerOutput.BackgroundImagePresignedUrlUpload.Identifier,
			MinioFormData:   registerOutput.BackgroundImagePresignedUrlUpload.MinioFormData,
			UploadExpiredAt: registerOutput.BackgroundImagePresignedUrlUpload.UploadExpiredAt,
			UploadUrl:       registerOutput.BackgroundImagePresignedUrlUpload.UploadURL,
		}
	}
	if registerOutput.ImageProfilePresignedUrlUpload != nil {
		resp.UploadProfileImage = &api.FileUploadResponse{
			Identifier:      registerOutput.ImageProfilePresignedUrlUpload.Identifier,
			MinioFormData:   registerOutput.ImageProfilePresignedUrlUpload.MinioFormData,
			UploadExpiredAt: registerOutput.ImageProfilePresignedUrlUpload.UploadExpiredAt,
			UploadUrl:       registerOutput.ImageProfilePresignedUrlUpload.UploadURL,
		}
	}

	h.writeJson(w, r, http.StatusOK, resp)
}
