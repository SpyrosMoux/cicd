package pipelines

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spyrosmoux/cicd/common/dto"
	"github.com/spyrosmoux/cicd/runner/dirmanagement"
	"github.com/spyrosmoux/cicd/runner/git"
)

type service struct {
	logger *logrus.Logger
}

func NewService(logger *logrus.Logger) Service {
	return &service{logger: logger}
}

// PrepareRun clones the repo into the source directory
func (svc *service) PrepareRun(repoMeta dto.Metadata) error {
	// TODO(@SpyrosMoux) create a unique dir for the run based on the unique build number (build number -> generated by the API)
	err := git.CloneRepo(repoMeta, dirmanagement.GlobalDM.GetSourceDir())
	if err != nil {
		return err
	}

	_, err = dirmanagement.GlobalDM.SetCurrentDir(dirmanagement.GlobalDM.GetSourceDir() + "/" + repoMeta.Repository)
	if err != nil {
		return err
	}

	err = git.CheckoutBranch(repoMeta.Branch)
	if err != nil {
		return err
	}

	_, err = dirmanagement.GlobalDM.SetCurrentDir(dirmanagement.GlobalDM.GetWorkDir())
	if err != nil {
		return err
	}

	return nil
}

// CleanupRun deletes the '_work' directory and all its contents
func (svc *service) CleanupRun() error {
	pathToRemove, err := dirmanagement.GlobalDM.SetCurrentDir(dirmanagement.GlobalDM.GetWorkDir())
	if err != nil {
		return err
	}

	err = os.RemoveAll(pathToRemove)
	if err != nil {
		svc.logger.WithError(err).Error("failed to remove temporary directory")
		return err
	}

	_, err = dirmanagement.GlobalDM.SetCurrentDir(dirmanagement.RUNNER_DIR)
	if err != nil {
		return err
	}

	svc.logger.WithFields(logrus.Fields{
		"directory": pathToRemove,
	}).Info("temporary directory removed successfully")
	return nil
}

// ExecuteStep executes a single step
func (svc *service) ExecuteStep(step Step, variables map[string]string) error {
	svc.logger.WithFields(logrus.Fields{
		"step": step.Name,
	}).Info("executing step")

	command := svc.SubstituteUserVariables(step.Run, variables)
	cmd := exec.Command("sh", "-c", command)

	svc.logger.WithFields(logrus.Fields{
		"cmd": cmd,
	}).Info("executing")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failure executing step step=%s, err=%s", step.Name, err.Error())
	}

	svc.logger.Info("output ", string(output))
	return nil
}

// ExecuteJob executes all steps in a job
func (svc *service) ExecuteJob(job Job, variables map[string]string) error {
	svc.logger.WithFields(logrus.Fields{
		"job": job.Name,
	}).Info("executing job")
	for _, step := range job.Steps {
		err := svc.ExecuteStep(step, variables)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunPipeline prepares, executes and cleans-up a run
func (svc *service) RunPipeline(pipeline Pipeline, runMetadata dto.Metadata) error {
	defer func(svc *service) {
		err := svc.CleanupRun()
		if err != nil {
			svc.logger.Error(err)
		}
	}(svc)

	err := svc.PrepareRun(runMetadata)
	if err != nil {
		return err
	}

	svc.SubstitutePredefinedVariables(pipeline, PredefinedVars)

	for _, job := range pipeline.Jobs {
		svc.logger.WithFields(logrus.Fields{
			"job": job.Name,
		}).Info("running job")
		err = svc.ExecuteJob(job, pipeline.Variables)
		if err != nil {
			svc.logger.WithFields(logrus.Fields{
				"name": job.Name,
				"err":  err,
			}).Error("failed to run pipeline")
			return err
		}
	}

	return nil
}

// SubstituteUserVariables substitutes user defined variables for a specific command
func (svc *service) SubstituteUserVariables(command string, variables map[string]string) string {
	for key, value := range variables {
		placeholder := fmt.Sprintf("${%s}", key)
		command = strings.ReplaceAll(command, placeholder, value)
	}
	return command
}

// SubstitutePredefinedVariables substitutes all predefined variables used in the pipeline
func (svc *service) SubstitutePredefinedVariables(pipeline Pipeline, predefinedVariables map[string]string) Pipeline {
	for i, job := range pipeline.Jobs {
		for j, step := range job.Steps {
			pipeline.Jobs[i].Steps[j].Run = replaceVariables(step.Run, predefinedVariables)
		}
	}
	return pipeline
}

func replaceVariables(input string, variables map[string]string) string {
	for key, value := range variables {
		input = strings.ReplaceAll(input, key, value)
	}
	return input
}
