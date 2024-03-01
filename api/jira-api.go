// vim: foldmethod=indent foldnestmax=1

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type IssuesPaginatedResponse struct {
	Issues     []JiraIssue
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

type TransitionsPaginatedResponse struct {
	Expand      string
	Transitions []JiraTransition
}

func JiraEndpoint(domain string) string {
	return fmt.Sprintf("https://%s.atlassian.net/rest/api/3", domain)
}

func JiraBrowse(domain string, key string) string {
	return fmt.Sprintf("https://%s.atlassian.net/browse/%s", domain, key)
}

// REST

func jiraApiGet(endpoint string) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", JiraEndpoint(viper.GetString("jira_domain")), endpoint)
	req, err := http.NewRequest("GET", url, nil)
	cobra.CheckErr(err)
	req.SetBasicAuth(viper.GetString("email"), viper.GetString("jira_token"))
	resp, err := client.Do(req)
	cobra.CheckErr(err)
	if resp.StatusCode != 200 {
		errBody, err := io.ReadAll(resp.Body)
		cobra.CheckErr(err)
		cobra.CheckErr(string(errBody))
	}
	body, err := io.ReadAll(resp.Body)
	cobra.CheckErr(err)
	return body
}

func _jiraApiPostPut(method string, endpoint string, body io.Reader) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", JiraEndpoint(viper.GetString("jira_domain")), endpoint)
	req, err := http.NewRequest(method, url, body)
	cobra.CheckErr(err)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.SetBasicAuth(viper.GetString("email"), viper.GetString("jira_token"))
	resp, err := client.Do(req)
	cobra.CheckErr(err)
	if resp.StatusCode != 204 && resp.StatusCode != 201 && resp.StatusCode != 200 {
		errBody, err := io.ReadAll(resp.Body)
		cobra.CheckErr(err)
		cobra.CheckErr(string(errBody))
	}
	responseBody, err := io.ReadAll(resp.Body)
	cobra.CheckErr(err)
	return responseBody
}

func jiraApiPost(endpoint string, body io.Reader) []byte {
	return _jiraApiPostPut("POST", endpoint, body)
}

func jiraApiPut(endpoint string, body io.Reader) []byte {
	return _jiraApiPostPut("PUT", endpoint, body)
}

// HIGH LEVEL METHODS

func GetIssue(key string) <-chan JiraIssue {
	channel := make(chan JiraIssue)
	go func() {
		defer close(channel)
		var issue JiraIssue
		response := jiraApiGet(fmt.Sprintf("/issue/%s", key))
		err := json.Unmarshal(response, &issue)
		cobra.CheckErr(err)
		channel <- issue
	}()
	return channel
}

func GetIssueList(nResults int, all bool, reporter bool, project string, statuses []string, types []string, searchTerm string, prioritySort bool) <-chan JiraIssue {
	channel := make(chan JiraIssue)
	go func() {
		defer close(channel)
		var paginatedReponse IssuesPaginatedResponse

		query := ""
		if !reporter && !all {
			query += "assignee=currentuser()"
		} else if reporter {
			query += "reporter=currentuser()"
		}
		if project != "" {
			if query != "" {
				query += "+AND+"
			}
			query += fmt.Sprintf("project=%s", url.QueryEscape(project))
		}
		if searchTerm != "" {
			if query != "" {
				query += "+AND+"
			}
			query += fmt.Sprintf("summary~\"%s\"", url.QueryEscape(searchTerm))
		}
		if len(statuses) > 0 {
			if query != "" {
				query += "+AND+"
			}
			query += "("
			for i, s := range statuses {
				if i == 0 {
					query += fmt.Sprintf("status=\"%s\"", url.QueryEscape(s))
				} else {
					query += fmt.Sprintf("+OR+status=\"%s\"", url.QueryEscape(s))
				}
			}
			query += ")"
		}
		if len(types) > 0 {
			if query != "" {
				query += "+AND+"
			}
			query += "("
			for i, s := range types {
				if i == 0 {
					query += fmt.Sprintf("type=\"%s\"", url.QueryEscape(s))
				} else {
					query += fmt.Sprintf("+OR+type=\"%s\"", url.QueryEscape(s))
				}
			}
			query += ")"
		}
		if prioritySort {
			query += "+order+by+priority+desc,status+asc"
		} else {
			query += "+order+by+status+asc,priority+desc"
		}

		response := jiraApiGet(fmt.Sprintf("search?maxResults=%d&fields=*all&jql=%s", nResults, query))
		err := json.Unmarshal(response, &paginatedReponse)
		cobra.CheckErr(err)
		for _, issue := range paginatedReponse.Issues {
			channel <- issue
		}
	}()
	return channel
}

func GetTransitions(key string) <-chan []JiraTransition {
	channel := make(chan []JiraTransition)
	go func() {
		defer close(channel)
		var data TransitionsPaginatedResponse
		response := jiraApiGet(fmt.Sprintf("/issue/%s/transitions", key))
		err := json.Unmarshal(response, &data)
		cobra.CheckErr(err)
		channel <- data.Transitions
	}()
	return channel
}

func PostTransitions(key string, transition string) {
	var transitionDTO = struct {
		Transition struct {
			Id string `json:"id"`
		} `json:"transition"`
	}{}
	transitionDTO.Transition.Id = transition
	content, err := json.Marshal(transitionDTO)
	cobra.CheckErr(err)
	jiraApiPost(fmt.Sprintf("/issue/%s/transitions", key), bytes.NewReader(content))
}

// Post worklog in seconds and set start time to seconds before current time
func PostWorklog(key string, seconds int) {
	var worklogDTO = struct {
		// TODO add a Started field that's the current time - seconds
		TimeSpent int    `json:"timeSpentSeconds"`
		Started   string `json:"started"`
	}{}
	worklogDTO.TimeSpent = seconds
	worklogDTO.Started = time.Now().Add(time.Duration(-seconds) * time.Second).Format(time.RFC3339)
	content, err := json.Marshal(worklogDTO)
	cobra.CheckErr(err)
	jiraApiPost(fmt.Sprintf("/issue/%s/worklog", key), bytes.NewReader(content))
}

func UpdateIssue(key string, data UpdateIssueRequestBody) JiraIssue {
	content, err := json.Marshal(data)
	cobra.CheckErr(err)
	fmt.Println(string(content))
	response := jiraApiPut(fmt.Sprintf("/issue/%s?returnIssue=true", key), bytes.NewReader(content))
	var issue JiraIssue
	err = json.Unmarshal(response, &issue)
	cobra.CheckErr(err)
	return issue
}
