package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
