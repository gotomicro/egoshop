package image

import (
	"encoding/base64"
	"errors"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/service"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/conf"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/mdw"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/goecology/egoshop/appgo/pkg/imagex"
	"go.uber.org/zap"
)

func Create(c *gin.Context) {
	// base64.StdEncoding.DecodeString(datasource)
	reqModel := ReqCreate{}
	if err := c.Bind(&reqModel); err != nil {
		mus.Logger.Error(err.Error(), zap.Int("code", 1))
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	if reqModel.Space == "" {
		reqModel.Space = "default"
	}

	if reqModel.OssType == 0 {
		reqModel.OssType = 1
	}

	ext := ".jpg"
	dataType := reqModel.Image[strings.IndexByte(reqModel.Image, ':')+1 : strings.IndexByte(reqModel.Image, ';')]
	b64data := reqModel.Image[strings.IndexByte(reqModel.Image, ',')+1:]

	b64dataDecode, _ := base64.StdEncoding.DecodeString(b64data) // 成图片文件并把文件写入到buffer

	fileName := imagex.GenerateUniqueMd5()
	rootPath, _ := imagex.GeneratePath(reqModel.Space)

	filePath := filepath.Join(rootPath, fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err := ioutil.WriteFile(filePath, b64dataDecode, 0666) // buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	prefix := service.Oss.OssPrefix(reqModel.OssType)
	if prefix == "" {
		base.JSONErr(c, 10003, errors.New("not exist type"))
		return
	}

	key := service.Oss.Key(prefix)
	err = service.Oss.PutObj(viper.GetString("oss.bucket"), key, filePath)
	if err != nil {
		base.JSONErr(c, 10004, err)
		return
	}

	url, err := service.Oss.GetObjURL(viper.GetString("oss.bucket"), key)
	if err != nil {
		base.JSONErr(c, 10005, err)
		return
	}

	var faImage = mysql.Image{
		CreatedBy: mdw.DefaultContextUser(c).Id,
		UpdatedBy: mdw.DefaultContextUser(c).Id,
		Name:      fileName,
		Type:      dataType,
		Url:       key,
	}
	if err := dao.Image.Create(c, mus.Db, &faImage); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	base.JSON(c, code.MsgOk, gin.H{
		"bmind": RespImage{
			Name: fileName,
			Path: url,
		},
		"origin": RespImage{
			Name: fileName,
			Path: url,
		},
		"small": RespImage{
			Name: fileName,
			Path: url,
		},
	})
}

func Add(c *gin.Context) {
	// base64.StdEncoding.DecodeString(datasource)
	reqModel := ReqAdd{}
	if err := c.Bind(&reqModel); err != nil {
		mus.Logger.Error(err.Error(), zap.Int("code", 1))
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	if reqModel.Space == "" {
		reqModel.Space = "default"
	}

	ext := ".jpg"
	dataType := reqModel.Image[strings.IndexByte(reqModel.Image, ':')+1 : strings.IndexByte(reqModel.Image, ';')]
	b64data := reqModel.Image[strings.IndexByte(reqModel.Image, ',')+1:]

	b64dataDecode, _ := base64.StdEncoding.DecodeString(b64data) // 成图片文件并把文件写入到buffer

	fileName := imagex.GenerateUniqueMd5()
	rootPath, month := imagex.GeneratePath(reqModel.Space)

	filePath := filepath.Join(rootPath, fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err := ioutil.WriteFile(filePath, b64dataDecode, 0666) // buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	showPath := conf.Conf.Image.Domain + "/" + reqModel.Space + "/" + month + "/" + fileName + ext

	var faImage = mysql.Image{
		CreatedBy: mdw.DefaultContextUser(c).Id,
		UpdatedBy: mdw.DefaultContextUser(c).Id,
		Name:      fileName,
		Type:      dataType,
		Url:       showPath + "/200_200",
	}
	if err := dao.Image.Create(c, mus.Db, &faImage); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	base.JSON(c, code.MsgOk, gin.H{
		"bmind": RespImage{
			Name: fileName,
			Path: showPath + "/120_120",
		},
		"origin": RespImage{
			Name: fileName,
			Path: showPath + "/200_200",
		},
		"small": RespImage{
			Name: fileName,
			Path: showPath + "/60_60",
		},
	})
}

func List(c *gin.Context) {
	req := ReqImageList{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	total, list := dao.Image.ListPage(c, mysql.Conds{}, &req.ReqPage)
	base.JSONList(c, list, req.ReqPage.Current, req.ReqPage.PageSize, total)
}
