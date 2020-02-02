package admincom

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/goecology/egoshop/appgo/pkg/imagex"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/mdw"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/goecology/egoshop/appgo/pkg/util"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
)

func List(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	conds := mysql.Conds{}

	// 搜索名称
	if req.Name != "" {
		conds["name"] = mysql.Cond{
			"like",
			req.Name,
		}
	}

	// 搜索繁育人
	if req.Author != "" {
		conds["author"] = mysql.Cond{
			"like",
			req.Author,
		}
	}

	// 搜索地址
	if req.Author != "" {
		conds["address"] = mysql.Cond{
			"like",
			req.Address,
		}
	}

	req.Sort = "id desc"
	total, list := dao.Com.ListPage(c, conds, &req.ReqPage)

	//处理封面图片
	for idx, comInfo := range list {
		resp, _ := dao.ComSku.List(c, mysql.Conds{
			"com_id": comInfo.Id,
		})
		comInfo.SkuList = resp
		comInfo.Cover = imagex.ShowImg(comInfo.Cover, "x1")
		list[idx] = comInfo
	}

	base.JSONList(c, list, req.ReqPage.Current, req.ReqPage.PageSize, total)
}

func One(c *gin.Context) {
	reqId := cast.ToInt(c.Param("id"))
	if reqId == 0 {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	comInfo, _ := dao.Com.Info(c, reqId)
	resp, _ := dao.ComSku.List(c, mysql.Conds{
		"com_id": comInfo.Id,
	})
	comInfo.SkuList = resp
	comInfo.Gallery = imagex.ShowImgArr(comInfo.Gallery,"")
	base.JSON(c, code.MsgOk, comInfo)
}

//获取或更新文档内容.
func Content(c *gin.Context) {
	reqId := cast.ToInt(c.Param("id"))
	if reqId == 0 {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	comInfo, err := dao.Com.Info(c, reqId)
	if err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	comStore, _ := dao.ComStore.InfoX(c, mysql.Conds{"com_id": comInfo.Id})
	comInfo.PreMarkdown = comStore.PreMarkdown
	comInfo.Markdown = comStore.Markdown
	comInfo.Html = comStore.Html
	comInfo.WechatHtml = comStore.WechatHtml

	base.JSON(c, code.MsgOk, comInfo)
}

func Create(c *gin.Context) {
	req := ReqCreateOrUpdate{}

	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	if len(req.Gallery) < 1 {
		base.JSON(c, code.MsgErr, "gallery length is error")
		return
	}
	createParam, err := filterCreateParam(req)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx := mus.Db.Begin()
	uid := mdw.AdminUid(c)
	createCom := &mysql.Com{
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		CreatedBy:   uid,
		UpdatedBy:   uid,
		Title:       createParam.Title,
		SubTitle:    createParam.SubTitle,
		Cid:         req.Cid,
		Cover:       createParam.Cover,
		Gallery:     createParam.Gallery,
		Stock:       createParam.Stock,
		SaleNum:     createParam.BaseSaleNum,
		SaleTime:    createParam.SaleTime,
		PayType:     createParam.PayType,
		FreightFee:  createParam.Freight,
		FreightId:   createParam.FreightTemplateId,
		BaseSaleNum: createParam.BaseSaleNum,
		IsOnSale:    createParam.IsOnSale,
		ImageSpecId: createParam.ImageSpecId,
		Price:       createParam.Price,
		SkuList:     createParam.SkuList,
		SpecList:    createParam.SpecList,
	}

	err = dao.Com.Create(c, tx, createCom)

	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	for i := range createCom.SkuList {
		_, specSign, specValueSign, title := skuSign(createCom.SkuList[i].Spec)
		createCom.SkuList[i].ComId = createCom.Id
		createCom.SkuList[i].Title = title
		createCom.SkuList[i].SpecSign = specSign
		createCom.SkuList[i].SpecValueSign = specValueSign
		err = dao.ComSku.Create(c, tx, &createCom.SkuList[i])
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}

	for _, img := range createCom.Gallery {
		err = dao.ComImage.Create(c, tx, &mysql.ComImage{
			Img:   img,
			ComId: createCom.Id,
		})
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}

	for _, cid := range createParam.Cids {
		err = dao.ComRelateCate.Create(c, tx, &mysql.ComRelateCate{
			ComId: createCom.Id,
			Cid:   cid,
		})
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}
	tx.Commit()

	reqPage := &trans.ReqPage{}
	total, list := dao.Com.ListPage(c, mysql.Conds{}, reqPage)
	base.JSONList(c, list, reqPage.Current, reqPage.PageSize, total)
}

func Update(c *gin.Context) {
	req := ReqCreateOrUpdate{}

	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	if len(req.Gallery) < 1 {
		base.JSON(c, code.MsgErr, "gallery length is error")
		return
	}
	createParam, err := filterCreateParam(req)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx := mus.Db.Begin()
	uid := mdw.AdminUid(c)

	err = dao.Com.Update(c, tx, createParam.Id, mysql.Ups{
		"title":             createParam.Title,
		"gallery":           createParam.Gallery,
		"cids":              createParam.Cids,
		"base_sale_num":     createParam.BaseSaleNum,
		"is_on_sale":        createParam.IsOnSale,
		"image_spec_id":     createParam.ImageSpecId,
		"image_spec_images": createParam.ImageSpecImages,
		"sku_list":          createParam.SkuList,
		"price":             createParam.Price,
		"stock":             createParam.Stock,
		"sale_time":         createParam.SaleTime,
		"spec_list":         createParam.SpecList,
		"cover":             createParam.Cover,
		"pay_type":          createParam.PayType,
		"freight_fee":       createParam.Freight,
		"freight_id":        createParam.FreightTemplateId,
		"updated_by":        uid,
	})
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	for i := range createParam.SkuList {
		var specSign []int
		var specValueSign []int
		for _, v := range createParam.SkuList[i].Spec {
			specSign = append(specSign, v.ID)
			specValueSign = append(specValueSign, v.ValueID)
		}
		specSignJ := util.JsonMarshal(specSign)
		specValueSignJ := util.JsonMarshal(specValueSign)
		// update.Skus[i].GoodsId = id
		createParam.SkuList[i].Title = createParam.Title
		createParam.SkuList[i].SpecSign = string(specSignJ)
		createParam.SkuList[i].SpecValueSign = string(specValueSignJ)
		err = dao.ComSku.Create(c, tx, &createParam.SkuList[i])
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}
	for _, sku := range createParam.SkuList {
		err = dao.ComSku.Update(c, tx, sku.Id, mysql.Ups{
			"title": createParam.Title,
		})
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}

	err = dao.ComImage.Delete(c, tx, createParam.Id)
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	for _, img := range createParam.Gallery {
		// TODO 需要Create吗？
		err = dao.ComImage.Create(c, tx, &mysql.ComImage{
			Img:   img,
			ComId: createParam.Id,
		})
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}

	err = dao.ComRelateCate.Delete(c, tx, createParam.Id)
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	for _, cid := range createParam.Cids {
		// TODO 需要Create吗？
		err = dao.ComRelateCate.Create(c, tx, &mysql.ComRelateCate{
			ComId: createParam.Id,
			Cid:   cid,
		})
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
	}
	tx.Commit()

	reqPage := &trans.ReqPage{}
	total, list := dao.Com.ListPage(c, mysql.Conds{}, &trans.ReqPage{})
	base.JSONList(c, list, reqPage.Current, reqPage.PageSize, total)

}

func Remove(c *gin.Context) {
	req := ReqRemove{}
	reqPage := &trans.ReqPage{}

	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	err := dao.Com.Update(c, mus.Db, req.Id, mysql.Ups{
		"deleted_at": time.Now(),
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	total, list := dao.Com.ListPage(c, mysql.Conds{}, &trans.ReqPage{})
	base.JSONList(c, list, reqPage.Current, reqPage.PageSize, total)
}

func filterCreateParam(req ReqCreateOrUpdate) (resp createMysql, err error) {
	resp.Id = req.Id
	resp.Title = req.Title
	resp.Gallery = req.Gallery
	if len(req.Gallery) > 0 {
		resp.Cover = req.Gallery[0]
	}

	// todo 需要验证分类id是否存在
	for _, value := range req.Cids {
		resp.Cids = append(resp.Cids, cast.ToInt(value))
	}

	//两个函数相比，不一样的地方在于，Parse()函数解析的时候，会默为UTC时间，获取的Time对象转换为Unix()对象后，会比当前时间多8小时。
	//
	//tm, err := time.Parse("2006-01-02T15:04:05Z", s) //转换后的时间，如果再转换为unix时间，需要-8小时
	//
	//
	//如果解析来源是GMT的时间的话，最好使用ParseInLocation()，并指定"*Location"为“time.Local”，比如：
	//
	//tm, err = time.ParseInLocation("2006-01-02T15:04:05Z", s, time.Local)  //转换后的时间如果再转换为unix时间，不需要处理。

	saleTime, err := time.Parse("2006-01-02T15:04:05Z", req.SaleTime)
	if err != nil {
		return
	}
	resp.IsOnSale = 0
	resp.BaseSaleNum = 0
	resp.FreightTemplateId = req.FreightId
	resp.Freight = req.FreightFee
	resp.SaleTime = saleTime
	resp.ImageSpecId = 0
	resp.SkuList = make([]mysql.ComSku, 0)
	resp.SpecMap = make(map[int]mysql.ComSpecOneInfo)
	resp.SpecList = make(mysql.ComSpecListJson, 0)

	if len(req.SkuList) > 0 {
		staticValueExistIds := make([]int, 0)
		for key, sku := range req.SkuList {
			// 初始化每个sku的图片（下面需要处理：如果选择了图片规格，自动设置为图片规格封面）
			sku.Cover = resp.Cover
			resp.SkuList = append(resp.SkuList, sku)

			// 价格
			if key == 0 {
				resp.Price = sku.Price
			} else if resp.Price > sku.Price {
				resp.Price = sku.Price
			}
			// 库存
			resp.Stock = sku.Stock

			// 规格层级集合json 不要重复，如色彩下有:xx色 xx 色
			for _, spec := range sku.Spec {
				if !funk.Contains(staticValueExistIds, spec.ValueID) {
					// 存着防止 sku 规格循环被重复记录
					staticValueExistIds = append(staticValueExistIds, spec.ValueID)
					if _, ok := resp.SpecMap[spec.ID]; !ok {
						resp.SpecMap[spec.ID] = mysql.ComSpecOneInfo{
							Id:        0,
							Name:      "",
							ValueList: make([]mysql.CreateSpecValue, 0),
						}
					}
					specTmp := resp.SpecMap[spec.ID]
					specTmp.Id = spec.ID
					specTmp.Name = spec.Name
					specTmp.ValueList = append(specTmp.ValueList, mysql.CreateSpecValue{
						Id:   spec.ValueID,
						Name: spec.ValueName,
					})

					// 规格图片,防止0或空进来 当默认没规格时是会为空
					if resp.ImageSpecId > 0 && resp.ImageSpecId == spec.ID {
						resp.ImageSpecImages = append(resp.ImageSpecImages, spec.ValueImg)
					}
					resp.SpecMap[spec.ID] = specTmp
				}
			}
		}
		for _, value := range resp.SpecMap {
			resp.SpecList = append(resp.SpecList, value)
		}
		if len(resp.SkuList) == 0 || len(resp.SpecList) == 0 {
			err = errors.New("skulist or speclist error")
			return
		}
	}

	resp.ImageSpecImages = resp.ImageSpecImages
	resp.Body = resp.Body
	resp.SpecList = resp.SpecList
	resp.SkuList = resp.SkuList
	resp.Cids = resp.Cids
	resp.Gallery = req.Gallery
	resp.SubTitle = req.SubTitle
	return
}

func skuSign(specs mysql.ComSkuSpecJson) (spec string, specSign string, specValueSign string, title string) {
	ids := make([]int, 0)
	valueIds := make([]int, 0)
	valueName := make([]string, 0)
	for _, value := range specs {
		ids = append(ids, value.ID)
		valueIds = append(valueIds, value.ValueID)
		valueName = append(valueName, value.ValueName)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	sort.Slice(valueIds, func(i, j int) bool {
		return valueIds[i] < valueIds[j]
	})
	spec = util.JsonMarshal(specs)
	specSign = util.JsonMarshal(ids)
	specValueSign = util.JsonMarshal(valueIds)
	title = strings.Join(valueName, " ")
	return
}

type createMysql struct {
	Id                int
	SubTitle          string
	Title             string
	ImageSpecImages   mysql.ComImageSpecImagesJson
	ImageSpecId       int
	Body              mysql.ComBodyJson
	Stock             int
	Freight           float64
	SaleTime          time.Time
	SpecList          mysql.ComSpecListJson
	SpecMap           map[int]mysql.ComSpecOneInfo
	BaseSaleNum       int
	Price             float64
	Cover             string
	SkuList           []mysql.ComSku
	Cids              mysql.ComCategoryIdsJson
	Gallery           mysql.ComGalleryJson
	FreightTemplateId int
	IsOnSale          int
	PayType           int
}
