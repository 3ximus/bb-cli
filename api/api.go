package api

import (
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func get(content string, spinner *spinner.Spinner) []byte {
	token := viper.GetString("token")
	username := viper.GetString("username")
	api_endpoint := viper.GetString("api")

	client := &http.Client{}

	url := fmt.Sprintf("%s/%s", api_endpoint, content)
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
func api_get(content string) []byte {
	s := spinner.New(
		spinner.CharSets[14],
		100*time.Millisecond,
		spinner.WithWriter(os.Stderr),
		spinner.WithSuffix(" Sending request..."),
		spinner.WithColor("fgHiBlue"),
	)
	s.Start()
	response := get("user", s)
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

func GetPr() {
}
