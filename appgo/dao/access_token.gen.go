package dao

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/muses/pkg/logger"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type accessToken struct {
	logger *logger.Client
	db     *gorm.DB
}

func InitAccessToken(logger *logger.Client, db *gorm.DB) *accessToken {
	return &accessToken{
		logger: logger,
		db:     db,
	}
}

// Create 新增一条记
func (g *accessToken) Create(c *gin.Context, data *mysql.AccessToken) (err error) {
	data.CreateTime = time.Now().Unix()
	if err = g.db.Create(data).Error; err != nil {
		g.logger.Error("create accessToken create error", zap.Error(err))
		return
	}
	return nil
}

// UpdateX Update的扩展方法，根据Cond更新一条或多条记录
func (g *accessToken) UpdateX(c *gin.Context, conds mysql.Conds, ups mysql.Ups) (err error) {

	sql, binds := mysql.BuildQuery(conds)
	if err = g.db.Table("access_token").Where(sql, binds...).Updates(ups).Error; err != nil {
		g.logger.Error("access_token update error", zap.Error(err))
		return
	}
	return
}

// DeleteX Delete的扩展方法，根据Cond删除一条或多条记录。如果有delete_time则软删除，否则硬删除。
func (g *accessToken) DeleteX(c *gin.Context, conds mysql.Conds) (err error) {
	sql, binds := mysql.BuildQuery(conds)

	if err = g.db.Table("access_token").Where(sql, binds...).Delete(&mysql.AccessToken{}).Error; err != nil {
		g.logger.Error("access_token delete error", zap.Error(err))
		return
	}

	return
}

// InfoX Info的扩展方法，根据Cond查询单条记录
func (g *accessToken) InfoX(c *gin.Context, conds mysql.Conds) (resp mysql.AccessToken, err error) {
	sql, binds := mysql.BuildQuery(conds)

	if err = g.db.Table("access_token").Where(sql, binds...).First(&resp).Error; err != nil {
		g.logger.Error("access_token info error", zap.Error(err))
		return
	}
	return
}

// List 查询list，extra[0]为sorts
func (g *accessToken) List(c *gin.Context, conds mysql.Conds, extra ...string) (resp []mysql.AccessToken, err error) {
	sql, binds := mysql.BuildQuery(conds)

	sorts := ""
	if len(extra) >= 1 {
		sorts = extra[0]
	}
	if err = g.db.Table("access_token").Where(sql, binds...).Order(sorts).Find(&resp).Error; err != nil {
		g.logger.Error("access_token info error", zap.Error(err))
		return
	}
	return
}

// ListPage 根据分页条件查询list
func (g *accessToken) ListPage(c *gin.Context, conds mysql.Conds, reqList *trans.ReqPage) (total int, respList []mysql.AccessToken) {
	if reqList.PageSize == 0 {
		reqList.PageSize = 10
	}
	if reqList.Current == 0 {
		reqList.Current = 1
	}
	sql, binds := mysql.BuildQuery(conds)

	db := g.db.Table("access_token").Where(sql, binds...)
	respList = make([]mysql.AccessToken, 0)
	db.Count(&total)
	db.Order(reqList.Sort).Offset((reqList.Current - 1) * reqList.PageSize).Limit(reqList.PageSize).Find(&respList)
	return
}
