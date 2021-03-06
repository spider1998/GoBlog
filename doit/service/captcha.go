package service

import (
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/resource"
	"bytes"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/pkg/errors"
	"image/color"
	"image/jpeg"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var Captcha = &CaptchaService{}

type CaptchaService struct{}

func file2Bytes(filename string) ([]byte, error) {

	// File
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// FileInfo:
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// []byte
	data := make([]byte, stats.Size())
	count, err := file.Read(data)
	if err != nil {
		return nil, err
	}
	fmt.Printf("read file %s len: %d \n", filename, count)
	return data, nil
}

func (s *CaptchaService) Generate() (token string, image []byte, err error) {
	token = RandString(32)
	captchaKey := s.captchaKey(token)
	c := captcha.New()
	fmt.Println(resource.FontBox.Path)
	data, err := file2Bytes("doit/resource/fonts/comic.ttf")
	if err != nil {
		return
	}
	err = c.AddFontFromBytes(data)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	c.SetBkgColor(color.RGBA{0xc8, 0xe1, 0xff, 1})
	value := RandString(4)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, c.CreateCustom(value), nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	err = app.Redis.Cmd("SET", captchaKey, strings.ToLower(value), "EX", 300).Err
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return token, buf.Bytes(), nil
}

func (s *CaptchaService) Validate(token, value string) error {
	key := s.captchaKey(token)
	v, err := app.Redis.Cmd("GET", key).Str()
	if err != nil {
		if err == redis.ErrRespNil {
			err = code.New(http.StatusBadRequest, code.CodeInvalidCaptcha)
			return err
		}
		err = errors.WithStack(err)
		return err
	}
	if v != strings.ToLower(value) {
		err = code.New(http.StatusBadRequest, code.CodeInvalidCaptcha)
		return err
	}
	err = app.Redis.Cmd("DEL", key).Err
	if err != nil {
		app.Logger.Error().Err(err).Msg("fail to delete captcha key in redis.")
	}
	return nil
}

func (s *CaptchaService) captchaKey(token string) string {
	return app.System + ":captcha:" + token
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
