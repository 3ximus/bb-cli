package api

type User struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	AccountId   string `json:"account_id"`
	Links       struct {
		Html struct {
			Href string
		}
	}
}
