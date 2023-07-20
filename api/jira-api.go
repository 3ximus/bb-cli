package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type IssuesPaginatedResponse struct {
	Issues     []JiraIssue
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

func JiraEndpoint(domain string) string {
	return fmt.Sprintf("https://%s.atlassian.net/rest/api/3", domain)
}

func JiraBrowse(domain string, key string) string {
	return fmt.Sprintf("https://%s.atlassian.net/browse/%s", domain, key)
}

func jiraApiGet(endpoint string) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", JiraEndpoint(viper.GetString("jira_domain")), endpoint)
	// fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	cobra.CheckErr(err)
	req.SetBasicAuth(viper.GetString("email"), viper.GetString("jira_token"))

	resp, err := client.Do(req)
	cobra.CheckErr(err)

	if resp.StatusCode != 200 {
		errBody, err := ioutil.ReadAll(resp.Body)
		cobra.CheckErr(err)
		cobra.CheckErr(string(errBody))
	}

	body, err := ioutil.ReadAll(resp.Body)
	cobra.CheckErr(err)

	return body
}

// HIGH LEVEL METHODS

func GetIssue(repository string, key string) <-chan JiraIssue {
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

func GetIssueList(repository string, nResults int, all bool, reporter bool, project string) <-chan JiraIssue {
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
			query += "project=DP"
		}

		// response := jiraApiGet(fmt.Sprintf("/issue/DP-1167"))
		response := jiraApiGet(fmt.Sprintf("search?maxResults=%d&fields=*all&jql=%s", nResults, query))
		err := json.Unmarshal(response, &paginatedReponse)
		cobra.CheckErr(err)
		for _, issue := range paginatedReponse.Issues {
			channel <- issue
		}
	}()
	return channel
}
