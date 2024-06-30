package models

type JobModel struct {
	BuildId     string `json:"build_id"`
	BuildNumber int    `json:"build_number"`
	Repository  string `json:"repository"`
	Branch      string `json:"branch"`
	Author      string `json:"author"`
}

// GhPushWebhook struct is the object returned
// by GitHub's Push Webhook
type GhPushWebhook struct {
	After      string     `json:"after"`
	BaseRef    string     `json:"base_ref"`
	Before     string     `json:"before"`
	Commits    []Commits  `json:"commits"`
	Compare    string     `json:"compare"`
	Created    bool       `json:"created"`
	Deleted    bool       `json:"deleted"`
	Forced     bool       `json:"forced"`
	HeadCommit HeadCommit `json:"head_commit"`
	Pusher     Pusher     `json:"pusher"`
	Ref        string     `json:"ref"`
	Repository `json:"repository"`
}

type Commits struct {
	Added     []string `json:"added"`
	Author    Author   `json:"author"`
	Commiter  Commiter `json:"committer"`
	Distinct  bool     `json:"distinct"`
	Id        string   `json:"id"`
	Message   string   `json:"message"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
	Timestamp string   `json:"timestamp"`
	TreeId    string   `json:"tree_id"`
	Url       string   `json:"url"`
}

type Author struct {
	Date     string `json:"date"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Commiter struct {
	Date     string `json:"date"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type HeadCommit struct {
	Added     []string `json:"added"`
	Author    Author   `json:"author"`
	Commiter  Commiter `json:"committer"`
	Distinct  bool     `json:"distinct"`
	Id        string   `json:"id"`
	Message   string   `json:"message"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
	Timestamp string   `json:"timestamp"`
	TreeId    string   `json:"tree_id"`
	Url       string   `json:"url"`
}

type Pusher struct {
	Date     string `json:"date"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Repository struct {
	Id          int    `json:"id"`
	NodeId      string `json:"node_id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Private     bool   `json:"private"`
	Owner       Owner  `json:"owner"` // TODO
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
	Fork        bool   `json:"fork"`
	Url         string `json:"url"`
	KeysUrl     string `json:"keys_url"`
}
