package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/internal/dto"
	"github.com/Yandex-Practicum/final/internal/infrastructure/env"
	httpService "github.com/Yandex-Practicum/final/internal/services/http"
	jwtService "github.com/Yandex-Practicum/final/pkg/jwt"
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

	if data.Password != env.EnvList.Password {
		httpService.WriteJson(w, dto.ErrorResponse{ErrorText: "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	token, err := jwtService.CreateToken(data.Password)
	if err != nil {
		log.Println("failed to sign jwt: %s\n" + err.Error())
		httpService.WriteJson(w, dto.ErrorResponse{ErrorText: "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   token,
		Expires: time.Now().Add(tokenExpired * time.Hour),
	}
	http.SetCookie(w, cookie)

	fmt.Println(token)
	httpService.WriteJsonOKResponse(w, SuccessResponse{Token: token})
}
