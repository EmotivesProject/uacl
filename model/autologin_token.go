package model

type AutologinToken struct {
	Username       string `json:"username"`
	AutologinToken string `json:"autologin_token"`
	Site           string `json:"site"`
}
