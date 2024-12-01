package dto

type VcsSource int

const (
	GITHUB VcsSource = iota
)

func (s VcsSource) String() string {
	switch s {
	case GITHUB:
		return "Github"
	default:
		return "unknown"
	}
}

type PublishRunDto struct {
	PipelineAsBytes []byte   `json:"pipelineAsBytes"`
	Metadata        Metadata `json:"metadata"`
}

type Metadata struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	RepoOwner  string `json:"repoOwner"`
	// event??
	VcsSource VcsSource `json:"vcsSource"`
	VcsToken  string    `json:"vcsToken"`
}
