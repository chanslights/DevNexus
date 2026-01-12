package types

// Web
type WebhookPayload struct {
	RepoName string `json:"repo_name"` // 仓库名
	Branch   string `json:"branch"`    // 分支
	CommitID string `json:"commit_id"` // 最新的Commit SHA
	Pusher   string `json:"pusher"`    // 推送人
}
