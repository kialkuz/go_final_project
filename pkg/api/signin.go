package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/env"
	jwtService "github.com/Yandex-Practicum/final/pkg/services/jwt"
	"github.com/golang-jwt/jwt/v5"
)

type UserData struct {
	Password string `json:"password"`
}

type SuccessResponse struct {
	Token string `json:"token"`
}

const (
	tokenExpired = 8
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data UserData

	if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data.Password != env.Lookup("TODO_PASSWORD", "") {
		writeJson(w, dto.ErrorResponse{ErrorText: "Неверный пароль"})
		return
	}

	token, err := jwtService.CreateToken(data.Password)
	if err != nil {
		log.Println("failed to sign jwt: %s\n" + err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "Неверный пароль"})
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   token,
		Expires: time.Now().Add(tokenExpired * time.Hour),
	}
	http.SetCookie(w, cookie)

	fmt.Println(token)
	writeJson(w, SuccessResponse{Token: token})
}

func createToken(password string) (string, error) {
	secret := []byte(password)

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
