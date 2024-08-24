package auth

import (
	"net/http"
	"trabalho-02-edges/lib/presentation/auth"
	"trabalho-02-edges/lib/presentation/auth/token"
)

type WebSocketQueryTokenAuthMiddleware struct {
	tokenService token.Service
}

func NewWebSocketQueryTokenAuthMiddleware(tokenService token.Service) *WebSocketQueryTokenAuthMiddleware {
	return &WebSocketQueryTokenAuthMiddleware{
		tokenService,
	}
}

func (m *WebSocketQueryTokenAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenContent, ok := r.URL.Query()[token.BroadcastTokenKey]

		if !ok {
			http.Error(w, token.InvalidTokenError.Error(), http.StatusUnauthorized)
			return
		}

		if len(tokenContent) != 1 {
			http.Error(w, token.InvalidTokenError.Error(), http.StatusUnauthorized)
			return
		}

		bt, err := m.tokenService.ParseTokenFromContent(tokenContent[0])

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !bt.IsValid() {
			http.Error(w, token.InvalidTokenError.Error(), http.StatusUnauthorized)
			return
		}

		if !bt.IsBroadcastToken() {
			http.Error(w, token.InvalidTokenError.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = auth.WithUserUuid(ctx, bt.Uid)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
