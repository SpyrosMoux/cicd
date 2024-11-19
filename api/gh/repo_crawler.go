package gh

import (
	"context"
	"errors"
	"fmt"
	"github.com/spyrosmoux/cicd/runner/pipelines"
	"io"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// FetchPipelineConfig scans a given repo for valid pipeline yamls in the '.flowforge' directory and returns an array with
// all the valid yamls.
// TODO(@SpyrosMoux) rewrite this so it makes more sense. Goal is to fetch all pipelines that need to run (FetchTriggeredPipelines???)
func FetchPipelineConfig(repoOwner string, repoName string, branchName string, installationId int64) ([]pipelines.Pipeline, error) {
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

	var validYAMLs []pipelines.Pipeline

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

		pipeline, err := pipelines.ValidateYAMLStructure([]byte(fileContent))
		if err != nil {
			// TODO(@SpyrosMoux) what happens if a yaml is invalid, and there are multiple yamls in the repo?
			// skipping invalid yaml for now
			fmt.Println(err.Error())
			continue // skip to next pipeline
		}

		validYAMLs = append(validYAMLs, pipeline)
	}

	return validYAMLs, nil
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
