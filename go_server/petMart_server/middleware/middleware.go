package middleware

import (
	"context"
	"net/http"
)

func InjectHttpMiddleware() func(http.Handler) http.Handler{
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				CookieWriter := CookieAccess{Writer: w}
				contextKey := HTTPkey("PetMart")
				ctx := context.WithValue(r.Context(), contextKey, &CookieWriter)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			})
		}
}
