package handler

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/jwt_claims_proto"
	jwt_util "github.com/mini-e-commerce-microservice/user-service/internal/util/jwt"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
	"net/http"
	"strings"
)

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
		h.httpOtel.Err(w, r, http.StatusBadRequest, err, err.Error())
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

func (h *handler) getUserFromBearerAuth(w http.ResponseWriter, r *http.Request) (*jwt_claims_proto.JwtAuthAccessTokenClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, collection.Err(errors.New("authorization header is missing")))
		return nil, false
	}

	bearerSplit := strings.Split(authHeader, " ")
	if len(bearerSplit) != 2 {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, collection.Err(errors.New("invalid authorization header format")))
		return nil, false
	}

	if bearerSplit[0] != "Bearer" {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, collection.Err(errors.New("authorization scheme must be Bearer")))
		return nil, false
	}

	authAccessTokenJWT := &jwt_util.AuthAccessTokenClaims{}
	err := authAccessTokenJWT.ClaimsHS256(bearerSplit[1], h.jwtAccessTokenConf.Key)
	if err != nil {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, collection.Err(err))
		return nil, false
	}

	return authAccessTokenJWT.JwtAuthAccessTokenClaims, true
}
