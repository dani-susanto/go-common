package rest

import (
	"context"
	"fmt"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/dani-susanto/go-common/jwt"
)

type Middleware func(gohttp.Handler) gohttp.Handler

func (r *rest) MiddlewareStack(xs ...Middleware) Middleware {
	return func(next gohttp.Handler) gohttp.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

type responseWriter struct {
	gohttp.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (s *rest) Log(next gohttp.Handler) gohttp.Handler {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		ctx := r.Context()
		start := time.Now()

		s.log.InfoContext(ctx, "INCOMING",
			"method", r.Method,
			"path", r.URL.Path,
			"time", start.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			"ip", r.RemoteAddr,
		)

		wrapped := &responseWriter{ResponseWriter: w, status: gohttp.StatusOK}
		next.ServeHTTP(wrapped, r)

		args := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"time", time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			"status", wrapped.status,
			"latency", time.Since(start),
		}

		if wrapped.status >= 400 {
			s.log.ErrorContext(ctx, "OUTGOING", args...)
		} else {
			s.log.InfoContext(ctx, "OUTGOING", args...)
		}
	})
}

func (s *rest) Authenticated(next gohttp.Handler) gohttp.Handler {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		authorization := r.Header.Get("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			s.responder.Error(w, gohttp.StatusUnauthorized, "unauthorized", nil)
			return
		}

		encodedToken := strings.TrimPrefix(authorization, "Bearer ")

		claims, err := s.accessTokenJWT.Validate(encodedToken)
		if err != nil {
			s.responder.Error(w, gohttp.StatusUnauthorized, err.Error(), nil)
			return
		}

		ctx := context.WithValue(r.Context(), jwt.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *rest) Role(roles ...string) Middleware {
	return func(next gohttp.Handler) gohttp.Handler {
		return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
			claims, ok := jwt.GetClaims(r.Context())
			if !ok {
				s.responder.Error(w, gohttp.StatusUnauthorized, "unauthorized", nil)
				return
			}

			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			s.responder.Error(w, gohttp.StatusForbidden, "forbidden", nil)
		})
	}
}

func (s *rest) ExecRateLimiter(
	ctx context.Context,
	limit int,
	duration string,
	method string,
	path string,
	IP string,
) (int, string) {
	key := fmt.Sprintf("http:limiter:ip:%s:%s:%s", IP, method, path)
	userID, ok := jwt.GetUserID(ctx)
	if ok {
		key = fmt.Sprintf("http:limiter:user_id:%d:%s:%s", userID, method, path)
	}

	timeDuration, err := time.ParseDuration(duration)
	if err != nil {
		return gohttp.StatusInternalServerError, err.Error()
	}

	count, err := s.cache.Incr(ctx, key).Result()
	if err != nil {
		return gohttp.StatusInternalServerError, err.Error()
	}

	if count == 1 {
		s.cache.Expire(ctx, key, timeDuration)
	}

	if count > int64(limit) {
		return gohttp.StatusTooManyRequests, "too many requests, please try again later."
	}

	return gohttp.StatusOK, ""
}

func (s *rest) GlobalRateLimiter(next gohttp.Handler) gohttp.Handler {
	return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		ctx := r.Context()

		if s.globalRateLimiterExcluded[fmt.Sprintf("%s:%s", r.Method, r.URL.Path)] {
			next.ServeHTTP(w, r)
			return
		}

		err, msg := s.ExecRateLimiter(
			ctx,
			s.globalRateLimiterAttempt,
			s.globalRateLimiterDuration,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		if err != gohttp.StatusOK {
			s.responder.Error(w, err, msg, nil)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (s *rest) ExcludeGlobalRateLimiter(method string, path string) Middleware {
	s.globalRateLimiterExcluded[fmt.Sprintf("%s:%s", method, path)] = true
	return func(next gohttp.Handler) gohttp.Handler {
		return next
	}

}

func (s *rest) RateLimiter(limit int, duration string) Middleware {
	return func(next gohttp.Handler) gohttp.Handler {
		return gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
			ctx := r.Context()

			err, msg := s.ExecRateLimiter(
				ctx,
				limit,
				duration,
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			)

			if err != gohttp.StatusOK {
				s.responder.Error(w, err, msg, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
