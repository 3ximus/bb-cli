package api

import (
	"errors"
	"fmt"
)

const JiraIssueKeyRegex = "[A-Z][A-Z0-9_]*-\\d+"

type IssueStatus string

const (
	TODO       IssueStatus = "todo"
	INPROGRESS IssueStatus = "inprogress"
	TESTING    IssueStatus = "testing"
	DONE       IssueStatus = "done"
	BLOCKED    IssueStatus = "blocked"
)

type JiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string
		Creator struct {
			AccountId   string `json:"accountId"`
			DisplayName string `json:"displayName"`
		}
		Reporter struct {
			AccountId   string `json:"accountId"`
			DisplayName string `json:"displayName"`
		}
		Type struct {
			Name    string
			Subtask bool
		} `json:"issuetype"`
		Assignee struct {
			AccountId   string `json:"accountId"`
			DisplayName string `json:"displayName"`
		}
		Status struct {
			Name string
		}
		Priority struct {
			Name string
			Id   string
		}
		Parent struct {
		}
		Project struct {
			Name string
			Key  string
		}
		Components []struct {
			Name        string
			Description string
		}
		Description struct {
			Content []struct {
				Type string
				// TODO return paragraphs here  ?
			}
		}
		TimeTracking struct {
			OriginalEstimate  string
			RemainingEstimate string
			TimeSpent         string
		}
		Comment struct {
			Total int
		}
	} `json:"fields"`
}

type JiraTransition struct {
	Id   string
	Name string
	To   struct {
		Id   string
		Name string
	}
}

// DEFAULT ACTIONS OVERRIDES

// String is used both by fmt.Print and by Cobra in help text
func (e *IssueStatus) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *IssueStatus) Set(v IssueStatus) error {
	switch v {
	case TODO, INPROGRESS, TESTING, DONE, BLOCKED:
		*e = IssueStatus(v)
		return nil
	default:
		return errors.New(fmt.Sprintf(`must be one of "%s", "%s", "%s", "%s" or "%s"`, TODO, INPROGRESS, TESTING, DONE, BLOCKED))
	}
}

// Type is only used in help text
func (e *IssueStatus) Type() string {
	return "state"
}
