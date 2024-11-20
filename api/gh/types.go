package gh

/*
INFO: https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation
*/

type AccessTokenResponse struct {
	AccessToken         string     `json:"token"`
	ExpiresAt           string     `json:"expires_at"`
	Permissions         Permission `json:"permissions"`
	RepositorySelection string     `json:"repository_selection"`
}

type Permission struct {
	Contents        string `json:"contents"`
	Metadata        string `json:"metadata"`
	PullRequests    string `json:"pull_requests"`
	RepositoryHooks string `json:"repository_hooks"`
	Statuses        string `json:"statuses"`
}
