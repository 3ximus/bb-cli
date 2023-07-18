package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// type PaginatedResponse[T any] struct {
// 	Values   []T
// 	Size     int    `json:"size"`
// 	Page     int    `json:"page"`
// 	PageLen  int    `json:"pagelen"`
// 	Next     string `json:"next"`
// 	Previous string `json:"previous"`
// }

func jiraApiGet(endpoint string) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("jira_api"), endpoint)
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

func GetIssueList(repository string) <-chan JiraIssue {
	channel := make(chan JiraIssue)
	go func() {
		defer close(channel)
		// response := jiraApiGet(fmt.Sprintf("/issue/DP-1167"))
		response := jiraApiGet(fmt.Sprintf("search?jql="))
		fmt.Println(string(response))
	}()
	return channel
}
