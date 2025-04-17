package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spyrosmoux/cicd/common/dto"
	"github.com/spyrosmoux/cicd/common/queue"
)

type GitClient interface {
	CloneRepo(repoMeta dto.Metadata, dir string) error
	CheckoutBranch(branchName string) error
}

type gitClient struct {
	logs *logrus.Entry
}

func NewGitClient(logs *logrus.Entry) GitClient {
	return &gitClient{
		logs: logs,
	}
}

func (gitClient gitClient) CloneRepo(repoMeta dto.Metadata, dir string) error {
	repoUrl := repoMeta.RepoOwner + "/" + repoMeta.Repository + ".git"

	isPrivate, err := isPrivate(repoMeta.RepoVisibility)
	if err != nil {
		return err
	}

	var baseUrl string
	if isPrivate {
		baseUrl = "https://x-access-token:" + repoMeta.VcsToken + "@github.com/"
	} else {
		baseUrl = "https://github.com/"
	}

	normalizedUrl := baseUrl + repoUrl

	targetDir := filepath.Join(dir, repoMeta.Repository)

	cmd := exec.Command("git", "clone", normalizedUrl, targetDir)

	pipelineRunId := gitClient.logs.Context.Value("pipelineRunId").(string)
	queue.PublishLog(pipelineRunId, dto.LogEntryDto{
		RunId:     pipelineRunId,
		Timestamp: time.Now().UTC().String(),
		LogLevel:  "INFO",
		Message:   fmt.Sprintf("executing cmd=%s", cmd.String()),
	})

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed output=%s", string(output))
	}

	queue.PublishLog(pipelineRunId, dto.LogEntryDto{
		RunId:     pipelineRunId,
		Timestamp: time.Now().UTC().String(),
		LogLevel:  "INFO",
		Message:   "git clone succeeded",
	})
	return nil
}

func (gitClient gitClient) CheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "fetch", "origin")

	pipelineRunId := gitClient.logs.Context.Value("pipelineRunId").(string)
	queue.PublishLog(pipelineRunId, dto.LogEntryDto{
		RunId:     pipelineRunId,
		Timestamp: time.Now().UTC().String(),
		LogLevel:  "INFO",
		Message:   fmt.Sprintf("executing cmd=%s", cmd.String()),
	})

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fetching origin, %s", err.Error())
	}

	queue.PublishLog(pipelineRunId, dto.LogEntryDto{
		RunId:     pipelineRunId,
		Timestamp: time.Now().UTC().String(),
		LogLevel:  "INFO",
		Message:   fmt.Sprintf("executed successfully cmd=%s output=%s", cmd.String(), string(output)),
	})

	// skip checkout if branch is already checked out
	skip, err := shouldSkipCheckout(branchName)
	if err != nil {
		return err
	}

	if skip {
		return nil
	}

	queue.PublishLog(pipelineRunId, dto.LogEntryDto{
		RunId:     pipelineRunId,
		Timestamp: time.Now().UTC().String(),
		LogLevel:  "INFO",
		Message:   fmt.Sprintf("will checkout branch=%s", branchName),
	})

	cmd = exec.Command("git", "switch", "-c", branchName, "origin/"+branchName)

	queue.PublishLog(pipelineRunId, dto.LogEntryDto{
		RunId:     pipelineRunId,
		Timestamp: time.Now().UTC().String(),
		LogLevel:  "INFO",
		Message:   fmt.Sprintf("executing cmd=%s", cmd.String()),
	})

	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error checking out remote branch %s, %s", branchName, err.Error())
	}

	queue.PublishLog(pipelineRunId, dto.LogEntryDto{
		RunId:     pipelineRunId,
		Timestamp: time.Now().UTC().String(),
		LogLevel:  "INFO",
		Message:   fmt.Sprintf("executed successfully cmd=%s output=%s", cmd.String(), string(output)),
	})

	return nil
}

func isPrivate(repoVisibility dto.RepoVisibility) (bool, error) {
	switch repoVisibility {
	case dto.PRIVATE:
		return true, nil
	case dto.PUBLIC:
		return false, nil
	default:
		return false, fmt.Errorf("unknown repo visibility %s\n", repoVisibility.String())
	}
}

// shouldSkipCheckout returns true if the checked out branch is the same
// as the one we want to checkout. If they are different, then continue with the checkout
func shouldSkipCheckout(desiredBranch string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	currentBranch, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error getting current branch name err=%s\n", err.Error())
	}

	currentBranch = bytes.Trim(currentBranch, "\n")

	if string(currentBranch) != desiredBranch {
		return false, nil
	}
	return true, nil
}
