package db

import (
	"database/sql"
	"go-image/model"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var dbClient *sql.DB
var err error

func init() {
	dbClient, err = sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatalln("open db file failed")
	}
}

func Insert(fileInfo *model.FileInfoModel) {
	sqlStr := "Insert into infos(fileid,mime,size,filename) values(?,?,?,?)"

	_, err := dbClient.Exec(sqlStr, fileInfo.FileID, fileInfo.Mime, fileInfo.Size, fileInfo.FileName)
	if err != nil {
		log.Println("fileInfo insert into failed", fileInfo.FileID)
	}
}

func GetAll() ([]*model.FileInfoModel, error) {
	sqlStr := "select fileid,filename,mime,size from infos"
	rows, err := dbClient.Query(sqlStr)
	if err != nil {
		log.Println("select all list failed", err)
		return nil, err
	}
	defer rows.Close()

	var result = make([]*model.FileInfoModel, 0)
	for rows.Next() {
		var fileInfo = &model.FileInfoModel{}
		err := rows.Scan(&fileInfo.FileID, &fileInfo.FileName, &fileInfo.Mime, &fileInfo.Size)
		if err != nil {
			log.Println("rows scan failed", err)
			return nil, err
		}
		result = append(result, fileInfo)
	}

	return result, nil
}

func Delete(fileId string) {
	sqlStr := "delete from infos where fileid=?"
	_, err := dbClient.Exec(sqlStr, fileId)
	if err != nil {
		log.Println("delete fileInfo failed", fileId)
	}
}
