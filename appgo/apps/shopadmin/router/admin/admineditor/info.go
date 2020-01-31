package admineditor

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/conf"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/mdw"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/constx"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/goecology/egoshop/appgo/pkg/imagex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//获取或更新文档内容.
func ContentSave(c *gin.Context) {
	req := ReqContentSave{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	if req.Id <= 0 {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	var err error

	switch req.Type {
	case constx.TypeCom:
		var one mysql.ComStore
		mus.Db.Where("com_id=?", req.Id).Find(&one)
		if one.ComId > 0 {
			err = dao.ComStore.UpdateX(c, mus.Db, mysql.Conds{"com_id": req.Id}, mysql.Ups{
				"pre_markdown": req.PreMarkdown,
				"pre_html":     req.PreHtml,
				"updated_by":   mdw.AdminUid(c),
			})
		} else {
			err = mus.Db.Create(&mysql.ComStore{
				ComId:       req.Id,
				PreMarkdown: req.PreHtml,
				PreHtml:     req.PreMarkdown,
				Markdown:    "",
				Html:        "",
				WechatHtml:  "",
				CreatedBy:   mdw.AdminUid(c),
				UpdatedBy:   mdw.AdminUid(c),
			}).Error
		}
		if err != nil {
			base.JSON(c, code.MsgErr, "updated is error")
			return
		}
	default:
		base.JSON(c, code.MsgErr, "type is error")
		return
	}
	base.JSONOK(c)

}

//获取或更新文档内容.
func Release(c *gin.Context) {
	req := ReqRelease{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	if req.Id <= 0 {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	var err error
	switch req.Type {
	case constx.TypeCom:
		var one mysql.ComStore
		err = mus.Db.Where("com_id=?", req.Id).Find(&one).Error
		if err != nil {
			base.JSON(c, code.MsgErr, "data is empty")
			return
		}
		err = dao.ComStore.UpdateX(c, mus.Db, mysql.Conds{"com_id": req.Id}, mysql.Ups{
			"markdown":   one.PreMarkdown,
			"html":       one.PreHtml,
			"updated_by": mdw.AdminUid(c),
		})
		if err != nil {
			base.JSON(c, code.MsgErr, "updated is error")
			return
		}
	default:
		base.JSON(c, code.MsgErr, "type is error")
		return
	}
	base.JSONOK(c)
}

//上传附件或图片.
func Upload(c *gin.Context) {
	req := ReqUpload{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	if req.Id <= 0 {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	var err error

	switch req.Type {
	case constx.TypeCom:
		req.Space = "com"
	default:
		base.JSON(c, code.MsgErr, "type is error")
		return
	}

	ext := ".jpg"
	//dataType := req.Image[strings.IndexByte(req.Image, ':')+1 : strings.IndexByte(req.Image, ';')]
	b64data := req.Image[strings.IndexByte(req.Image, ',')+1:]

	b64dataDecode, _ := base64.StdEncoding.DecodeString(b64data) // 成图片文件并把文件写入到buffer

	fileName := imagex.GenerateUniqueMd5()
	rootPath, month := imagex.GeneratePath(req.Space)

	filePath := filepath.Join(rootPath, fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = ioutil.WriteFile(filePath, b64dataDecode, 0666) // buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	// showPath := conf.Conf.Image.Domain + "/" + req.Space + "/" + month + "/" + fileName + ext

	//
	//if !conf.IsAllowUploadFileExt(ext) {
	//	base.JSONErr(c, code.FeedUploadErr5, errors.New("ext is not allowed"))
	//	return
	//}

	attachment := mysql.Attachment{
		GoodsId:   req.Id,
		GoodsType: req.Type,
		FileName:  fileName + ext,
		FilePath:  "/" + req.Space + "/" + month + "/" + fileName + ext,
		FileSize:  0,
		HttpPath:  conf.Conf.Image.Domain + "/" + req.Space + "/" + month + "/" + fileName + ext,
		FileExt:   ext,
		CreatedBy: mdw.DefaultContextUser(c).Id,
	}

	if fileInfo, err := os.Stat(filePath); err == nil {
		attachment.FileSize = float64(fileInfo.Size())
	}

	err = dao.Attachment.Create(c, mus.Db, &attachment)

	if err != nil {
		os.Remove(filePath)
		base.JSONErr(c, code.FeedUploadErr7, err)
		return
	}
	result := map[string]interface{}{
		"url":    attachment.HttpPath,
		"alt":    attachment.FileName,
		"attach": attachment,
	}
	base.JSON(c, code.MsgOk, result)
}
