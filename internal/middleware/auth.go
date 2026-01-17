package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/Yandex-Practicum/final/internal/infrastructure/env"
	jwtService "github.com/Yandex-Practicum/final/pkg/jwt"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			// получаем куку
			cookie, err := r.Cookie("token")
			if err != nil {
				log.Println("failed to sign jwt: %s\n" + err.Error())
				http.Error(w, "Ошибка аутентификации", http.StatusBadRequest)
				return
			}

			jwt := cookie.Value

			token, err := jwtService.CreateToken(env.EnvList.Password)
			if err != nil {
				log.Println("failed to sign jwt: %s\n" + err.Error())
				http.Error(w, "Ошибка аутентификации", http.StatusBadRequest)
				return
			}

			if token != jwt {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
