package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func tempoApiGet(endpoint string) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("tempo_api"), endpoint)
	req, err := http.NewRequest("GET", url, nil)
	cobra.CheckErr(err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("tempo_token")))
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

func tempoApiPost(endpoint string, body io.Reader) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("tempo_api"), endpoint)
	req, err := http.NewRequest("POST", url, body)
	cobra.CheckErr(err)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("tempo_token")))
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

func ListWorklogs(user Myself, start, end time.Time) []Worklog {
	from := start.Format("2006-01-02")
	to := end.Format("2006-01-02")

	respBody := tempoApiGet(fmt.Sprintf("/worklogs/user/%s?from=%s&to=%s", user.AccountID, from, to))
	var result struct {
		Results []Worklog `json:"results"`
	}
	err := json.Unmarshal(respBody, &result)
	cobra.CheckErr(err)
	return result.Results
}

func PostWorklog(user Myself, issueId int, seconds int, start time.Time) Worklog {
	worklog := struct {
		IssueId          int    `json:"issueId"`
		TimeSpentSeconds int    `json:"timeSpentSeconds"`
		StartDate        string `json:"startDate"`
		StartTime        string `json:"startTime"`
		AuthorId         string `json:"authorAccountId"`
	}{
		IssueId:          issueId,
		TimeSpentSeconds: seconds,
		StartDate:        start.Format("2006-01-02"),
		StartTime:        start.Format("15:04:05"), // 24-hour format
		AuthorId:         user.AccountID,
	}
	content, err := json.Marshal(worklog)
	cobra.CheckErr(err)

	resp := tempoApiPost("/worklogs", bytes.NewReader(content))

	result := Worklog{}
	err = json.Unmarshal(resp, &result)
	cobra.CheckErr(err)
	return result
}
