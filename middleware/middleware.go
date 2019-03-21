package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

const (
	ACCESS_TOKEN_KEY = "Access-Token"
)

type goMiddleware struct {
	// another stuff , may be needed by middleware
}

type responseError struct {
	Message string `json:"message"`
}

func (m *goMiddleware) CORS(next http.Handler) http.Handler {
	return cors.Default().Handler(next)
}

func InitMiddleware() *goMiddleware {
	return &goMiddleware{}
}
