package git

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"

	"github.com/spyrosmoux/cicd/common/dto"
)

func CloneRepo(repoMeta dto.Metadata, dir string) error {
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

	slog.Debug("executing", "cmd", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed output=%s", string(output))
	}

	slog.Debug("git clone succeeded", "output", string(output))
	return nil
}

func CheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "fetch", "origin")

	slog.Debug("executing", "cmd", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fetching origin, %s", err.Error())
	}
	slog.Debug("executed successfully", "cmd", cmd.String(), "output", string(output))

	// skip checkout if branch is already checked out
	skip, err := shouldSkipCheckout(branchName)
	if err != nil {
		return err
	}

	if skip {
		return nil
	}

	slog.Debug("will checkout branch ", "branch", branchName)

	cmd = exec.Command("git", "switch", "-c", branchName, "origin/"+branchName)

	slog.Debug("executing", "cmd", cmd.String())

	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error checking out remote branch %s, %s", branchName, err.Error())
	}

	slog.Debug("executed successfully", "cmd", cmd.String(), "output", string(output))
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
