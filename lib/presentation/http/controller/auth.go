package controller

import (
	"encoding/json"
	"net/http"
	"trabalho-02-edges/lib/application/case/user"
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/presentation/auth"
	"trabalho-02-edges/lib/presentation/auth/token"
	userrequest "trabalho-02-edges/lib/presentation/request/user"
	userresponse "trabalho-02-edges/lib/presentation/response/user"
)

const (
	refreshRoutePath   = "/api/v1/auth/refresh"
	globalApiRoutePath = "/api/v1"
	rootPath           = "/"
)

type AuthController struct {
	signUpUseCase   *user.SignUpUseCase
	signInUseCase   *user.SignInUseCase
	findUserUseCase *user.FindUserUseCase
	tokenService    token.Service
}

func NewAuthController(
	signUpUseCase *user.SignUpUseCase,
	signInUseCase *user.SignInUseCase,
	findUserUseCase *user.FindUserUseCase,
	tokenService token.Service,
) *AuthController {
	return &AuthController{
		signUpUseCase,
		signInUseCase,
		findUserUseCase,
		tokenService,
	}
}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var input userrequest.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := c.signUpUseCase.Execute(input.Name, input.Username, input.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.authenticate(u, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parsedResponse := userresponse.NewUserResponse(u, nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(parsedResponse)
}

func (c *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	var input userrequest.SignInRequest
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := c.signInUseCase.Execute(input.Username, input.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.authenticate(u, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parsedResponse := userresponse.NewUserResponse(u, nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(parsedResponse)
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	rt, err := c.getRefreshTokenFromCookie(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !rt.IsValid() {
		http.Error(w, token.InvalidTokenError.Error(), http.StatusBadRequest)
		return
	}

	u, err := c.findUserUseCase.Execute(rt.Uid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.authenticate(u, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	userUuid, ok := auth.GetUserUuid(r.Context())

	if !ok {
		http.Error(w, entity.UserNotFoundError.Error(), http.StatusUnauthorized)
		return
	}

	u, err := c.findUserUseCase.Execute(userUuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedResponse := userresponse.NewUserResponse(u, nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(parsedResponse)
}

func (c *AuthController) authenticate(u *entity.User, w http.ResponseWriter) error {
	at, err := c.tokenService.GenerateAccessToken(u)

	if err != nil {
		return err
	}

	rt, err := c.tokenService.GenerateRefreshToken(u)

	if err != nil {
		return err
	}

	bt, err := c.tokenService.GenerateBroadcastToken(u)

	if err != nil {
		return err
	}

	c.setAccessTokenCookies(w, at)
	c.setRefreshTokenCookies(w, rt)
	c.setBroadcastTokenCookies(w, bt)

	return nil
}

func (c *AuthController) setAccessTokenCookies(w http.ResponseWriter, at *token.Token) {
	http.SetCookie(w, &http.Cookie{
		Name:     token.AccessTokenKey,
		Value:    at.Content,
		Expires:  at.ExpiresAt,
		HttpOnly: true,
		Path:     globalApiRoutePath,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    token.AccessTokenKey + "_duration",
		Expires: at.ExpiresAt,
		Path:    rootPath,
	})
}

func (c *AuthController) setRefreshTokenCookies(w http.ResponseWriter, rt *token.Token) {
	http.SetCookie(w, &http.Cookie{
		Name:     token.RefreshTokenKey,
		Value:    rt.Content,
		Expires:  rt.ExpiresAt,
		HttpOnly: true,
		Path:     refreshRoutePath,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    token.RefreshTokenKey + "_duration",
		Expires: rt.ExpiresAt,
		Path:    rootPath,
	})
}

func (c *AuthController) setBroadcastTokenCookies(w http.ResponseWriter, bt *token.Token) {
	http.SetCookie(w, &http.Cookie{
		Name:    token.BroadcastTokenKey,
		Value:   bt.Content,
		Expires: bt.ExpiresAt,
		Path:    rootPath,
	})
}

func (c *AuthController) getRefreshTokenFromCookie(r *http.Request) (*token.Token, error) {
	cookie, err := r.Cookie(token.RefreshTokenKey)

	if err != nil {
		return nil, err
	}

	return c.tokenService.ParseTokenFromContent(cookie.Value)
}
