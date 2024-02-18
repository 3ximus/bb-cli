// vim: foldmethod=indent foldnestmax=1

package api

import (
	"errors"
	"fmt"
	"time"
)

type ErrorResponse struct {
	Error struct {
		Message string
		Detail  string
	}
}

type User struct {
	UUID        string `json:"uuid"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	AccountId   string `json:"account_id"`
	Nickname    string `json:"nickname"`
	Links       struct{ Html struct{ Href string } }
}

type PrState string

const (
	OPEN       PrState = "open"
	MERGED     PrState = "merged"
	DECLINED   PrState = "declined"
	SUPERSEDED PrState = "superseded"
)

type Branch struct {
	Branch struct {
		Name string `json:"name"`
	} `json:"branch"`
}

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
	Destination  Branch
	Source       Branch
	Links        struct{ Html struct{ Href string } }
	Status       CommitStatus
	CreatedOn    time.Time `json:"created_on"`
	UpdatedOn    time.Time `json:"updated_on"`
	Participants []struct {
		User           User `json:"user"`
		Role           string
		Approved       bool
		State          string
		ParticipatedOn time.Time `json:"participated_on"`
	}
}

type CreatePullRequestBody struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CloseSource bool    `json:"close_source_branch"`
	Destination *Branch `json:"destination,omitempty"`
	Source      *Branch `json:"source,omitempty"`
	Reviewers   []struct {
		AccountId string `json:"account_id"`
	} `json:"reviewers,omitempty"`
}

type RunPipelineRequestBody struct {
	Target struct {
		RefType     string                   `json:"ref_type"`
		Type        string                   `json:"type"`
		RefName     string                   `json:"ref_name"`
		Commit      *PipelineCommitRefBody   `json:"commit"`
		PullRequest *PipelinePullRequestBody `json:"pullrequest"`
		Selector    *PipelineSelectorBody    `json:"selector"`
	} `json:"target"`
}

type PipelineSelectorBody struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

type PipelineCommitRefBody struct {
	Hash string `json:"hash"`
	Type string `json:"type"`
}

type PipelinePullRequestBody struct {
	Id string `json:"id"`
}

type CommitStatus struct {
	RefName   string    `json:"refname"`
	Name      string    `json:"name"`
	State     string    `json:"state"`
	Url       string    `json:"url"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

type Pipeline struct {
	UUID        string
	BuildNumber int `json:"build_number"`
	State       struct {
		Name   string `json:"name"`
		Result struct {
			Name string
		}
	} `json:"state"`
	Target struct {
		Source      string
		Destination string
		RefName     string `json:"ref_name"`
		PullRequest struct {
			Id    int
			Title string
			Links struct{ Html struct{ Href string } }
		} `json:"pullrequest"`
	}
	Trigger struct {
		Name string
	}
	Author            User      `json:"creator"`
	DurationInSeconds int       `json:"duration_in_seconds"`
	CompletedOn       time.Time `json:"completed_on"`
	CreatedOn         time.Time `json:"created_on"`
}

type PipelineStep struct {
	UUID              string
	Name              string
	DurationInSeconds int `json:"duration_in_seconds"`
	State             struct {
		Name   string
		Result struct {
			Name string
		}
		Stage struct {
			Name string
		}
	}
	SetupCommands    []StepCommand `json:"setup_commands"`
	ScriptCommands   []StepCommand `json:"script_commands"`
	TeardownCommands []StepCommand `json:"teardown_commands"`
	Image            struct {
		Name string
	}
	Pipeline struct {
		UUID string
	}
}

type PipelineReport struct {
	Total   int `json:"number_of_test_cases"`
	Success int `json:"number_of_successful_test_cases"`
	Failed  int `json:"number_of_failed_test_cases"`
	Error   int `json:"number_of_error_test_cases"`
	Skipped int `json:"number_of_skipped_test_cases"`
}

type PipelineReportCase struct {
	UUID               string
	Name               string
	FullyQualifiedName string `json:"fully_qualified_name"`
	PackageName        string `json:"package_name"`
	Status             string `json:"status"`
	Duration           string `json:"duration"`
}

type StepCommand struct {
	Name        string
	Command     string
	CommandType string
}

type PrComment struct {
	Id      int
	Content struct {
		Raw  string
		Html string
	}
	User      User
	Deleted   bool
	Type      string
	Links     struct{ Html struct{ Href string } }
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

type Environment struct {
	UUID     string `json:"uuid"`
	Name     string
	Category struct {
		Name string
	}
	EnvironmentType struct {
		Name string
	} `json:"environment_type"`
	Lock struct {
		Triggerer struct {
			PipelineUUID string `json:"pipeline_uuid"`
		}
	}
	Status Pipeline
}

type EnvironmentVariable struct {
	Key     string
	Value   string
	Secured bool
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
