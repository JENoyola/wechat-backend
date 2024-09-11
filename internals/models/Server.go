package models

type ServerResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	DATA    interface{} `json:"data"`
}
