// vim: foldmethod=indent foldnestmax=1

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BBPaginatedResponse[T any] struct {
	Values   []T
	Size     int    `json:"size"`
	Page     int    `json:"page"`
	PageLen  int    `json:"pagelen"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func BBBrowsePipelines(repository string, id int) string {
	return fmt.Sprintf("https://bitbucket.org/%s/pipelines/results/%d", repository, id)
}

// REST

func bbApiGet(endpoint string) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("bb_api"), endpoint)
	// fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	cobra.CheckErr(err)
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("bb_token"))

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

func _bbApiPostPut(method string, endpoint string, body io.Reader) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("bb_api"), endpoint)

	req, err := http.NewRequest(method, url, body)
	cobra.CheckErr(err)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("bb_token"))

	resp, err := client.Do(req)
	cobra.CheckErr(err)

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		errBody, err := io.ReadAll(resp.Body)
		cobra.CheckErr(err)
		cobra.CheckErr(string(errBody))
	}

	responseBody, err := io.ReadAll(resp.Body)
	cobra.CheckErr(err)

	return responseBody
}

func bbApiPost(endpoint string, body io.Reader) []byte {
	return _bbApiPostPut("POST", endpoint, body)
}

func bbApiPut(endpoint string, body io.Reader) []byte {
	return _bbApiPostPut("PUT", endpoint, body)
}

func bbApiDelete(endpoint string) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("bb_api"), endpoint)
	req, err := http.NewRequest("DELETE", url, nil)
	cobra.CheckErr(err)
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("bb_token"))

	resp, err := client.Do(req)
	cobra.CheckErr(err)

	if resp.StatusCode != 204 {
		errBody, err := io.ReadAll(resp.Body)
		cobra.CheckErr(err)
		cobra.CheckErr(string(errBody))
	}

	body, err := io.ReadAll(resp.Body)
	cobra.CheckErr(err)

	return body
}

// HIGH LEVEL METHODS

func GetUser() User {
	response := bbApiGet("user")

	// decode response
	var user User
	err := json.Unmarshal(response, &user)
	cobra.CheckErr(err)

	return user
}

func GetPrList(
	repository string,
	states []string,
	author string,
	search string,
	source string,
	destination string,
	pages int,
	status bool,
	participants bool,
) <-chan PullRequest {
	channel := make(chan PullRequest)
	go func() {
		defer close(channel)

		stateQuery := ""
		if len(states) > 0 {
			stateQuery = "("
			for i, s := range states {
				if i == 0 {
					stateQuery += fmt.Sprintf("state = \"%s\"", s)
				} else {
					stateQuery += fmt.Sprintf(" OR state = \"%s\"", s)
				}
			}
			stateQuery += ")"
		}
		authorQuery := ""
		if author != "" {
			authorQuery = fmt.Sprintf(" AND author.nickname = \"%s\"", author)
		}
		searchQuery := ""
		if search != "" {
			searchQuery = fmt.Sprintf(" AND title ~ \"%s\"", search)
		}
		sourceQuery := ""
		if source != "" {
			sourceQuery = fmt.Sprintf(" AND source.branch.name = \"%s\"", source)
		}
		destinationQuery := ""
		if destination != "" {
			destinationQuery = fmt.Sprintf(" AND destination.branch.name = \"%s\"", destination)
		}
		participantsExpansion := ""
		if participants {
			// this should be fields=* but it doesn't work
			participantsExpansion = "&fields=values.id,values.title,values.description,values.state,values.comment_count,values.task_count,values.author,values.closed_by,values.close_source_branch,values.destination,values.source,values.links,values.status,values.created_on,values.updated_on,values.participants"
		}

		var prevResponse BBPaginatedResponse[PullRequest]
		for i := 0; i < pages; i++ {
			var response []byte
			if i == 0 {
				response = bbApiGet(fmt.Sprintf("repositories/%s/pullrequests?sort=-id%s&q=%s", repository, participantsExpansion, url.QueryEscape(stateQuery+authorQuery+searchQuery+sourceQuery+destinationQuery)))
			} else {
				newUrl := strings.Replace(prevResponse.Next, "https://api.bitbucket.org/2.0/", "", 1)
				if newUrl == "" {
					break // there's no next page
				}
				response = bbApiGet(newUrl)
			}
			err := json.Unmarshal(response, &prevResponse)
			cobra.CheckErr(err)

			// yield the value on the channel
			for _, pr := range prevResponse.Values {
				if status {
					status := <-GetPrStatuses(repository, pr.ID)
					if status != nil && len(status) > 0 {
						// TODO FIX instead of getting the first one get the latest one
						pr.Status = status[0] // only get the first one
					}
				}
				channel <- pr
			}
		}
	}()
	return channel
}

func GetPr(repository string, id int) <-chan PullRequest {
	channel := make(chan PullRequest)
	go func() {
		defer close(channel)
		var pr PullRequest
		response := bbApiGet(fmt.Sprintf("repositories/%s/pullrequests/%d", repository, id))
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
		var paginatedResponse BBPaginatedResponse[CommitStatus]
		response := bbApiGet(fmt.Sprintf("repositories/%s/pullrequests/%d/statuses", repository, id))
		err := json.Unmarshal(response, &paginatedResponse)
		cobra.CheckErr(err)
		channel <- paginatedResponse.Values
	}()
	return channel
}

func GetPrComments(repository string, id int) <-chan []PrComment {
	channel := make(chan []PrComment)
	go func() {
		defer close(channel)
		var paginatedResponse BBPaginatedResponse[PrComment]
		response := bbApiGet(fmt.Sprintf("repositories/%s/pullrequests/%d/comments", repository, id))
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
		var paginatedResponse BBPaginatedResponse[User]
		response := bbApiGet(fmt.Sprintf("repositories/%s/effective-default-reviewers", repository))
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
		var paginatedResponse BBPaginatedResponse[struct {
			User User `json:"user"`
		}]
		response := bbApiGet(fmt.Sprintf("workspaces/%s/members", workspace))
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

func PostPr(repository string, data CreatePullRequest) PullRequest {
	content, err := json.Marshal(data)
	cobra.CheckErr(err)
	response := bbApiPost(fmt.Sprintf("repositories/%s/pullrequests", repository), bytes.NewReader(content))

	// decode response
	var pr PullRequest
	err = json.Unmarshal(response, &pr)
	cobra.CheckErr(err)
	return pr
}

func UpdatePr(repository string, id int, data CreatePullRequest) PullRequest {
	content, err := json.Marshal(data)
	cobra.CheckErr(err)
	response := bbApiPut(fmt.Sprintf("repositories/%s/pullrequests/%d", repository, id), bytes.NewReader(content))

	// decode response
	var pr PullRequest
	err = json.Unmarshal(response, &pr)
	cobra.CheckErr(err)
	return pr
}

func ApprovePr(repository string, id int) {
	bbApiPost(fmt.Sprintf("repositories/%s/pullrequests/%d/approve", repository, id), nil)
}

func MergePr(repository string, id int, message string) {
	content, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	cobra.CheckErr(err)
	payload := bytes.NewReader(content)
	if message == "" {
		payload = nil
	}
	bbApiPost(fmt.Sprintf("repositories/%s/pullrequests/%d/merge", repository, id), payload)
}

func UnnaprovePr(repository string, id int) {
	bbApiDelete(fmt.Sprintf("repositories/%s/pullrequests/%d/approve", repository, id))
}

func DeclinePr(repository string, id int) {
	bbApiPost(fmt.Sprintf("repositories/%s/pullrequests/%d/decline", repository, id), nil)
}

func RequestChangesPr(repository string, id int) {
	bbApiPost(fmt.Sprintf("repositories/%s/pullrequests/%d/request-changes", repository, id), nil)
}

func GetPipelineList(repository string, nResults int, targetBranch string) <-chan Pipeline {
	channel := make(chan Pipeline)
	go func() {
		defer close(channel)

		query := ""
		if targetBranch != "" {
			query += fmt.Sprintf("&target.branch=%s", targetBranch)
		}

		var pipelineResponse BBPaginatedResponse[Pipeline]
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines?sort=-created_on&pagelen=%d%s", repository, nResults, query))
		err := json.Unmarshal(response, &pipelineResponse)
		cobra.CheckErr(err)
		for _, pipeline := range pipelineResponse.Values {
			channel <- pipeline
		}
	}()
	return channel
}

func GetPipeline(repository string, id string) <-chan Pipeline {
	channel := make(chan Pipeline)
	go func() {
		defer close(channel)
		var pipeline Pipeline
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines/%s", repository, id))
		err := json.Unmarshal(response, &pipeline)
		cobra.CheckErr(err)
		channel <- pipeline
	}()
	return channel
}

func GetPipelineSteps(repository string, id string) <-chan []PipelineStep {
	channel := make(chan []PipelineStep)
	go func() {
		defer close(channel)
		var steps BBPaginatedResponse[PipelineStep]
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines/%s/steps", repository, id))
		err := json.Unmarshal(response, &steps)
		cobra.CheckErr(err)
		channel <- steps.Values
	}()
	return channel
}

func GetPipelineStepLogs(repository string, id string, stepUUID string) <-chan string {
	channel := make(chan string)
	go func() {
		defer close(channel)
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines/%s/steps/%s/log", repository, id, stepUUID))
		channel <- string(response)
	}()
	return channel
}

func GetPipelineVariables(repository string) <-chan EnvironmentVariable {
	channel := make(chan EnvironmentVariable)
	go func() {
		defer close(channel)

		var environmentResponse BBPaginatedResponse[EnvironmentVariable]
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines_config/variables", repository))
		err := json.Unmarshal(response, &environmentResponse)
		cobra.CheckErr(err)
		for _, envVar := range environmentResponse.Values {
			channel <- envVar
		}
	}()
	return channel
}

func GetEnvironmentList(repository string, status bool) <-chan Environment {
	channel := make(chan Environment)
	go func() {
		defer close(channel)

		var environmentResponse BBPaginatedResponse[Environment]
		response := bbApiGet(fmt.Sprintf("repositories/%s/environments", repository))
		err := json.Unmarshal(response, &environmentResponse)
		cobra.CheckErr(err)
		for _, env := range environmentResponse.Values {
			if status {
				env.Status = <-GetPipeline(repository, env.Lock.Triggerer.PipelineUUID)
			}
			channel <- env
		}
	}()
	return channel
}

func GetEnvironmentVariables(repository string, envName string) <-chan EnvironmentVariable {
	channel := make(chan EnvironmentVariable)
	go func() {
		defer close(channel)

		for env := range GetEnvironmentList(repository, false) {
			if env.Name == envName {
				var environmentResponse BBPaginatedResponse[EnvironmentVariable]
				response := bbApiGet(fmt.Sprintf("repositories/%s/deployments_config/environments/%s/variables", repository, env.UUID))
				err := json.Unmarshal(response, &environmentResponse)
				cobra.CheckErr(err)
				for _, envVar := range environmentResponse.Values {
					channel <- envVar
				}
				break
			}
		}
	}()
	return channel
}
