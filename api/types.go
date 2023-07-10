package api

import "time"

type User struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	AccountId   string `json:"account_id"`
	Nickname    string `json:"nickname"`
	Links       struct {
		Html struct {
			Href string
		}
	}
}

type PullRequest struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	CommentCount int       `json:"comment_count"`
	TaskCount    int       `json:"task_count"`
	CreatedOn    time.Time `json:"created_on"`
	UpdatedOn    time.Time `json:"updated_on"`
	Author       User      `json:"author"`
	Destination  struct {
		Branch struct {
			Name string `json:"name"`
		}
	}
	Source struct {
		Branch struct {
			Name string `json:"name"`
		}
	}
}
