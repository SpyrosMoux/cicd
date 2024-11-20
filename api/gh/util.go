package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/spyrosmoux/cicd/common/helpers"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ghAppClientId   = helpers.LoadEnvVariable("GITHUB_APP_CLIENT_ID")
	ghAppPrivateKey = helpers.LoadEnvVariable("GITHUB_APP_PRIVATE_KEY_PATH")
)

// getInstallationToken Uses an installationId and a generated JWT token to get an access token
func getInstallationToken(installationId int64) (string, error) {
	token := generateJWT()

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

// generateJWT generates a JWT token from a given GitHub App private key and pem file.
// TODO(@SpyrosMoux) should return error instead of fatalling
func generateJWT() string {
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

// downloadYAMLContent downloads the content of a given raw GitHub url
func downloadYAMLContent(downloadUrl string, installationId int64) ([]byte, error) {
	token, err := getInstallationToken(installationId)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pipeline config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pipeline config: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// A range of events could trigger a pipeline. i.e push a new branch, make a commit
// For now we only support push events
// These events should be matched as follows
// - If a commit is made to the specified branch -> run
// - If a branch is created, and specified in the triggers -> run
// - If a tag is create, and the tag is specified in the trigges -> run
// Apart from creating stuff the push event also represents deletion events. Such as deleting a tag or branch.
// These events should be ignored
func matchPushEventWithBranch(event *github.PushEvent, branches []string) bool {
	branchName, err := getBranchNameFromRef(event.GetRef())
	if err != nil {
		log.Printf("Failed to get branch name from ref %s: %v", event.GetRef(), err)
		return false
	}

	shouldRun := false
	for _, branch := range branches {
		if branch == "*" {
			shouldRun = true
			break
		}

		if branchName == branch {
			fmt.Printf("Matching push event for branch %s\n", branch)
			shouldRun = true
		}
	}

	return shouldRun
}

func getBranchNameFromRef(ref string) (string, error) {
	if !strings.Contains(ref, "refs/heads/") {
		return "", errors.New("ref does not contain refs/heads/")
	}

	parts := strings.SplitAfter(ref, "refs/heads/")
	if len(parts) != 2 {
		return "", errors.New("error getting branch name from ref " + ref)
	}

	return parts[1], nil
}
