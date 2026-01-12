package jwt

import "github.com/golang-jwt/jwt/v5"

func CreateToken(password string) (string, error) {
	secret := []byte(password)

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
