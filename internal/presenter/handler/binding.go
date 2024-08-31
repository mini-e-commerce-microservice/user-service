package handler

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"github.com/mini-e-commerce-microservice/user-service/internal/primitive"
	"io"
	"net/http"
)

func (h *handler) bodyRequestBindToStruct(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		Error(w, r, http.StatusUnprocessableEntity, err)
		return false
	}
	defer r.Body.Close()

	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		Error(w, r, http.StatusUnprocessableEntity, err, err.Error())
		return false
	}

	err = h.validator.Struct(v)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return false
	}
	return true
}

func (h *handler) bindUploadFileRequest(w http.ResponseWriter, r *http.Request, input api.FileUploadRequest) (output primitive.PresignedFileUpload, ok bool) {
	fileUploadInput := primitive.NewPresignedFileUploadInput{
		Identifier:       input.Identifier,
		OriginalFileName: input.OriginalFilename,
		MimeType:         primitive.MimeType(input.MimeType),
		Size:             input.Size,
		ChecksumSHA256:   input.ChecksumSha256,
	}
	output, err := primitive.NewPresignedFileUpload(fileUploadInput)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err, err.Error())
		return output, false
	}
	return output, true
}

func (h *handler) bindUploadFileResponse(input primitive.PresignedFileUploadOutput) (output api.FileUploadResponse) {
	return api.FileUploadResponse{
		Identifier:      input.Identifier,
		UploadExpiredAt: input.UploadExpiredAt,
		UploadUrl:       input.UploadURL,
		MinioFormData:   input.MinioFormData,
	}
}

func (h *handler) writeJson(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	respByte, err := json.Marshal(v)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respByte)
}

var decoder = schema.NewDecoder()

func (h *handler) queryParamBindToStruct(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := r.ParseForm(); err != nil {
		Error(w, r, http.StatusBadRequest, err, err.Error())
		return false
	}

	decoder.SetAliasTag("json")
	if err := decoder.Decode(v, r.Form); err != nil {
		Error(w, r, http.StatusBadRequest, err, err.Error())
		return false
	}

	err := h.validator.Struct(v)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return false
	}
	return true
}
