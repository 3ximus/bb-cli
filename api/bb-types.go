package api

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	UUID        string `json:"uuid"`
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

type PrState string

const (
	OPEN       PrState = "open"
	MERGED     PrState = "merged"
	DECLINED   PrState = "declined"
	SUPERSEDED PrState = "superseded"
)

type PullRequest struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	State        PrState `json:"state"`
	CommentCount int     `json:"comment_count"`
	TaskCount    int     `json:"task_count"`
	Author       User    `json:"author"`
	ClosedBy     User    `json:"closed_by"`
	CloseSource  bool    `json:"close_source_branch"`
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
	Links struct {
		Html struct {
			Href string
		}
	}
	Status    CommitStatus
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

type CreatePullRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CloseSource bool   `json:"close_source_branch"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"destination"`
	Source struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"source"`
	Reviewers []struct {
		AccountId string `json:"account_id"`
	} `json:"reviewers,omitempty"`
}

type CommitStatus struct {
	RefName   string    `json:"refname"`
	Name      string    `json:"name"`
	State     string    `json:"state"`
	Url       string    `json:"url"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

// DEFAULT ACTIONS OVERRIDES

// String is used both by fmt.Print and by Cobra in help text
func (e *PrState) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *PrState) Set(v PrState) error {
	switch v {
	case OPEN, MERGED, DECLINED, SUPERSEDED:
		*e = PrState(v)
		return nil
	default:
		return errors.New(fmt.Sprintf(`must be one of "%s", "%s", "%s" or "%s"`, OPEN, MERGED, DECLINED, SUPERSEDED))
	}
}

// Type is only used in help text
func (e *PrState) Type() string {
	return "state"
}
