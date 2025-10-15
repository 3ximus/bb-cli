package api

type Worklog struct {
	JiraWorklogID  int `json:"jiraWorklogId"`
	TempoWorklogID int `json:"tempoWorklogId"`
	Issue          struct {
		// Key string `json:"key"`
		ID int `json:"id"`
	} `json:"issue"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	StartDate        string `json:"startDate"`
	StartDateTimeUtc string `json:"startDateTimeUtc"`
	StartTime        string `json:"startTime"`
	Description      string `json:"description"`
	Author           struct {
		AccountID string `json:"accountId"`
	} `json:"author"`
}
