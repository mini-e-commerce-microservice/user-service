package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mini-e-commerce-microservice/user-service/generated/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"strings"
)

func getMsg(msg []string, code int) string {
	if msg != nil && len(msg) > 0 {
		return strings.Join(msg, ". ")
	}

	return defaultStatusCodeMessages[code]
}

func Error(w http.ResponseWriter, r *http.Request, code int, err error, msg ...string) {
	ctx, span := otel.Tracer("error").Start(r.Context(), fmt.Sprintf("error code: %d", code))
	defer span.End()
	span.SetAttributes(attribute.String("error.server", err.Error()))

	r = r.WithContext(ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errMsgByte := make([]byte, 0)

	switch code {
	case http.StatusInternalServerError:
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		errMsg := api.Error{
			Message: getMsg(msg, code),
		}
		errMsgByte, err = json.Marshal(errMsg)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			w.Write([]byte(`{"error": "internal server error"}`))
			return
		}
	default:
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errMsg := api.Error400{
				Errors: make(map[string][]string),
			}

			for _, validationError := range validationErrors {
				fieldName := validationError.Field()
				errMsg.Errors[fieldName] = []string{
					validationError.Error(),
				}
			}

			errMsgByte, err = json.Marshal(errMsg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				w.Write([]byte(`{"error": "internal server error"}`))
				return
			}
		} else {
			errMsg := api.Error{
				Message: getMsg(msg, code),
			}
			errMsgByte, err = json.Marshal(errMsg)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				w.Write([]byte(`{"error": "internal server error"}`))
				return
			}
		}
	}

	span.SetAttributes(attribute.String("error.response", string(errMsgByte)))
	w.Write(errMsgByte)
}

var defaultStatusCodeMessages = map[int]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusMethodNotAllowed:    "Method Not Allowed",
	http.StatusConflict:            "Conflict",
	http.StatusInternalServerError: "Internal Status Error",
	http.StatusUnprocessableEntity: "Unprocessable Entity",
}
