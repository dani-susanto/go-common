package rest

import (
	"fmt"
	"net/http"

	"github.com/dani-susanto/go-common/jwt"
	"github.com/dani-susanto/go-common/log"
	"github.com/redis/go-redis/v9"
)

type Rest interface {
	GetMux() *http.ServeMux

	// HTTP Method
	Get(path string, handler http.HandlerFunc, middlewares ...Middleware)
	Post(path string, handler http.HandlerFunc, middlewares ...Middleware)
	Patch(path string, handler http.HandlerFunc, middlewares ...Middleware)

	// Middleware
	Log(next http.Handler) http.Handler
	Authenticated(next http.Handler) http.Handler
	Role(roles ...string) Middleware
	RateLimiter(limit int, duration string) Middleware
	GlobalRateLimiter(next http.Handler) http.Handler
	ExcludeGlobalRateLimiter(method string, path string) Middleware
}

type rest struct {
	mux                       *http.ServeMux
	cache                     *redis.Client
	globalRateLimiterExcluded map[string]bool
	globalRateLimiterAttempt  int
	globalRateLimiterDuration string
	accessTokenJWT            jwt.JWT
	responder                 Responder
	log                       log.Log
}

func New(
	mux *http.ServeMux,
	cache *redis.Client,
	globalRateLimiterAttempt int,
	globalRateLimiterDuration string,
	accessTokenJWT jwt.JWT,
	responder Responder,
	log log.Log,
) Rest {
	return &rest{
		mux:                       mux,
		cache:                     cache,
		globalRateLimiterExcluded: make(map[string]bool),
		globalRateLimiterAttempt:  globalRateLimiterAttempt,
		globalRateLimiterDuration: globalRateLimiterDuration,
		accessTokenJWT:            accessTokenJWT,
		responder:                 responder,
		log:                       log,
	}
}

func (r *rest) Get(
	path string,
	handler http.HandlerFunc,
	middlewares ...Middleware,
) {
	r.Route("GET", path, handler, middlewares...)
}

func (r *rest) Post(
	path string,
	handler http.HandlerFunc,
	middlewares ...Middleware,
) {
	r.Route("POST", path, handler, middlewares...)
}

func (r *rest) Patch(
	path string,
	handler http.HandlerFunc,
	middlewares ...Middleware,
) {
	r.Route("PATCH", path, handler, middlewares...)
}

func (r *rest) Route(
	method string,
	path string,
	handler http.HandlerFunc,
	middlewares ...Middleware,
) {
	middleware := r.MiddlewareStack(middlewares...)
	r.mux.Handle(fmt.Sprintf("%s %s", method, path), middleware(handler))
}

func (r *rest) GetMux() *http.ServeMux {
	return r.mux
}
