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

type RepoVisibility int

const (
	PUBLIC RepoVisibility = iota
	PRIVATE
)

func (rV RepoVisibility) String() string {
	switch rV {
	case PUBLIC:
		return "Public"
	case PRIVATE:
		return "Private"
	default:
		return "Unknown"
	}
}

type PublishRunDto struct {
	PipelineAsBytes []byte   `json:"pipelineAsBytes"`
	Metadata        Metadata `json:"metadata"`
}

type Metadata struct {
	Repository     string         `json:"repository"`
	Branch         string         `json:"branch"`
	RepoOwner      string         `json:"repoOwner"`
	RepoVisibility RepoVisibility `json:"repoVisibility"`
	VcsSource      VcsSource      `json:"vcsSource"`
	VcsToken       string         `json:"vcsToken"`
}
