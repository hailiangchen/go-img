package controller

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-img/cache"
	"go-img/conf"
	"go-img/db"
	"go-img/model"
	"go-img/util"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AdminController struct {
}

func (ac *AdminController) RegisterRouter(r *gin.Engine, handles ...gin.HandlerFunc) {
	gr := r.Group("/admin", handles...)
	{
		gr.POST("/upload", ac.Upload)
		gr.POST("/delete/:fileId", ac.Delete)
		gr.GET("/getall", ac.GetAll)
	}
}

func (ac *AdminController) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewResponse(false, "参数错误", nil))
		return
	}

	files := form.File["files"]
	results := make([]*model.Response, 0, len(files))
	for _, file := range files {
		response := model.NewResponse(false, "", &model.FileInfo{})
		response.Data.Size = uint32(file.Size)
		response.Data.FileName = file.Filename
		fileBytes, mime, err := verifyFile(file)
		response.Data.Mime = mime
		if err != nil {
			response.Message = err.Error()
			results = append(results, response)
			break
		}
		fileId, err := saveFile(fileBytes)
		if err != nil {
			response.Message = err.Error()
			results = append(results, response)
			break
		}
		response.Data.FileID = fileId
		response.Success = true
		results = append(results, response)

		db.Insert(response.Data)
	}

	c.JSON(http.StatusOK, results)
}

func (ac *AdminController) Delete(c *gin.Context) {
	fileId := c.Param("fileId")
	fileId = strings.Trim(fileId, "/")
	if len(fileId) == 0 || !util.IsMd5Str(fileId) {
		c.JSON(http.StatusBadRequest, model.NewResponse(false, "参数错误", nil))
		return
	}

	fileDir := util.Md5Path(fileId)
	file_path := filepath.Join(conf.AppConf.ImageConf.UploadPath, fileDir)

	if b, err := util.PathExists(file_path); !b {
		log.Println("os stat fialed:", err)
		c.JSON(http.StatusOK, model.NewResponse(false, "文件不存在", nil))
		return
	}

	err := os.RemoveAll(file_path)
	if err != nil {
		log.Println("remove all failed:", err)
		c.JSON(http.StatusOK, model.NewResponse(false, "删除失败", nil))
		return
	}

	if conf.AppConf.CacheConf.Enable {
		cache.NewRedisCache().Del(fileId)
	}

	db.Delete(fileId)

	c.JSON(http.StatusOK, model.NewResponse(true, "", nil))
}

type getAllInput struct {
	Page     uint `form:"page"`
	PageSize uint `form:"pageSize"`
}

func (ac *AdminController) GetAll(c *gin.Context) {
	var getAllInput getAllInput
	if err := c.ShouldBindQuery(&getAllInput); err != nil {
		c.JSON(http.StatusBadRequest, "error")
		return
	}

	var response = new(model.ResponseFiles)
	response.Success = true
	data, count, err := db.GetAll(getAllInput.Page, getAllInput.PageSize)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
	}
	response.Total = count
	response.Data = data

	c.JSON(http.StatusOK, response)
}

func verifyFile(file *multipart.FileHeader) (*[]byte, string, error) {
	fileFd, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("打开文件失败")
	}
	defer fileFd.Close()

	fileBytes := make([]byte, file.Size)
	_, err = bufio.NewReader(fileFd).Read(fileBytes)
	if err != nil && err != io.EOF {
		log.Println("读取文件内容失败：", err)
		return nil, "", fmt.Errorf("读取文件内容错误")
	}

	mime := http.DetectContentType(fileBytes[:16])
	if !util.IsContains(mime) {
		return nil, mime, fmt.Errorf("文件类型不允许上传")
	}

	return &fileBytes, mime, nil
}

func saveFile(fileBytes *[]byte) (string, error) {
	fileMd5 := fmt.Sprintf("%x", md5.Sum(*fileBytes))
	savePath := filepath.Join(conf.AppConf.ImageConf.UploadPath, util.Md5Path(fileMd5), "0_0")
	err := util.SaveFile(fileBytes, savePath)
	return fileMd5, err
}
