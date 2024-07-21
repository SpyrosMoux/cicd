package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"spyrosmoux/api/internal/helpers"
	"time"
)

var (
	ghAppClientId   = helpers.LoadEnvVariable("GITHUB_APP_CLIENT_ID")
	ghAppPrivateKey = helpers.LoadEnvVariable("GITHUB_APP_PRIVATE_KEY_PATH")
)

func GenerateJWT() string {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(ghAppPrivateKey))
	if err != nil {
		log.Fatal("ERROR: Could not parse private key with error: ", err)
	}

	claims := jwt.MapClaims{
		"client_id": ghAppClientId,
		"exp":       time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
		"iat":       time.Now().Unix(),                    // Issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(key)
	if err != nil {
		log.Fatal("ERROR: Could not generate token with error: ", err)
	}

	return signedToken
}
