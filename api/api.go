package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type PaginatedResponse[T any] struct {
	Values   []T
	Size     int    `json:"size"`
	Page     int    `json:"page"`
	PageLen  int    `json:"pagelen"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func api_get(endpoint string) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("api"), endpoint)
	// fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	cobra.CheckErr(err)
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("token"))

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

func api_post(endpoint string, body io.Reader) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("api"), endpoint)

	req, err := http.NewRequest("POST", url, body)
	cobra.CheckErr(err)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("token"))

	resp, err := client.Do(req)
	cobra.CheckErr(err)

	if resp.StatusCode != 201 {
		errBody, err := ioutil.ReadAll(resp.Body)
		cobra.CheckErr(err)
		cobra.CheckErr(string(errBody))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	cobra.CheckErr(err)

	return responseBody
}

// HIGH LEVEL GET METHODS

func GetUser() User {
	response := api_get("user")

	// decode response
	var user User
	err := json.Unmarshal(response, &user)
	cobra.CheckErr(err)

	return user
}

func GetPrList(repository string, state string, author string, search string, pages int) []PullRequest {
	stateQuery := ""
	if state != "" {
		stateQuery = fmt.Sprintf("state = \"%s\"", state)
	}
	authorQuery := ""
	if author != "" {
		authorQuery = fmt.Sprintf(" AND author.nickname = \"%s\"", author)
	}
	searchQuery := ""
	if search != "" {
		searchQuery = fmt.Sprintf(" AND title ~ \"%s\"", search)
	}

	var prs []PullRequest
	var prevResponse PaginatedResponse[PullRequest]
	for i := 0; i < pages; i++ {
		var response []byte
		// fmt.Printf("%v\n", prevResponse)
		if i == 0 {
			response = api_get(fmt.Sprintf("repositories/%s/pullrequests?q=%s", repository, url.QueryEscape(stateQuery+authorQuery+searchQuery+"")))
		} else {
			newUrl := strings.Replace(prevResponse.Next, "https://api.bitbucket.org/2.0/", "", 1)
			if newUrl == "" {
				break // there's no next page
			}
			response = api_get(newUrl)
		}
		err := json.Unmarshal(response, &prevResponse)
		cobra.CheckErr(err)
		prs = append(prs, prevResponse.Values...)
	}

	return prs
}

func GetPr(repository string, id int) <-chan PullRequest {
	channel := make(chan PullRequest)
	go func() {
		defer close(channel)
		var pr PullRequest
		response := api_get(fmt.Sprintf("repositories/%s/pullrequests/%d", repository, id))
		err := json.Unmarshal(response, &pr)
		cobra.CheckErr(err)
		channel <- pr
	}()
	return channel
}

func GetPrStatuses(repository string, id int) <-chan []CommitStatus {
	channel := make(chan []CommitStatus)
	go func() {
		defer close(channel)
		var paginatedResponse PaginatedResponse[CommitStatus]
		response := api_get(fmt.Sprintf("repositories/%s/pullrequests/%d/statuses", repository, id))
		err := json.Unmarshal(response, &paginatedResponse)
		cobra.CheckErr(err)
		channel <- paginatedResponse.Values
	}()
	return channel
}

func GetReviewers(repository string) <-chan []User {
	channel := make(chan []User)
	go func() {
		defer close(channel)
		var paginatedResponse PaginatedResponse[User]
		response := api_get(fmt.Sprintf("repositories/%s/effective-default-reviewers", repository))
		err := json.Unmarshal(response, &paginatedResponse)
		cobra.CheckErr(err)
		channel <- paginatedResponse.Values
	}()
	return channel
}

func GetWorkspaceMembers(workspace string) <-chan []User {
	channel := make(chan []User)
	go func() {
		defer close(channel)
		var paginatedResponse PaginatedResponse[struct {
			User User `json:"user"`
		}]
		response := api_get(fmt.Sprintf("workspaces/%s/members", workspace))
		err := json.Unmarshal(response, &paginatedResponse)
		cobra.CheckErr(err)
		var users []User
		for _, r := range paginatedResponse.Values {
			users = append(users, r.User)
		}
		channel <- users
	}()
	return channel
}

// HIGH LEVEL POST METHODS

func PostPr(repository string, data CreatePullRequest) PullRequest {
	content, err := json.Marshal(data)
	cobra.CheckErr(err)
	// fmt.Println(content)
	cobra.CheckErr(string(content))
	response := api_post(fmt.Sprintf("repositories/%s/pullrequests", repository), bytes.NewReader(content))

	// decode response
	var pr PullRequest
	err = json.Unmarshal(response, &pr)
	cobra.CheckErr(err)
	return pr
}
