package api

type User struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	AccountId   string `json:"account_id"`
	Links       struct {
		Html struct {
			Href string
		}
	}
}

type PullRequest struct {
	ID           int `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	State        string `json:"state"`
	CommentCount int `json:"comment_count"`
	TaskCount    int `json:"task_count"`
}
