package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dani-susanto/go-common/jwt"
)

type Middleware func(http.Handler) http.Handler

func (r *rest) MiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (s *rest) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.log.Info("INCOMING method=%s path=%s time=%s ip=%s",
			r.Method,
			r.URL.Path,
			start.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			r.RemoteAddr,
		)

		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		msg := fmt.Sprintf("OUTGOING method=%s path=%s time=%s status=%d latency=%s",
			r.Method,
			r.URL.Path,
			time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			wrapped.status,
			time.Since(start),
		)

		switch {
		case wrapped.status >= 500:
			s.log.Error("%v", msg)
		case wrapped.status >= 400:
			s.log.Error("%v", msg)
		default:
			s.log.Info("%v", msg)
		}

	})
}

func (s *rest) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			s.responder.Error(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}

		encodedToken := strings.TrimPrefix(authorization, "Bearer ")

		claims, err := s.accessTokenJWT.Validate(encodedToken)
		if err != nil {
			s.responder.Error(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		ctx := context.WithValue(r.Context(), jwt.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *rest) Role(roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := jwt.GetClaims(r.Context())
			if !ok {
				s.responder.Error(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}

			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			s.responder.Error(w, http.StatusForbidden, "forbidden", nil)
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
	key := fmt.Sprintf("rest:limiter:ip:%s:%s:%s", IP, method, path)
	userID, ok := jwt.GetUserID(ctx)
	if ok {
		key = fmt.Sprintf("rest:limiter:user_id:%d:%s:%s", userID, method, path)
	}

	timeDuration, err := time.ParseDuration(duration)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	count, err := s.cache.Incr(ctx, key).Result()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	if count == 1 {
		s.cache.Expire(ctx, key, timeDuration)
	}

	if count > int64(limit) {
		return http.StatusTooManyRequests, "too many requests, please try again later."
	}

	return http.StatusOK, ""
}

func (s *rest) GlobalRateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		if err != http.StatusOK {
			s.responder.Error(w, err, msg, nil)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (s *rest) ExcludeGlobalRateLimiter(method string, path string) Middleware {
	s.globalRateLimiterExcluded[fmt.Sprintf("%s:%s", method, path)] = true
	return func(next http.Handler) http.Handler {
		return next
	}

}

func (s *rest) RateLimiter(limit int, duration string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			err, msg := s.ExecRateLimiter(
				ctx,
				limit,
				duration,
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			)

			if err != http.StatusOK {
				s.responder.Error(w, err, msg, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
