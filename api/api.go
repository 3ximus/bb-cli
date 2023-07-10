package api

import (
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func get(endpoint string, spinner *spinner.Spinner) []byte {
	token := viper.GetString("token")
	username := viper.GetString("username")
	api_endpoint := viper.GetString("api")
	url := fmt.Sprintf("%s/%s", api_endpoint, endpoint)

	client := &http.Client{}

	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	cobra.CheckErr(err)
	req.SetBasicAuth(username, token)

	resp, err := client.Do(req)
	cobra.CheckErr(err)

	if resp.StatusCode != 200 {
		errBody, err := ioutil.ReadAll(resp.Body)
		spinner.Stop()
		cobra.CheckErr(err)
		cobra.CheckErr(string(errBody))
	}

	body, err := ioutil.ReadAll(resp.Body)
	cobra.CheckErr(err)

	return body
}

// api get request wrapper with a loading spinner
func api_get(endpoint string) []byte {
	s := spinner.New(
		spinner.CharSets[14],
		100*time.Millisecond,
		spinner.WithWriter(os.Stderr),
		spinner.WithSuffix(" Sending request..."),
		spinner.WithColor("fgHiBlue"),
	)
	s.Start()
	response := get(endpoint, s)
	s.Stop()
	return response
}

func GetUser() User {
	response := api_get("user")

	// decode response
	var user User
	err := json.Unmarshal(response, &user)
	cobra.CheckErr(err)

	return user
}

func GetPr(repository string, state string, author string, search string) []PullRequest {
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
	fmt.Println(stateQuery + authorQuery + searchQuery)
	response := api_get(fmt.Sprintf("repositories/%s/pullrequests?q=%s", repository, url.QueryEscape(stateQuery+authorQuery+searchQuery)))

	type PRResponse struct {
		Values []PullRequest
	}

	// fmt.Println(string(response))

	// decode response
	var prs PRResponse
	err := json.Unmarshal(response, &prs)
	cobra.CheckErr(err)

	return prs.Values
}
