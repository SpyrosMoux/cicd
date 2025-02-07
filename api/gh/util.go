package gh

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v68/github"
)

// downloadYAMLContent downloads the content of a given raw GitHub url
func downloadYAMLContent(downloadUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	req.Header.Set("Authorization", "Bearer "+GhToken)

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
		return false
	}

	shouldRun := false
	for _, branch := range branches {
		if branch == "*" {
			shouldRun = true
			break
		}

		if branchName == branch {
			shouldRun = true
		}
	}

	return shouldRun
}

// matchPullRequestEventWithBranch matches the base branch (where the PR points to) with the list
// of branches mentioned in the pr: section of the yaml. If there is match return true else false
func matchPullRequestEventWithBranch(event *github.PullRequestEvent, branches []string) bool {
	branchName := event.GetPullRequest().GetBase().GetRef()

	shouldRun := false
	for _, branch := range branches {
		if branch == "*" {
			shouldRun = true
			break
		}

		if branchName == branch {
			shouldRun = true
		}
	}

	return shouldRun
}

func getBranchNameFromRef(ref string) (string, error) {
	if !strings.Contains(ref, "refs/heads/") {
		return ref, nil
	}

	parts := strings.SplitAfter(ref, "refs/heads/")
	if len(parts) != 2 {
		return "", fmt.Errorf("failed to get branch name from ref %v\n", ref)
	}

	return parts[1], nil
}
