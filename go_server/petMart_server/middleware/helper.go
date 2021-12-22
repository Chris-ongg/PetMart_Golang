package middleware

import (
	"context"
	"net/http"
)

type HTTPkey string

type CookieAccess struct {
	Writer http.ResponseWriter
}

func (C *CookieAccess) SetToken(token string) {
	http.SetCookie(C.Writer , &http.Cookie{
		Name:       "PET_MART_USER",
		Value:      token,
		HttpOnly:   false,
	})
}

func GetCookieAccess(ctx context.Context) *CookieAccess {
	contextKey := HTTPkey("PetMart")
	return ctx.Value(contextKey).(*CookieAccess)
}