package utils

import (
	"go-image/config"
	"log"
	"regexp"
	"strings"
)

var regexpURLParse *regexp.Regexp
var ImageTypes []string
var ImagePath string
var AdminIPs string

func init() {
	var err error

	ImagePath = config.GetSetting("image.path")

	if len(ImagePath) == 0 {
		ImagePath = "image/"
	} else {
		if ImagePath[len(ImagePath):] != "/" {
			ImagePath = ImagePath + "/"
		}
	}

	regexpURLParse, err = regexp.Compile("[a-z0-9]{32}")
	if err != nil {
		log.Println("regexpUrlParse:", err)
	}

	ImageTypes = strings.Split(config.GetSetting("image.type"), ",")
	AdminIPs = config.GetSetting("server.admin_ips")
}
