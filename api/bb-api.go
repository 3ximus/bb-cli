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
	return bbApiRangedGet(endpoint, "")
}

func bbApiRangedGet(endpoint string, dataRange string) []byte {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", viper.GetString("bb_api"), endpoint)
	req, err := http.NewRequest("GET", url, nil)
	cobra.CheckErr(err)
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("bb_token"))

	if dataRange != "" {
		req.Header.Add("Range", fmt.Sprintf("bytes=%s", dataRange))
	}

	resp, err := client.Do(req)
	cobra.CheckErr(err)

	if resp.StatusCode == 404 {
		var errResponse ErrorResponse
		body, err := io.ReadAll(resp.Body)
		cobra.CheckErr(err)
		err = json.Unmarshal(body, &errResponse)
		cobra.CheckErr(err)
		cobra.CheckErr(errResponse.Error.Detail)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 206 && resp.StatusCode != 416 {
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

	if resp.StatusCode != 201 && resp.StatusCode != 200 && resp.StatusCode != 204 {
		errBody, err := io.ReadAll(resp.Body)
		cobra.CheckErr(strings.Join([]string{fmt.Sprintf("%v", err), string(errBody)}, ","))
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

func PostPr(repository string, data CreatePullRequestBody) PullRequest {
	content, err := json.Marshal(data)
	cobra.CheckErr(err)
	response := bbApiPost(fmt.Sprintf("repositories/%s/pullrequests", repository), bytes.NewReader(content))

	// decode response
	var pr PullRequest
	err = json.Unmarshal(response, &pr)
	cobra.CheckErr(err)
	return pr
}

func UpdatePr(repository string, id int, data CreatePullRequestBody) PullRequest {
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

func GetPipelineStep(repository string, id string, stepId string) <-chan PipelineStep {
	channel := make(chan PipelineStep)
	go func() {
		defer close(channel)
		var step PipelineStep
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines/%s/steps/%s", repository, id, stepId))
		err := json.Unmarshal(response, &step)
		cobra.CheckErr(err)
		channel <- step
	}()
	return channel
}

func GetPipelineStepLogs(repository string, id string, stepId string, offset int) <-chan string {
	channel := make(chan string)
	go func() {
		defer close(channel)
		response := bbApiRangedGet(fmt.Sprintf("repositories/%s/pipelines/%s/steps/%s/log", repository, id, stepId), fmt.Sprintf("%d-", offset))
		if bytes.Index(response, []byte("Range Not Satisfiable")) != -1 {
			channel <- ""
		} else {
			channel <- string(response)
		}
	}()
	return channel
}

func GetPipelineReport(repository string, id string, stepId string) <-chan PipelineReport {
	channel := make(chan PipelineReport)
	go func() {
		defer close(channel)
		var report PipelineReport
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines/%s/steps/%s/test_reports", repository, id, stepId))
		err := json.Unmarshal(response, &report)
		cobra.CheckErr(err)
		channel <- report
	}()
	return channel
}

func GetPipelineReportCases(repository string, id string, stepId string) <-chan PipelineReportCase {
	channel := make(chan PipelineReportCase)
	go func() {
		defer close(channel)
		var report BBPaginatedResponse[PipelineReportCase]
		// TODO pagelen is hardcoded, this should be changed if the number of tests are too big
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines/%s/steps/%s/test_reports/test_cases?pagelen=300", repository, id, stepId))
		err := json.Unmarshal(response, &report)
		cobra.CheckErr(err)
		for _, rep := range report.Values {
			channel <- rep
		}
	}()
	return channel
}

func RunPipeline(repository string, data RunPipelineRequestBody) Pipeline {
	content, err := json.Marshal(data)
	cobra.CheckErr(err)
	response := bbApiPost(fmt.Sprintf("repositories/%s/pipelines", repository), bytes.NewReader(content))

	var pipeline Pipeline
	err = json.Unmarshal(response, &pipeline)
	cobra.CheckErr(err)
	return pipeline
}

func StopPipeline(repository string, id string) {
	bbApiPost(fmt.Sprintf("repositories/%s/pipelines/%s/stopPipeline", repository, id), nil)
}

func GetPipelineVariables(repository string) <-chan []EnvironmentVariable {
	channel := make(chan []EnvironmentVariable)
	go func() {
		defer close(channel)
		var environmentResponse BBPaginatedResponse[EnvironmentVariable]
		response := bbApiGet(fmt.Sprintf("repositories/%s/pipelines_config/variables?pagelen=200", repository))
		err := json.Unmarshal(response, &environmentResponse)
		cobra.CheckErr(err)
		channel <- environmentResponse.Values
	}()
	return channel
}

func CreatePipelineVariable(repository string, key string, value string, secure bool) <-chan EnvironmentVariable {
	channel := make(chan EnvironmentVariable)
	go func() {
		defer close(channel)
		body := EnvironmentVariable{
			Key:     key,
			Value:   value,
			Secured: secure,
		}
		content, err := json.Marshal(body)
		cobra.CheckErr(err)
		response := bbApiPost(fmt.Sprintf("repositories/%s/pipelines_config/variables", repository), bytes.NewReader(content))
		var newVar EnvironmentVariable
		err = json.Unmarshal(response, &newVar)
		cobra.CheckErr(err)
		channel <- newVar
	}()
	return channel
}

func UpdatePipelineVariable(repository string, varUUID string, key string, value string, secure bool) <-chan EnvironmentVariable {
	channel := make(chan EnvironmentVariable)
	go func() {
		defer close(channel)
		body := EnvironmentVariable{
			Key:     key,
			Value:   value,
			Secured: secure,
		}
		content, err := json.Marshal(body)
		cobra.CheckErr(err)
		response := bbApiPut(fmt.Sprintf("repositories/%s/pipelines_config/variables/%s", repository, varUUID), bytes.NewReader(content))
		var newVar EnvironmentVariable
		err = json.Unmarshal(response, &newVar)
		cobra.CheckErr(err)
		channel <- newVar
	}()
	return channel
}

func DeletePipelineVariable(repository string, varUUID string) {
	bbApiDelete(fmt.Sprintf("repositories/%s/pipelines_config/variables/%s", repository, varUUID))
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
