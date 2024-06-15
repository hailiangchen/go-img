package util

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"go-img/conf"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func IsContains(fileType string) bool {
	for _, v := range conf.AppConf.ImageConf.Exts {
		if strings.Contains(fileType, strings.ToLower(v)) {
			return true
		}
	}
	return false
}

func GetHash(r *bufio.Reader) string {
	h := md5.New()

	if _, err := io.Copy(h, r); err != nil {
		log.Println(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetBigHash(f *os.File) string {
	info, _ := f.Stat()
	fileSize := info.Size()
	var fileChunk int64 = 8192
	blocks := int64(math.Ceil(float64(fileSize / fileChunk)))

	hash := md5.New()
	for i := int64(0); i < blocks; i++ {
		blocksize := int(math.Min(float64(fileChunk), float64(fileSize-i*fileChunk)))
		buf := make([]byte, blocksize)

		f.Read(buf)
		hash.Write(buf)
		//io.WriteString(hash, string(buf)) // append into the hash
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Md5Path(md5Str string) string {
	first, err := strconv.ParseUint(md5Str[:3], 16, 32)
	if err != nil {
		log.Println(err)
		return ""
	}

	second, err := strconv.ParseUint(md5Str[3:6], 16, 32)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("%d/%d/%s", first/4, second/4, md5Str)
}

func IsAllow(ip string) bool {
	if conf.AppConf.HttpConf.AdminIPS == "*" {
		return true
	} else {
		return strings.Contains(conf.AppConf.HttpConf.AdminIPS, ip)
	}
}

func IsMd5Str(str string) bool {
	regexpURLParse, err := regexp.Compile("[a-z0-9]{32}")
	if err != nil {
		log.Println("正则表达式错误：", err)
		return false
	}

	return regexpURLParse.MatchString(str)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func SaveFile(byte *[]byte, savePath string) error {
	err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("创建目录失败")
	}

	file, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE, 0664)
	defer file.Close()
	file.Write(*byte)

	return nil
}
