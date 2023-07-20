package api

const JiraIssueKeyRegex = "[A-Z][A-Z0-9_]*-\\d+"

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
		IssueType struct {
			Name string `json:"name"`
		}
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
		Project   struct {
			Name string
			Key  string
		}
		Components []struct {
			Name string
			Description string
		}
		Description struct {
			Content []struct {
				Type string
				// TODO return paragraphs here  ?
			}
		}
		TimeTracking struct {
			OriginalEstimate string
			RemainingEstimate string
			TimeSpent string
		}
		Comment struct {
			Total int
		}
		// TODO add worklog ?
	} `json:"fields"`
}

