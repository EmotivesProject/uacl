package model

type AutologinToken struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	AutologinToken string `json:"autologin_token"`
	Site           string `json:"site"`
}
