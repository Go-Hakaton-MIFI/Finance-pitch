package middleware

import (
	"context"
	"encoding/json"
	delivery "finance-backend/internal/delivery/http"
	"finance-backend/internal/domain"
	"finance-backend/pkg/logger"
	"finance-backend/pkg/utils"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTParserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		publicKey := strings.ReplaceAll(os.Getenv("JWT_PUBLIC"), `\n`, "\n")

		log := logger.NewLogger()

		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "no auth token in request"})
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid token format"})
			return
		}

		jwtString := strings.Split(tokenString, "Bearer ")[1]

		token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				log.Error(r.Context(), "error on parsing jwt token", nil)
				return nil, delivery.ErrJwtTokenParsing
			}

			parsedSecter, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))

			if err != nil {
				return nil, err
			}

			return parsedSecter, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "error on parsing jwt token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			var ctx context.Context
			if claims["is_admin"] != "" {
				userLogin, _ := claims["sub"].(string)
				isAdmin, _ := claims["is_admin"].(bool)
				user := domain.User{
					Login:   userLogin,
					IsAdmin: isAdmin,
				}
				ctx = context.WithValue(r.Context(), utils.ContextKeyUser, user)
			} else {
				ctx = context.WithValue(r.Context(), utils.ContextKeyUser, nil)
			}
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			log.Error(r.Context(), "error on getting jwt token claims", nil)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "error on parsing jwt token"})
			return
		}
	})
}
