package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"spyrosmoux/api/internal/helpers"
	"time"
)

/*
INFO: https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app
*/

var (
	ghAppClientId   = helpers.LoadEnvVariable("GITHUB_APP_CLIENT_ID")
	ghAppPrivateKey = helpers.LoadEnvVariable("GITHUB_APP_PRIVATE_KEY_PATH")
)

func GenerateJWT() string {
	pemFileData, err := os.ReadFile(ghAppPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(pemFileData)
	if err != nil {
		log.Fatal("ERROR: Could not parse private key with error: ", err)
	}

	claims := jwt.MapClaims{
		"iat": time.Now().Unix() - 60,  // Issued at, 60 seconds in the past to allow for clock drift
		"exp": time.Now().Unix() + 600, // Token expires in 10 minutes
		"iss": ghAppClientId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(key)
	if err != nil {
		log.Fatal("ERROR: Could not generate token with error: ", err)
	}

	return signedToken
}
