package gh

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spyrosmoux/api/internal/helpers"

	"github.com/golang-jwt/jwt/v5"
)

/*
INFO: https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation
*/

type AccessTokenResponse struct {
	AccessToken         string     `json:"token"`
	ExpiresAt           string     `json:"expires_at"`
	Permissions         Permission `json:"permissions"`
	RepositorySelection string     `json:"repository_selection"`
}

type Permission struct {
	Contents        string `json:"contents"`
	Metadata        string `json:"metadata"`
	PullRequests    string `json:"pull_requests"`
	RepositoryHooks string `json:"repository_hooks"`
	Statuses        string `json:"statuses"`
}

var (
	ghAppClientId   = helpers.LoadEnvVariable("GITHUB_APP_CLIENT_ID")
	ghAppPrivateKey = helpers.LoadEnvVariable("GITHUB_APP_PRIVATE_KEY_PATH")
)

func GetInstallationToken(installationId int64) (string, error) {
	token := GenerateJWT()

	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationId)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse AccessTokenResponse
	err = json.Unmarshal(body, &tokenResponse)

	return tokenResponse.AccessToken, nil
}

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
