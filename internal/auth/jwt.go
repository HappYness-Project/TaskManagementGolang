package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
)

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		tokenString = tokenString[len("Bearer "):]
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId, _ := strconv.Atoi(claims["nameid"].(string))
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userId)
		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.AccessToken), nil
	})
}
func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}
