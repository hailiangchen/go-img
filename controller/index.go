package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-img/cache"
	"go-img/conf"
	"go-img/image"
	"go-img/model"
	"go-img/util"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type IndexController struct {
}

func (ic *IndexController) RegisterRouter(r *gin.Engine, handles ...gin.HandlerFunc) {
	gr := r.Group("", handles...)
	{
		gr.GET("/:fileId", ic.Index)
	}
}

func (ic *IndexController) Index(c *gin.Context) {
	//设置初始默认值
	req := &model.Request{Format: "jpeg", X: -1, Y: -1, Proportion: -1, Quality: 75}
	req.FileId = c.Param("fileId")
	if err := c.ShouldBindQuery(req); err != nil || len(req.FileId) == 0 || !util.IsMd5Str(req.FileId) {
		c.JSON(http.StatusBadRequest, "params error")
		return
	}

	c.Header("Cache-Control", "max-age=3600")
	if req.Download {
		c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s.%s", req.FileId, req.Format))
	}

	dirPath := filepath.Join(conf.AppConf.ImageConf.UploadPath, util.Md5Path(req.FileId))
	if fileBytes := readFile(req, dirPath); fileBytes != nil {
		c.Data(http.StatusOK, http.DetectContentType(fileBytes[:16]), fileBytes)
		return
	}

	newFileBytes, err := image.Convert(req, filepath.Join(dirPath, "0_0"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	util.SaveFile(&newFileBytes, filepath.Join(dirPath, req.FormatFileName()))
	if conf.AppConf.CacheConf.Enable {
		cache.NewRedisCache().Set(fmt.Sprintf("%s:%s", req.FileId, req.FormatFileName()), newFileBytes)
	}

	c.Data(http.StatusOK, http.DetectContentType(newFileBytes[:16]), newFileBytes)
}

func readFile(req *model.Request, dirPath string) []byte {
	//获取原图
	if req.Proportion == 0 && req.Width == 0 && req.Height == 0 {
		file, err := os.Open(filepath.Join(dirPath, "0_0"))
		if err != nil {
			log.Println("open original image failed:", err)
			return nil
		}
		defer file.Close()
		buffer := bytes.NewBuffer(make([]byte, 0))
		io.Copy(buffer, file)

		return buffer.Bytes()
	}

	//从缓存获取
	if conf.AppConf.CacheConf.Enable {
		fileByte, _ := cache.NewRedisCache().Get(fmt.Sprintf("%s:%s", req.FileId, req.FormatFileName()))
		return fileByte
	}

	//从本地获取
	if ok, _ := util.PathExists(filepath.Join(dirPath, req.FormatFileName())); ok {
		file, err := os.Open(filepath.Join(dirPath, req.FormatFileName()))
		if err != nil {
			log.Println("open image failed:", err)
			return nil
		}
		defer file.Close()
		buffer := bytes.NewBuffer(make([]byte, 0))
		io.Copy(buffer, file)
		if conf.AppConf.CacheConf.Enable {
			cache.NewRedisCache().Set(fmt.Sprintf("%s:%s", req.FileId, req.FormatFileName()), buffer.Bytes())
		}

		return buffer.Bytes()
	}

	return nil
}
