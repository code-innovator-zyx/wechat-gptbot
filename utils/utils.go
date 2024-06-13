package utils

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/sirupsen/logrus"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

var r *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
}

type Element interface {
	~int | ~string | ~float64 | ~float32
}

func EleInArray[T Element](slice []T, elem T) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

func HasPrefixes(s string, profiles []string) bool {
	for _, profile := range profiles {
		if strings.HasPrefix(s, profile) {
			return true
		}
	}
	return false
}

var client = http.Client{}

func DownloadImage(url string, reader *bytes.Buffer) error {
	response, err := client.Get(url)
	if err != nil {
		logrus.Infof("downloadImage failed, err=%+v", err)
		return err
	}
	_, err = io.Copy(reader, response.Body)
	if nil != err {
		return err
	}
	return nil
}

func CompressImage(url string, reader *bytes.Buffer) error {
	response, err := client.Get(url)
	if err != nil {
		logrus.Infof("downloadImage failed, err=%+v", err)
		return err
	}
	defer response.Body.Close()

	return GetImageSuffix(url)(response.Body, reader)
}
func GetImageSuffix(uri string) compress {
	u, _ := url.Parse(uri)
	// 获取文件名
	name := path.Base(u.Path)
	// 获取文件后缀
	switch path.Ext(name) {
	case ".webp":
		return compressWebp
	case ".jpeg":
		return compressJpeg
	default:
		return compressPng
	}
}

type compress func(io.Reader, *bytes.Buffer) error

func compressJpeg(reader io.Reader, buffer *bytes.Buffer) error {
	img, err := jpeg.Decode(reader)
	if nil != err {
		return err
	}
	// Compress image
	options := jpeg.Options{Quality: 50}
	err = jpeg.Encode(buffer, img, &options)
	if err != nil {
		logrus.Infof("compress image failed, err=%+v", err)
		return err
	}
	return nil
}

//func compressPng(reader io.Reader, buffer *bytes.Buffer) error {
//	// Png 无损压缩
//	io.Copy(buffer, reader) //
//	return nil
//}

func compressPng(reader io.Reader, buffer *bytes.Buffer) error {
	img, err := png.Decode(reader)
	if nil != err {
		return err
	}
	// Compress image
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression, // 1872330
	}
	err = encoder.Encode(buffer, img)
	if err != nil {
		logrus.Infof("compress image failed, err=%+v", err)
		return err
	}
	return nil
}

func compressWebp(reader io.Reader, buffer *bytes.Buffer) error {
	img, err := webp.Decode(reader)
	if nil != err {
		return err
	}
	// Compress image
	options := webp.Options{Quality: 50}
	err = webp.Encode(buffer, img, &options)
	if err != nil {
		logrus.Infof("compress image failed, err=%+v", err)
		return err
	}
	return nil
}

func BuildPersonalMessage(userName, content string) string {
	builder := strings.Builder{}
	builder.WriteString("【")
	builder.WriteString(userName)
	builder.WriteString("】:")
	builder.WriteString(content)
	return builder.String()
}

func BuildResponseMessage(userName, content, reply string) string {
	builder := strings.Builder{}
	//builder.WriteString("[")
	//builder.WriteString(userName)
	//builder.WriteString("]:")
	//builder.WriteString(content)
	//builder.WriteString("\n---------------------------------------------\n")
	builder.WriteString(reply)
	return builder.String()
}

func FakeIP() string {
	// 随便找的国内IP段：223.64.0.0 - 223.117.255.255
	return fmt.Sprintf("223.%d.%d.%d", r.Intn(54)+64, r.Intn(254), r.Intn(254))
}

func GetFromData(data map[string]string) io.Reader {
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	return strings.NewReader(formData.Encode())
}

func GetRandInt64(n int64) int64 {
	return r.Int63n(n)
}
