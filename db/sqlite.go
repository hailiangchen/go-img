package db

import (
	"database/sql"
	"go-img/model"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var dbClient *sql.DB
var err error

func init() {
	dbClient, err = sql.Open("sqlite3", "db/store.db?loc=Asia/Shanghai")
	if err != nil {
		log.Fatalln("open db file failed", err)
	}
}

func Insert(fileInfo *model.FileInfo) {
	sqlStr := "Insert into infos(fileid,mime,size,filename) values(?,?,?,?)"
	_, err := dbClient.Exec(sqlStr, fileInfo.FileID, fileInfo.Mime, fileInfo.Size, fileInfo.FileName)
	if err != nil {
		log.Println("fileInfo insert into failed: ", fileInfo.FileID, err)
	}
}

func GetAll(page, size uint) ([]*model.FileInfo, uint, error) {
	sqlStr := "select fileid,filename,mime,size,creation from infos limit ?,?"
	count_sql := "select count(*) as count from infos"
	var count uint = 0
	err := dbClient.QueryRow(count_sql).Scan(&count)
	if err != nil {
		log.Println("select all list failed", err)
		return nil, count, err
	}
	rows, err := dbClient.Query(sqlStr, page-1, size)
	if err != nil {
		log.Println("select all list failed", err)
		return nil, count, err
	}
	defer rows.Close()

	var result = make([]*model.FileInfo, 0, size)
	for rows.Next() {
		var fileInfo = &model.FileInfo{}
		err := rows.Scan(&fileInfo.FileID, &fileInfo.FileName, &fileInfo.Mime, &fileInfo.Size, &fileInfo.Creation)
		if err != nil {
			log.Println("rows scan failed", err)
			return nil, count, err
		}
		result = append(result, fileInfo)
	}

	return result, count, nil
}

func Delete(fileId string) {
	sqlStr := "delete from infos where fileid=?"
	_, err := dbClient.Exec(sqlStr, fileId)
	if err != nil {
		log.Println("delete fileInfo failed", fileId)
	}
}
