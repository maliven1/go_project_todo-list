package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maliven1/go_final_project/entity"
)

var secret = []byte("dsa53219nlxvnju")

var TODO_PASSWORD = "1324657980"

var claims = jwt.MapClaims{
	"password": TODO_PASSWORD,
}

func SignToken() string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		fmt.Printf("failed to sign jwt: %s\n", err)
	}
	return signedToken
}

func AuthorizationGenerateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var UserPass entity.UserPass
	var TokenJson entity.TokenJson
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		responseWhithError(w, "Ошибка чтения")
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &UserPass); err != nil {
		responseWhithError(w, "Ошибка чтения")
		return
	}

	if UserPass.Password == TODO_PASSWORD {
		TokenJson.Token = SignToken()
		responseWithConfirmPas(w, TokenJson)
		return
	}
	http.Error(w, "Не верный пароль", http.StatusUnauthorized)
}
