package controller

import (
	"fmt"
	"go-image/cache"
	"go-image/db"
	"go-image/filehandler"
	"go-image/model"
	"go-image/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Index 处理首页路径
func Index(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.String()
	if urlStr == "/favicon.ico" {
		return
	}

	var req = new(model.Goimg_req_t)

	parse, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	path := utils.ParseUrlPath(parse.Path[1:])
	if path == "" {
		http.NotFound(w, r)
		return
	}

	model.ParamHandler(req, r)

	dirPath := utils.ImagePath + path
	sourceFilePath := dirPath + "/0_0"
	md5Str := parse.Path[1:]
	var cacheKey string

	//参数d=1时，直接下载文件。
	if req.Download == 1 {
		w.Header().Set("Content-Disposition", "attachment;filename="+md5Str+"."+req.Format)
	}

	w.Header().Set("Cache-Control", "max-age=3600") //强制浏览器缓存
	//w.Header().Set("Expires", time.Now().Add(10*time.Hour).UTC().Format(http.TimeFormat))
	if req.P == 1 {
		file, err := os.Open(sourceFilePath)
		if err != nil {
			log.Println(err)
			http.Error(w, "未找到文件", http.StatusNotFound)
			return
		}
		defer file.Close()
		io.Copy(w, file)
		return
	}

	//从缓存读取
	if cache.IsCache {
		cacheKey = fmt.Sprintf("%s:%d_%d_g%d_r%.f_p%d_x%d_y%d_q%d.%s", md5Str, req.Width, req.Height, req.Grayscale, req.Rotate, req.P, req.X, req.Y, req.Quality, req.Format)
		cacheValue := cache.Get(cacheKey)
		if *cacheValue != nil {
			w.Write(*cacheValue)
			return
		}
	}

	//从硬盘读取
	filePath := fmt.Sprintf("%s/%d_%d_g%d_r%.f_p%d_x%d_y%d_q%d.%s", dirPath, req.Width, req.Height, req.Grayscale, req.Rotate, req.P, req.X, req.Y, req.Quality, req.Format)
	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		if cache.IsCache {
			cache.Set(cacheKey, b)
		}
		w.Write(b)
		return
	}

	if _, err = os.Stat(sourceFilePath); err != nil && !os.IsExist(err) {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}

	b, err := filehandler.ResizeImage(sourceFilePath, req, filePath)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cache.IsCache {
		cache.Set(cacheKey, *b)
	}

	w.Write(*b)
}

// Uploads upload files function.
func Upload(w http.ResponseWriter, r *http.Request) {
	//允许跨域
	//w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseMultipartForm(1024 << 14)
	if r.MultipartForm == nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["files"]
	var response = make([]*model.ResponseModel, 0)
	for i := 0; i < len(files); i++ {
		resp := model.NewResponseModel()
		// fileInfo := new(model.FileInfoModel)
		// resp.Data = fileInfo

		file, err := files[i].Open()
		if err != nil {
			resp.Success = false
			resp.Message = "读取file数据失败"
			response = append(response, resp)
			break
		}
		defer file.Close()

		resp.Data.FileName = files[i].Filename

		b, err := ioutil.ReadAll(file)
		if err != nil {
			resp.Success = false
			resp.Message = "读取file数据失败"
			response = append(response, resp)
			break
		}

		resp.Data.Mime = http.DetectContentType(b[:512])
		resp.Data.Size = uint(len(b))

		if !utils.IsType(resp.Data.Mime) {
			resp.Success = false
			resp.Message = "图片类型不符合"
			response = append(response, resp)
			break
		}

		md5Str := filehandler.GetHash(&b)
		md5Path := utils.SavePath(md5Str)

		file.Seek(0, 0)

		dirPath := utils.ImagePath + md5Path + "/"

		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			resp.Success = false
			resp.Message = "创建目录失败"
			response = append(response, resp)
			break
		}

		err = ioutil.WriteFile(dirPath+"0_0", b, 0660)
		if err != nil {
			resp.Success = false
			resp.Message = err.Error()
			response = append(response, resp)
			break
		}

		// err = filehandler.CompressionImage(b, dirPath+"0_0", 75, resp.Data)
		// if err != nil {
		// 	resp.Success = false
		// 	resp.Message = err.Error()
		// 	response = append(response, resp)
		// 	break
		// }

		resp.Success = true
		resp.Message = "OK"
		resp.Data.FileID = md5Str
		response = append(response, resp)

		db.Insert(resp.Data)
	}

	w.Write(model.ResponseJson(response))
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	var resp = &model.ResponseFileModel{}
	resp.Success = true
	res, err := db.GetAll()
	if err != nil {
		resp.Success = false
		resp.Message = "获取信息失败"
	}
	resp.Data = res
	w.Write(model.ResponseFileModelJson(resp))
}

func Delete(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.String()
	if urlStr == "/favicon.ico" {
		return
	}

	parse, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}
	md5Str := parse.Path[strings.LastIndex(parse.Path, "/")+1:]

	if !utils.IsMd5Str(md5Str) {
		http.NotFound(w, r)
		return
	}

	md5Path := utils.SavePath(md5Str)
	if _, err = os.Stat(utils.ImagePath + md5Path); err == nil || os.IsExist(err) {
		err = os.RemoveAll(utils.ImagePath + md5Path)
		if err != nil {
			http.Error(w, "删除失败", http.StatusInternalServerError)
			return
		}

		if cache.IsCache {
			cache.Del(md5Str)
		}

		db.Delete(md5Str)

		fmt.Fprintln(w, "ok")
	} else {
		log.Println(err)
		fmt.Fprintln(w, "文件不存在")
	}
}

func responseError(w http.ResponseWriter, code int) {
	html := fmt.Sprintf("<html><body><h1>%s</h1></body></html>", model.StatusText(code))
	w.WriteHeader(code)
	fmt.Fprintln(w, html)
}
