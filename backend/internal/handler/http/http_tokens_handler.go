package http

import (
	"encoding/json"
	"fmt"
	"github.com/chempik1234/availability-checker-web/internal/ports/tokens"
	"github.com/chempik1234/availability-checker-web/internal/types"
	"net/http"
	"strings"
)

type TokensHttpHandler struct {
	tokensRepository tokens.TokensRepository
}

func NewTokensHttpHandler(tokensRepository tokens.TokensRepository) *TokensHttpHandler {
	return &TokensHttpHandler{tokensRepository: tokensRepository}
}

func (h *TokensHttpHandler) CheckTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaderValues := strings.Split(r.Header.Get(types.AuthorizationHeader), " ")
		if len(authHeaderValues) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if authHeaderValues[0] != types.TokenPrefix {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(types.TokenPrefixErrorMessage))
			return
		}

		token := authHeaderValues[1]
		ok, err := h.tokensRepository.Check(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Error while checking token: %v", err.Error())))
			return
		}

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r)
	})
}

func (h *TokensHttpHandler) NewTokensHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			token, err := h.tokensRepository.Create()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Errorf("Failed to save response: %v", err.Error()).Error()))
				break
			}

			result := make(map[string]string)
			result["token"] = token

			resultJson, err := json.Marshal(result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Errorf("Failed to save response: %v", err.Error()).Error()))
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(resultJson)
		case http.MethodDelete:
			tokenToDelete := r.URL.Query().Get("token")
			if tokenToDelete == "" {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("query argument token must be provided"))
				break
			}
			err := h.tokensRepository.Delete(tokenToDelete)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Errorf("Failed to delete: %v", err.Error()).Error()))
			}

			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet, http.MethodPut, http.MethodPatch:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
