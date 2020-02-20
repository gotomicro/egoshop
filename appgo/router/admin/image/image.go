package image

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/dao"
	"github.com/i2eco/egoshop/appgo/model/constx"
	"github.com/i2eco/egoshop/appgo/model/mysql"
	"github.com/i2eco/egoshop/appgo/pkg/base"
	"github.com/i2eco/egoshop/appgo/pkg/code"
	"github.com/i2eco/egoshop/appgo/pkg/conf"
	"github.com/i2eco/egoshop/appgo/pkg/imagex"
	"github.com/i2eco/egoshop/appgo/pkg/mus"
	"github.com/i2eco/egoshop/appgo/router/mdw"
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

	srcPath := filepath.Join(rootPath, fileName+ext)

	path := filepath.Dir(srcPath)

	os.MkdirAll(path, os.ModePerm)

	err := ioutil.WriteFile(srcPath, b64dataDecode, 0666) // buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	prefix, _ := constx.OssMap[reqModel.OssType]
	if prefix == "" {
		base.JSONErr(c, 10003, errors.New("not exist type"))
		return
	}

	dstPath := mus.Oss.GenerateKey(prefix)
	err = mus.Oss.PutObjectFromFile(dstPath, srcPath)
	if err != nil {
		base.JSONErr(c, 10004, err)
		return
	}

	url, err := mus.Oss.SignURL(dstPath, "GET", 120)
	if err != nil {
		base.JSONErr(c, 10005, err)
		return
	}

	var faImage = mysql.Image{
		CreatedBy: mdw.DefaultContextUser(c).Id,
		UpdatedBy: mdw.DefaultContextUser(c).Id,
		Name:      fileName,
		Type:      dataType,
		Url:       dstPath,
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
