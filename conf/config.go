package conf

import (
	"gopkg.in/ini.v1"
	"log"
	"strings"
)

type appConf struct {
	HttpConf  httpConf  `ini:"http"`
	ImageConf imageConf `ini:"image"`
	CacheConf cacheConf `ini:"cache"`
}

type httpConf struct {
	Address  string `ini:"addr"`
	AdminIPS string `ini:"admin_ips"`
	RunMode  string `ini:"run_mode"`
}

type imageConf struct {
	Ext        string `ini:"ext"`
	UploadPath string `ini:"upload_path"`
	Exts       []string
}

type cacheConf struct {
	Enable     bool   `ini:"enable"`
	Address    string `ini:"addr"`
	Password   string `ini:"password"`
	MaxCache   uint32 `ini:"max_cache"`
	ExpireTime uint32 `ini:"expire_time"`
}

var AppConf = new(appConf)

func init() {
	err := ini.MapTo(AppConf, "conf/app.conf")
	if err != nil {
		log.Fatalln("Failed to read config file:", err)
	}

	AppConf.ImageConf.Exts = strings.Split(AppConf.ImageConf.Ext, ",")
}
