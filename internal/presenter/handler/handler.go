package handler

import (
	httplogwrap "github.com/SyaibanAhmadRamadhan/http-log-wrap"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
)

type handler struct {
	r         *chi.Mux
	validator *validator.Validate
}

func NewHandler(r *chi.Mux) {
	r.Use(httplogwrap.HttpOtel)
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	h := &handler{
		r:         r,
		validator: nil,
	}
	h.route()
}

func (h *handler) route() {
	h.r.Get("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello World"))
	})
}
