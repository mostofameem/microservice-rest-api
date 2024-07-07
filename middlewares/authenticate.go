package middlewares

import (
	"fmt"
	"net/http"
	"order_service/config"
	"order_service/web/utils"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var userKey = "user"

type AuthClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func unauthorizedResponse(w http.ResponseWriter) {
	utils.SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
}
func VerifyToken(tokenStr string) (AuthClaims, error) {
	conf := config.GetConfig()
	var claims = AuthClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(conf.JwtSecret), nil
		},
	)

	if !token.Valid {
		err = fmt.Errorf("unauthorized")
	}
	return claims, err
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// collect token from header
		header := r.Header.Get("authorization")
		tokenStr := ""

		// collect token from query
		if len(header) == 0 {
			tokenStr = r.URL.Query().Get("auth")
		} else {
			tokens := strings.Split(header, " ")
			if len(tokens) != 2 {
				unauthorizedResponse(w)
				return
			}
			tokenStr = tokens[1]
		}
		claims, err := VerifyToken(tokenStr)

		// set user id in the context
		if err != nil {
			unauthorizedResponse(w)
			return
		}

		r.Header.Set("id", strconv.Itoa(claims.Id))
		r.Header.Set("email", claims.Email)
		next.ServeHTTP(w, r)
	})
}

func GetUserId(r *http.Request) (int, error) {
	userIdVal := r.Context().Value(userKey)
	userId, ok := userIdVal.(int)
	if !ok {
		return 0, fmt.Errorf("unauthorized")
	}
	return userId, nil
}
