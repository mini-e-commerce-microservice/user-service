package presenter

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type StatusError int

const (
	InternalServerError = StatusError(http.StatusInternalServerError)
)

type errorDetail struct {
	Err             string `json:"err"`
	StatusCode      int    `json:"status_code"`
	LogRequestBody  bool   `json:"log_request_body"`
	LogResponseBody bool   `json:"log_response_body"`
	LogParams       bool   `json:"log_params"`
}

func (e *errorDetail) Error() string {
	return ""
}

func newError(w http.ResponseWriter, r *http.Request, statusCode StatusError, err error, opts ...Option) {
	ed := &errorDetail{
		Err:             err.Error(),
		StatusCode:      int(statusCode),
		LogParams:       true,
		LogRequestBody:  true,
		LogResponseBody: true,
	}

	//for _, opt := range opts {
	//	opt(ed)
	//}

	//_, span := observability.TracerUserService.Start(r.Context(), "response body", trace.WithAttributes(
	//	attribute.String("asd", "asd"),
	//))
	//span.End()
	w.WriteHeader(int(statusCode))
	b, err := json.Marshal(ed)
	if err != nil {
		log.Err(err)
	}

	w.Write(b)
}
