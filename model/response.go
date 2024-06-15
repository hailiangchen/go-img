package model

const version = "v1.0.0"

type Response struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Version string    `json:"version"`
	Data    *FileInfo `json:"data"`
}

type ResponseFiles struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Total   uint        `json:"total"`
	Data    []*FileInfo `json:"data"`
}

func NewResponse(success bool, message string, fileInfo *FileInfo) *Response {
	response := new(Response)
	response.Version = version
	response.Success = success
	response.Message = message
	response.Data = fileInfo
	return response
}
