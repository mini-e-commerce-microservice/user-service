package handler

import (
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	primitive2 "github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
	"net/http"
	"strconv"
)

func (h *handler) bindUploadFileRequest(w http.ResponseWriter, r *http.Request, input api.FileUploadRequest) (output primitive2.PresignedFileUpload, ok bool) {
	fileUploadInput := primitive2.NewPresignedFileUploadInput{
		Identifier:       input.Identifier,
		OriginalFileName: input.OriginalFilename,
		MimeType:         primitive2.MimeType(input.MimeType),
		Size:             input.Size,
		ChecksumSHA256:   input.ChecksumSha256,
	}
	output, err := primitive2.NewPresignedFileUpload(fileUploadInput)
	if err != nil {
		h.httpOtel.Err(w, r, http.StatusBadRequest, err, err.Error())
		return output, false
	}
	return output, true
}

func (h *handler) bindUploadFileResponse(input primitive2.PresignedFileUploadOutput) (output api.FileUploadResponse) {
	return api.FileUploadResponse{
		Identifier:      input.Identifier,
		UploadExpiredAt: input.UploadExpiredAt,
		UploadUrl:       input.UploadURL,
		MinioFormData:   input.MinioFormData,
	}
}

func (h *handler) getUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	userIDStr := r.Header.Get("X-User-Id")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.httpOtel.Err(w, r, http.StatusForbidden, err)
		return 0, false
	}

	return userID, true
}
