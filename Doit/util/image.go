package util

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

//图片转码
func DecodeImage(reader io.Reader) (image.Image, error) {
	//解码已注册格式编码的图像
	img, _, err := image.Decode(reader)
	return img, err
}
