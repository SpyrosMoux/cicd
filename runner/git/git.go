package git

import (
	"fmt"
	"github.com/spyrosmoux/cicd/common/dto"
	"log/slog"
	"os/exec"
)

func CloneRepo(repoMeta dto.Metadata, dir string) error {
	repoUrl := "github.com/" + repoMeta.RepoOwner + "/" + repoMeta.Repository + ".git"
	baseUrl := "https://x-access-token:" + repoMeta.VcsToken
	normalizedUrl := baseUrl + "@" + repoUrl

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
