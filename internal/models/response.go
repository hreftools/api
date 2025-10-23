package models

type ResponseSuccess struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type ResponseError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
