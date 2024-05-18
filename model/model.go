package model

import "encoding/json"

const version = "v0.1.1"

type ResponseBaseModel struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version"`
}

type FileInfoModel struct {
	Size     uint   `json:"size" db:"size"`
	Mime     string `json:"mime" db:"mime"`
	FileID   string `json:"fileId" db:"fileid"`
	FileName string `json:"fileName" db:"fielname"`
}

type ResponseModel struct {
	ResponseBaseModel
	Data *FileInfoModel `json:"data"`
}

type ResponseFileModel struct {
	ResponseBaseModel
	Data []*FileInfoModel `json:"data"`
}

func ResponseFileModelJson(r *ResponseFileModel) []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}

	return data
}

func ResponseJson(r []*ResponseModel) []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}

	return data
}

func NewResponseModel() *ResponseModel {
	rm := new(ResponseModel)
	rm.Version = version
	rm.Data = &FileInfoModel{}
	return rm
}
