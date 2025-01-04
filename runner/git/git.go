package git

import (
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

	slog.Info("executing", "cmd", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed output=%s", string(output))
	}

	slog.Info("git clone succeeded", "output", string(output))
	return nil
}

func CheckoutBranch(branchName string) error {
	cmd := exec.Command("git", "fetch", "origin")

	slog.Info("executing", "cmd", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fetching origin, %s", err.Error())
	}
	slog.Info("executed successfully", "cmd", cmd.String(), "output", string(output))

	// TODO(@spyrosmoux) add check -> if branch already exists, skip checkout

	cmd = exec.Command("git", "switch", "-c", branchName, "origin/"+branchName)

	slog.Info("executing", "cmd", cmd.String())

	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error checking out remote branch %s, %s", branchName, err.Error())
	}

	slog.Info("executed successfully", "cmd", cmd.String(), "output", string(output))
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
