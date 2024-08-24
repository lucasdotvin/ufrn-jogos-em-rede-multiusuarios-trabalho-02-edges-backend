package auth

import (
	"net/http"
	"trabalho-02-edges/lib/presentation/auth"
	"trabalho-02-edges/lib/presentation/auth/token"
)

type CookieTokenAuthMiddleware struct {
	tokenService token.Service
}

func NewCookieTokenAuthMiddleware(tokenService token.Service) *CookieTokenAuthMiddleware {
	return &CookieTokenAuthMiddleware{
		tokenService,
	}
}

func (m *CookieTokenAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenContent, err := r.Cookie(token.AccessTokenKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		at, err := m.tokenService.ParseTokenFromContent(tokenContent.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !at.IsValid() {
			http.Error(w, token.InvalidTokenError.Error(), http.StatusUnauthorized)
			return
		}

		if !at.IsAccessToken() {
			http.Error(w, token.InvalidTokenError.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = auth.WithUserUuid(ctx, at.Uid)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
