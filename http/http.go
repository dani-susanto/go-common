package http

import (
	"fmt"
	gohttp "net/http"

	"github.com/dani-susanto/go-common/jwt"
	"github.com/dani-susanto/go-common/log"
	"github.com/redis/go-redis/v9"
)

type HTTP interface {
	GetMux() *gohttp.ServeMux

	// HTTP Method
	Get(path string, handler gohttp.HandlerFunc, middlewares ...Middleware)
	Post(path string, handler gohttp.HandlerFunc, middlewares ...Middleware)
	Patch(path string, handler gohttp.HandlerFunc, middlewares ...Middleware)

	// Middleware
	Log(next gohttp.Handler) gohttp.Handler
	Authenticated(next gohttp.Handler) gohttp.Handler
	Role(roles ...string) Middleware
	RateLimiter(limit int, duration string) Middleware
	GlobalRateLimiter(next gohttp.Handler) gohttp.Handler
	ExcludeGlobalRateLimiter(method string, path string) Middleware
}

type http struct {
	mux                       *gohttp.ServeMux
	cache                     *redis.Client
	globalRateLimiterExcluded map[string]bool
	globalRateLimiterAttempt  int
	globalRateLimiterDuration string
	accessTokenJWT            jwt.JWT
	responder                 Responder
	log                       log.Log
}

func New(
	mux *gohttp.ServeMux,
	cache *redis.Client,
	globalRateLimiterAttempt int,
	globalRateLimiterDuration string,
	accessTokenJWT jwt.JWT,
	responder Responder,
	log log.Log,
) HTTP {
	return &http{
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

func (r *http) Get(
	path string,
	handler gohttp.HandlerFunc,
	middlewares ...Middleware,
) {
	r.Route("GET", path, handler, middlewares...)
}

func (r *http) Post(
	path string,
	handler gohttp.HandlerFunc,
	middlewares ...Middleware,
) {
	r.Route("POST", path, handler, middlewares...)
}

func (r *http) Patch(
	path string,
	handler gohttp.HandlerFunc,
	middlewares ...Middleware,
) {
	r.Route("PATCH", path, handler, middlewares...)
}

func (r *http) Route(
	method string,
	path string,
	handler gohttp.HandlerFunc,
	middlewares ...Middleware,
) {
	middleware := r.MiddlewareStack(middlewares...)
	r.mux.Handle(fmt.Sprintf("%s %s", method, path), middleware(handler))
}

func (r *http) GetMux() *gohttp.ServeMux {
	return r.mux
}
