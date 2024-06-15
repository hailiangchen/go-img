package image

import (
	"fmt"
	"go-img/model"
	"gopkg.in/gographics/imagick.v3/imagick"
	"log"
	"math"
)

func init() {
	imagick.Initialize()
}

type imageHandler func(*imagick.MagickWand) (*imagick.MagickWand, error)

func Convert(req *model.Request, imagePath string) ([]byte, error) {
	var err error
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImage(imagePath)
	if err != nil {
		log.Println("read image failed:", err)
		return nil, fmt.Errorf("读取文件失败")
	}

	mw.ResetIterator()

	mw, err = Handler(mw, Verify(), CropAndResize(req.Width, req.Height, req.X, req.Y, req.Proportion),
		Rotate(req.Rotate), Gray(req.Grayscale), Quality(req.Quality), Format(req.GetFormat()))
	if err != nil {
		return nil, err
	}

	mw.StripImage()

	return mw.GetImageBlob(), nil
}

func Handler(mw *imagick.MagickWand, handlers ...imageHandler) (*imagick.MagickWand, error) {
	var err error
	for _, handler := range handlers {
		mw, err = handler(mw)
		if err != nil {
			return mw, err
		}
	}
	return mw, nil
}

func Verify() imageHandler {
	return func(mw *imagick.MagickWand) (*imagick.MagickWand, error) {
		imageType := mw.GetImageColorspace()
		if imageType == imagick.COLORSPACE_CMYK {
			return mw, fmt.Errorf("该图像为CMYK，暂不支持处理")
		}

		ret := mw.GetImageOrientation()
		if ret > imagick.ORIENTATION_UNDEFINED {
			err := mw.AutoOrientImage()
			if err != nil {
				log.Println("orientation auto orientat faile", err)
				return mw, fmt.Errorf("orientation auto orientat faild")
			}
		}

		return mw, nil
	}
}

func CropAndResize(w, h uint, x, y, p int) imageHandler {
	return func(mw *imagick.MagickWand) (*imagick.MagickWand, error) {
		//都不为0
		if !(w == 0 && h == 0) {
			if x == -1 && y == -1 {
				return Proportion(mw, w, h, p)
			} else {
				return Crop(mw, w, h, x, y)
			}
		}

		return mw, nil
	}
}

// Rotate 旋转图像
func Rotate(rotate float64) imageHandler {
	return func(mw *imagick.MagickWand) (*imagick.MagickWand, error) {
		if rotate != 0 {
			pw := imagick.NewPixelWand()
			defer pw.Destroy()

			pw.SetColor("withe")
			err := mw.RotateImage(pw, rotate)
			if err != nil {
				log.Println("rotote image failed:", err)
				return mw, fmt.Errorf("旋转图像失败")
			}
		}

		return mw, nil
	}
}

// Gray 灰阶图像
func Gray(isGray bool) imageHandler {
	return func(mw *imagick.MagickWand) (*imagick.MagickWand, error) {
		if isGray {
			err := mw.SetImageType(imagick.IMAGE_TYPE_GRAYSCALE)
			if err != nil {
				log.Println("gray image failed:", err)
				return mw, fmt.Errorf("图像灰阶失败")
			}
		}

		return mw, nil
	}
}

// Quality 图像质量
func Quality(quality uint) imageHandler {
	return func(mw *imagick.MagickWand) (*imagick.MagickWand, error) {
		currentQuality := mw.GetImageCompressionQuality()
		if currentQuality == 0 || currentQuality > quality {
			err := mw.SetImageCompressionQuality(quality)
			if err != nil {
				log.Println("Compression quality image failed:", err)
				return mw, fmt.Errorf("图像压缩失败")
			}
		}

		return mw, nil
	}
}

// Format 转换格式
func Format(format string) imageHandler {
	return func(mw *imagick.MagickWand) (*imagick.MagickWand, error) {
		err := mw.SetImageFormat(format)
		if err != nil {
			log.Println("set image format failed:", err)
			return mw, fmt.Errorf("设置图像格式失败")
		}

		return mw, nil
	}
}

func Proportion(mw *imagick.MagickWand, w, h uint, p int) (*imagick.MagickWand, error) {
	img_w := mw.GetImageWidth()
	img_h := mw.GetImageHeight()

	//按照给定宽度和高度处理
	if p == 0 {
		if w == 0 || h == 0 {
			return mw, fmt.Errorf("p=0时,需要指定宽度和高度")
		}
		return Resize(mw, w, h)

		//以中心部分裁切
	} else if p == 2 {
		x := (img_w - w) / 2
		y := (img_h - h) / 2
		return Crop(mw, w, h, int(x), int(y))

		//按照图片百分比缩放
	} else if p == 3 {
		if w == 0 || h == 0 {
			var rate uint
			if w > 0 {
				rate = w
			} else {
				rate = h
			}

			w = uint(math.Round(float64(img_w) * float64(rate) / 100))
			h = uint(math.Round(float64(img_h) * float64(rate) / 100))

			return Resize(mw, w, h)
		} else {
			w = uint(math.Round(float64(img_w) * float64(w) / 100))
			h = uint(math.Round(float64(img_h) * float64(h) / 100))

			return Resize(mw, w, h)
		}

		//按照尺寸等比例缩放
	} else {
		if w == 0 || h == 0 {
			if w > 0 {
				h = uint(math.Floor((float64(w) / float64(img_w)) * float64(img_h)))
			} else {
				w = uint(math.Floor((float64(h) / float64(img_h)) * float64(img_w)))
			}

			return Resize(mw, w, h)
		} else {
			var rate float64
			rate_w := float64(w) / float64(img_w)
			rate_h := float64(h) / float64(img_h)

			if rate_w < rate_h {
				rate = rate_w
			} else {
				rate = rate_h
			}

			w = uint(math.Round(rate * float64(img_w)))
			h = uint(math.Round(rate * float64(img_h)))

			return Resize(mw, w, h)
		}
	}

	return nil, nil
}

// Resize 缩放
func Resize(mw *imagick.MagickWand, w, h uint) (*imagick.MagickWand, error) {
	img_w := mw.GetImageWidth()
	img_h := mw.GetImageHeight()
	if w > img_w {
		w = img_w
	}
	if h > img_h {
		h = img_h
	}
	err := mw.ResizeImage(w, h, imagick.FILTER_LANCZOS)
	if err != nil {
		log.Println("resize image failed:", err)
		return mw, fmt.Errorf("图像缩放失败")
	}
	return mw, nil
}

// CropHandler 裁切
func Crop(mw *imagick.MagickWand, w, h uint, x, y int) (*imagick.MagickWand, error) {
	img_w := mw.GetImageWidth()
	img_h := mw.GetImageHeight()
	if uint(x) >= img_w || uint(y) >= img_h {
		return mw, fmt.Errorf("图像坐标错误")
	}

	if x < 0 {
		x = 0
	}

	if y < 0 {
		y = 0
	}

	if w == 0 || img_w < uint(x)+w {
		w = img_w - uint(x)
	}

	if h == 0 || img_h < uint(y)+h {
		h = img_h - uint(y)
	}

	err := mw.CropImage(w, h, x, y)
	if err != nil {
		log.Println("crop failed:", err)
		return mw, fmt.Errorf("图像裁切失败")
	}

	return mw, nil
}
