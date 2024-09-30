package gh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/spyrosmoux/core-engine/pkg/models"
	"golang.org/x/oauth2"
)

// FetchPipelineConfig scans a given repo for valid pipeline yamls in the '.flowforge' directory and returns an array with
// all the valid yamls.
func FetchPipelineConfig(repoOwner string, repoName string, branchName string, installationId int64) ([]string, error) {
	token, err := GetInstallationToken(installationId)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		TokenType:   "Bearer",
		AccessToken: token,
	})
	tokenClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(tokenClient)
	options := &github.RepositoryContentGetOptions{
		Ref: branchName,
	}

	_, contents, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, ".flowforge", options) // Only look in .flowforge dir for pipelines
	if err != nil {
		return nil, err
	}

	var validYAMLs []string

	for _, file := range contents {
		fmt.Printf("Found file: %s\n", file.GetName())

		downloadURL := file.GetDownloadURL()
		if downloadURL == "" {
			return nil, errors.New("Could not get download URL for pipeline " + file.GetName())
		}
		fileContent, err := downloadYAMLContent(downloadURL, installationId)
		if err != nil {
			return nil, err
		}

		_, err = models.ValidateYAMLStructure([]byte(fileContent))
		if err != nil {
			// TODO(spyrosmoux) what happens if a yaml is invalid, and there are multiple yamls in the repo?
			// skip for now
			fmt.Println(err.Error())
			continue // skip to next pipeline
		}

		// TODO(spyrosmoux) check if yaml has trigger of type same as eventType

		validYAMLs = append(validYAMLs, string(fileContent))
	}

	return validYAMLs, nil
}

func validatePipelineTrigger(trigger string) error {
	//TODO(spyrosmoux) implement me
	panic("implement me")
}

// downloadYAMLContent downloads the content of a given raw GitHub url
func downloadYAMLContent(downloadUrl string, installationId int64) ([]byte, error) {
	token, err := GetInstallationToken(installationId)
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
