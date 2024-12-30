package git

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/spyrosmoux/cicd/common/dto"
)

func CloneRepo(repoMeta dto.Metadata, dir string) error {
	repoUrl := "github.com/" + repoMeta.RepoOwner + "/" + repoMeta.Repository + ".git"

	isPrivate, err := isPrivate(repoMeta.RepoVisibility)
	if err != nil {
		return err
	}

	var baseUrl string
	if isPrivate {
		baseUrl = "https://x-access-token:" + repoMeta.VcsToken + "@"
	} else {
		baseUrl = "https://"
	}

	normalizedUrl := baseUrl + repoUrl

	cmd := exec.Command("git", "clone", normalizedUrl, dir+"/"+repoMeta.Repository)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	slog.Info("Output: " + string(output))
	return nil
}

func CheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "fetch", "origin")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fetching origin, %s", err.Error())
	}
	slog.Info("Output: " + string(output))

	cmd = exec.Command("git", "checkout", "-b", branchName, "origin/"+branchName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error checking out remote branch %s, %s", branchName, err.Error())
	}

	slog.Info("Output: " + string(output))
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
