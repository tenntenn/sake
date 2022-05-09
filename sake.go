package sake

import (
	"context"
	"net/http"
)

type Handler[T any, P Param[T]] interface {
	ServeHTTP(ctx context.Context, w http.ResponseWriter, r *Request[T, P]) error
}

type HandlerFunc[T any, P Param[T]] func(ctx context.Context, w http.ResponseWriter, r *Request[T, P]) error

func (f HandlerFunc[T, P]) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *Request[T, P]) error {
	return f(ctx, w, r)
}

func _[T any, P Param[T]]() {
	var _ Handler[T, P] = HandlerFunc[T, P](nil)
}

type Request[T any, P Param[T]] struct {
	*http.Request
	Param P
}

type Param[T any] interface {
	*T
	Set(r *http.Request) error
}

type ErrorHandler func(w http.ResponseWriter, err error)

func DefaultErrorHandler(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

func Standard[T any, P Param[T]](h Handler[T, P], errhandler ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		pr := &Request[T, P]{
			Request: r,
			Param:   P(new(T)),
		}
		pr.Param.Set(r)
		if err := h.ServeHTTP(ctx, w, pr); err != nil {
			if errhandler != nil {
				errhandler(w, err)
			} else {
				DefaultErrorHandler(w, err)
			}
		}
	})
}

func StandardFunc[T any, P Param[T]](h HandlerFunc[T, P], errhandler ErrorHandler) http.Handler {
	return Standard[T, P](h, errhandler)
}
