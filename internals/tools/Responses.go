package tools

import "wechat-back/internals/models"

func FormatSuccessResponse(data any, code int, message string) models.ServerResponse {
	var res models.ServerResponse
	res.Code = code
	res.DATA = data
	res.Error = false
	res.Message = message

	return res
}

func FormatErrResponse(code int, err error) models.ServerResponse {
	var res models.ServerResponse

	res.Code = code
	res.Error = true
	res.Message = err.Error()

	return res
}

func FormatCustomErrResponse(message string, code int) models.ServerResponse {
	var res models.ServerResponse

	res.Code = code
	res.Error = true
	res.Message = message

	return res
}
