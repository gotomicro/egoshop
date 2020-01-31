package imagex

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/goecology/muses/pkg/system"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	Version           = "1.7"
	StoreLocal string = "local"
	StoreOss   string = "oss"
)

//操作图片显示
//如果用的是oss存储，这style是avatar、cover可选项
func ShowImg(img string, style ...string) (url string) {
	if strings.HasPrefix(img, "https://") || strings.HasPrefix(img, "http://") {
		return img
	}
	img = strings.TrimLeft(img, "./")
	switch viper.GetString("app.storeType") {
	case StoreOss:
		s := ""
		if len(style) > 0 && strings.TrimSpace(style[0]) != "" {
			s = "/" + style[0]
		}
		url = img + s
	case StoreLocal:
		url = img
	}
	// 说明没有图片，给个默认图片
	if url == "/" {
		url = "/static/images/book.png"
	}
	// 适应小程序
	url = viper.GetString("oss.cdnName") + url
	return
}

func ShowImgArr(imgs []string, style ...string) (urlArr []string) {
	urlArr = make([]string, 0)
	for _, img := range imgs {
		urlArr = append(urlArr, ShowImg(img, style...))
	}
	return
}

func FilterImgArr(imgs []string) (urlArr []string) {
	urlArr = make([]string, 0)
	for _, img := range imgs {
		// 如果存在cdn，截取掉对应信息
		if strings.HasPrefix(img, viper.GetString("oss.cdnName")) {
			img = strings.ReplaceAll(img, viper.GetString("oss.cdnName"), "")
			// 如果url里存在/,但斜线后面无.jpg，那么需要裁剪
			lastIndex := strings.LastIndex(img, "/")
			if !strings.HasPrefix(img[lastIndex:], ".") {
				img = img[:lastIndex]
			}
		}
		urlArr = append(urlArr, img)
	}
	return
}

func FilterOneImg(img string) string {
	// 如果存在cdn，截取掉对应信息
	if strings.HasPrefix(img, viper.GetString("oss.cdnName")) {
		img = strings.ReplaceAll(img, viper.GetString("oss.cdnName"), "")
		// 如果url里存在/,但斜线后面无.jpg，那么需要裁剪
		lastIndex := strings.LastIndex(img, "/")
		if !strings.HasPrefix(img[lastIndex:], ".") {
			img = img[:lastIndex]
		}
	}
	return img
}

// Substr returns the substr from start to length.
func Substr(s string, length int) string {
	bt := []rune(s)
	start := 0
	dot := false

	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
		dot = true
	}

	str := string(bt[start:end])
	if dot {
		str = str + "..."
	}
	return str
}

func GeneratePath(space string) (string, string) {
	month := time.Now().Format("200601")
	return fmt.Sprintf("%s/%s/%s/", viper.GetString("image.path"), space, month), month
}

func GenerateUniqueMd5() string {
	date := time.Now().Format("20060102150405")
	uniqueID := GenerateUniqueID()
	sno := date + system.RunInfo.HostName + string(system.RunInfo.Pid) + uniqueID

	return fmt.Sprintf("%x", md5.Sum([]byte(sno)))
}

func GenerateUniqueID() string {
	b := make([]byte, 16)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
