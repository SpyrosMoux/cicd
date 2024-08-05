package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spyrosmoux/api/internal/auth"
)

/*
INFO: https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation
*/

type accessTokenResponse struct {
	AccessToken         string     `json:"token"`
	ExpiresAt           string     `json:"expires_at"`
	Permissions         permission `json:"permissions"`
	RepositorySelection string     `json:"repository_selection"`
}

type permission struct {
	Contents        string `json:"contents"`
	Metadata        string `json:"metadata"`
	PullRequests    string `json:"pull_requests"`
	RepositoryHooks string `json:"repository_hooks"`
	Statuses        string `json:"statuses"`
}

func GetInstallationToken(installationId int64) (string, error) {
	token := auth.GenerateJWT()

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

	var tokenResponse accessTokenResponse
	err = json.Unmarshal(body, &tokenResponse)

	return tokenResponse.AccessToken, nil
}
