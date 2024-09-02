package middlewares

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("dsa53219nlxvnju")

var TODO_PASSWORD = "1324657980"

var claims = jwt.MapClaims{
	"password": TODO_PASSWORD,
}

func CheckToken(token string) bool {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		fmt.Printf("failed to sign jwt: %s\n", err)
	}
	return token == signedToken
}
func NewAuthMeddlewares() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// смотрим наличие пароля
			if len(TODO_PASSWORD) > 0 {
				var jwtToken string // JWT-токен из куки
				// получаем куку
				cookie, err := r.Cookie("token")
				if err == nil {
					jwtToken = cookie.Value
				}
				if err != nil {
					fmt.Printf("failed to parse token: %s\n", err)
					return
				}

				var valid bool
				if CheckToken(jwtToken) {
					valid = true
				}

				if !valid {
					// возвращаем ошибку авторизации 401
					http.Error(w, "Authentification required", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
		return http.HandlerFunc(fn)
	}

}
