package model

type Response struct {
	Result  interface{} `json:"result"`
	Message *string     `json:"message"`
}
