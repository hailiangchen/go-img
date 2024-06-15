package model

type FileInfo struct {
	Size     uint32 `json:"size" db:"size"`
	Mime     string `json:"mime" db:"mime"`
	FileID   string `json:"fileId" db:"fileid"`
	FileName string `json:"fileName" db:"fielname"`
	Creation string `json:"creation"`
}
