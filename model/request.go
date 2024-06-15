package model

import (
	"fmt"
	"strings"
)

type Request struct {
	Grayscale  bool    `form:"g"` //灰阶
	Rotate     float64 `form:"r"` //旋转
	Width      uint    `form:"w"`
	Height     uint    `form:"h"`
	Quality    uint    `form:"q"` //图片质量
	X          int     `form:"x"`
	Y          int     `form:"y"`
	Proportion int     `form:"p"` //分辨率
	Download   bool    `form:"d"` //下载图片
	Format     string  `form:"f"` //图片格式
	FileId     string
}

func (r *Request) FormatFileName() string {
	return fmt.Sprintf("%d_%d_g%d_r%.f_p%d_x%d_y%d_q%d.%s", r.Width, r.Height, BtoI(r.Grayscale), r.Rotate, r.Proportion, r.X, r.Y, r.Quality, r.GetFormat())
}

func (r *Request) GetFormat() string {
	switch strings.ToLower(r.Format) {
	case "jpeg", "jpg":
		r.Format = "jpeg"
	case "png":
		r.Format = "png"
	case "gif":
		r.Format = "gif"
	case "webp":
		r.Format = "webp"
	default:
		r.Format = "jpeg"
	}

	return r.Format
}

func BtoI(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
