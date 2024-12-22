package dirmanagement

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	RUNNER_DIR   = "/flowforge"
	WORK_DIR     = "_work"
	ARTIFACT_DIR = "a"
	SOURCE_DIR   = "s"
)

var GlobalDM *DirManagement

type DirManagement struct {
	CurrentDir  string
	WorkDir     string
	ArtifactDir string
	SourceDir   string
}

type DirManager interface {
	GetCurrentDir() string
	GetWorkDir() string
	GetArtifactDir() string
	GetSourceDir() string
	SetCurrentDir(dirName string) (string, error)
	SetWorkDir(dirName string) (string, error)
	SetArtifactDir(dirName string) (string, error)
	SetSourceDir(dirName string) (string, error)
}

func InitGlobalDM() error {
	GlobalDM = &DirManagement{
		CurrentDir:  "",
		WorkDir:     "",
		ArtifactDir: "",
		SourceDir:   "",
	}

	err := GlobalDM.CreateDirectory(RUNNER_DIR)
	if err != nil {
		return err
	}

	_, err = GlobalDM.SetCurrentDir(RUNNER_DIR)
	if err != nil {
		return err
	}

	return nil
}

func (dm *DirManagement) GetCurrentDir() string {
	return dm.CurrentDir
}

func (dm *DirManagement) GetWorkDir() string {
	return dm.WorkDir
}

func (dm *DirManagement) GetArtifactDir() string {
	return dm.ArtifactDir
}

func (dm *DirManagement) GetSourceDir() string {
	return dm.SourceDir
}

func (dm *DirManagement) SetCurrentDir(dirName string) (string, error) {
	abs, err := findDir(dirName)
	if err != nil {
		return "", err
	}
	dm.CurrentDir = abs

	err = os.Chdir(abs)
	if err != nil {
		return "", fmt.Errorf("error changing to %s directory, %s", abs, err.Error())
	}
	return dm.CurrentDir, nil
}

func (dm *DirManagement) SetWorkDir(dirName string) (string, error) {
	abs, err := findDir(dirName)
	if err != nil {
		return "", err
	}
	dm.WorkDir = abs
	return dm.WorkDir, nil
}

func (dm *DirManagement) SetArtifactDir(dirName string) (string, error) {
	abs, err := findDir(dirName)
	if err != nil {
		return "", err
	}
	dm.ArtifactDir = abs
	return dm.ArtifactDir, nil
}

func (dm *DirManagement) SetSourceDir(dirName string) (string, error) {
	abs, err := findDir(dirName)
	if err != nil {
		return "", err
	}
	dm.SourceDir = abs
	return dm.SourceDir, nil
}

// CreateDirectory creates a directory from a given string. The string can be the directory's
//
//	or a path.
func (dm *DirManagement) CreateDirectory(dirName string) error {
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory %s, %s", dirName, err.Error())
	}
	return nil
}

// findDir returns the absolute path of a given directory name, else errors
func findDir(dirName string) (string, error) {
	abs, err := filepath.Abs(dirName)
	if err != nil {
		return "", fmt.Errorf("error finding absolute path for directory %s, %s", dirName, err.Error())
	}
	return abs, err
}
