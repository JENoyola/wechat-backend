package tools

import "wechat-back/internals/models"

func FormatResponse(code int, err error) models.ServerResponse {
	var res models.ServerResponse

	res.Code = code
	res.Error = true
	res.Message = err.Error()

	return res
}

func FormatCustomResponse(message string, code int) models.ServerResponse {
	var res models.ServerResponse

	res.Code = code
	res.Error = true
	res.Message = message

	return res
}
