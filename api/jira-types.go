package api

const JiraIssueKeyRegex = "[A-Z][A-Z0-9_]*-\\d+"

type Myself struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
}

type JiraIssue struct {
	ID     string
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
			Key    string `json:"key"`
			Fields struct {
				Summary string
				Type    struct {
					Name    string
					Subtask bool
				} `json:"issuetype"`
				Priority struct {
					Name string
					Id   string
				}
			}
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

type UpdateIssueRequestBody struct {
	Fields struct {
		// This on might be enough when set on project
		TimeTracking *TimeTracking `json:"timetracking,omitempty"`
		Summary      string        `json:"summary,omitempty"`
		Priority     struct {
			Id string `json:"id,omitempty"`
		} `json:"priority,omitempty"`
	} `json:"fields,omitempty"`
	Update struct {
		TimeTracking []UpdateType[TimeTracking] `json:"timetracking,omitempty"`
	} `json:"update,omitempty"`
}

type UpdateType[K TimeTracking] struct {
	Edit *K `json:"edit,omitempty"`
	Set  *K `json:"set,omitempty"`
}

type TimeTracking struct {
	OriginalEstimate  string `json:"originalEstimate,omitempty"`
	RemainingEstimate string `json:"remainingEstimate,omitempty"`
	TimeSpent         string `json:"timeSpent,omitempty"`
}

type JiraTransition struct {
	Id   string
	Name string
	To   struct {
		Id   string
		Name string
	}
}
