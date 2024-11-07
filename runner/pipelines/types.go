package pipelines

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

// Job represents an individual job in the CI/CD pipeline
type Job struct {
	Name  string   `yaml:"name"`
	Needs []string `yaml:"needs,omitempty"`
	Steps []Step   `yaml:"steps"`
}

// Step represents an individual step within a job
type Step struct {
	Name string `yaml:"name"`
	Run  string `yaml:"run"`
}

type Triggers struct {
	Branch []string `yaml:"branch,omitempty"`
	PR     []string `yaml:"pr,omitempty"`
}

// UnifiedCI represents the top-level structure containing jobs
type UnifiedCI struct {
	Triggers  Triggers          `yaml:"triggers,omitempty"`
	Variables map[string]string `yaml:"variables,omitempty"`
	Jobs      []Job             `yaml:"jobs"`
}

// ValidateYAMLStructure validates the structure of a given YAML content.
func ValidateYAMLStructure(yamlData []byte) (UnifiedCI, error) {
	// Parse YAML data into UnifiedCI struct
	var ci UnifiedCI
	err := yaml.Unmarshal(yamlData, &ci)
	if err != nil {
		return UnifiedCI{}, fmt.Errorf("error unmarshalling YAML: %w", err)
	}

	// Perform validation checks
	if len(ci.Jobs) == 0 {
		return UnifiedCI{}, errors.New("no jobs defined in the YAML")
	}

	for _, job := range ci.Jobs {
		if job.Name == "" {
			return UnifiedCI{}, errors.New("job name is required")
		}
		if len(job.Steps) == 0 {
			return UnifiedCI{}, fmt.Errorf("job '%s' has no steps defined", job.Name)
		}

		for _, step := range job.Steps {
			if step.Name == "" {
				return UnifiedCI{}, fmt.Errorf("step in job '%s' is missing a name", job.Name)
			}
			if step.Run == "" {
				return UnifiedCI{}, fmt.Errorf("step '%s' in job '%s' is missing a run command", step.Name, job.Name)
			}
		}

		// Validate 'Needs' - check if all dependencies are defined as jobs
		for _, need := range job.Needs {
			found := false
			for _, j := range ci.Jobs {
				if j.Name == need {
					found = true
					break
				}
			}
			if !found {
				return UnifiedCI{}, fmt.Errorf("job '%s' has an undefined dependency '%s'", job.Name, need)
			}
		}
	}

	// Additional validation for variables, if needed
	// e.g., checking for the presence of required variables

	return ci, nil
}
