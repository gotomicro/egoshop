package image

import (
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
)

type ReqImageList struct {
	trans.ReqPage
	mysql.Image
}

type ReqCreate struct {
	Image   string `json:"image"`
	IsSave  int    `json:"isSave"`
	Space   string `json:"space"`
	OssType int    `json:"OssType"`
}

type ReqAdd struct {
	Image  string `json:"image"`
	IsSave int    `json:"is_save"`
	Space  string `json:"space"`
}

type RespImage struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}
