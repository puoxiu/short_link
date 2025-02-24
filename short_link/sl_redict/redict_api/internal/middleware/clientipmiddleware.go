package middleware

import "net/http"

type ClientIPMiddleware struct {
}

func NewClientIPMiddleware() *ClientIPMiddleware {
	return &ClientIPMiddleware{}
}

func (m *ClientIPMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
