package dao

import (
	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/model/mysql"
	"github.com/i2eco/egoshop/appgo/model/trans"
	"github.com/i2eco/muses/pkg/logger"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type cart struct {
	logger *logger.Client
	db     *gorm.DB
}

func InitCart(logger *logger.Client, db *gorm.DB) *cart {
	return &cart{
		logger: logger,
		db:     db,
	}
}

// Create 新增一条记
func (g *cart) Create(c *gin.Context, db *gorm.DB, data *mysql.Cart) (err error) {

	if err = db.Create(data).Error; err != nil {
		g.logger.Error("create cart create error", zap.Error(err))
		return
	}
	return nil
}

// Update 根据主键更新一条记录
func (g *cart) Update(c *gin.Context, db *gorm.DB, paramId int, ups mysql.Ups) (err error) {
	var sql = "`id`=?"
	var binds = []interface{}{paramId}

	if err = db.Table("cart").Where(sql, binds...).Updates(ups).Error; err != nil {
		g.logger.Error("cart update error", zap.Error(err))
		return
	}
	return
}

// UpdateX Update的扩展方法，根据Cond更新一条或多条记录
func (g *cart) UpdateX(c *gin.Context, db *gorm.DB, conds mysql.Conds, ups mysql.Ups) (err error) {

	sql, binds := mysql.BuildQuery(conds)
	if err = db.Table("cart").Where(sql, binds...).Updates(ups).Error; err != nil {
		g.logger.Error("cart update error", zap.Error(err))
		return
	}
	return
}

// Delete 根据主键删除一条记录。如果有delete_time则软删除，否则硬删除。
func (g *cart) Delete(c *gin.Context, db *gorm.DB, paramId int) (err error) {
	var sql = "`id`=?"
	var binds = []interface{}{paramId}

	if err = db.Table("cart").Where(sql, binds...).Delete(&mysql.Cart{}).Error; err != nil {
		g.logger.Error("cart delete error", zap.Error(err))
		return
	}

	return
}

// DeleteX Delete的扩展方法，根据Cond删除一条或多条记录。如果有delete_time则软删除，否则硬删除。
func (g *cart) DeleteX(c *gin.Context, db *gorm.DB, conds mysql.Conds) (err error) {
	sql, binds := mysql.BuildQuery(conds)

	if err = db.Table("cart").Where(sql, binds...).Delete(&mysql.Cart{}).Error; err != nil {
		g.logger.Error("cart delete error", zap.Error(err))
		return
	}

	return
}

// Info 根据PRI查询单条记录
func (g *cart) Info(c *gin.Context, paramId int) (resp mysql.Cart, err error) {
	var sql = "`id`=?"
	var binds = []interface{}{paramId}

	if err = g.db.Table("cart").Where(sql, binds...).First(&resp).Error; err != nil {
		g.logger.Error("cart info error", zap.Error(err))
		return
	}
	return
}

// InfoX Info的扩展方法，根据Cond查询单条记录
func (g *cart) InfoX(c *gin.Context, conds mysql.Conds) (resp mysql.Cart, err error) {
	sql, binds := mysql.BuildQuery(conds)

	if err = g.db.Table("cart").Where(sql, binds...).First(&resp).Error; err != nil {
		g.logger.Error("cart info error", zap.Error(err))
		return
	}
	return
}

// List 查询list，extra[0]为sorts
func (g *cart) List(c *gin.Context, conds mysql.Conds, extra ...string) (resp []mysql.Cart, err error) {
	sql, binds := mysql.BuildQuery(conds)

	sorts := ""
	if len(extra) >= 1 {
		sorts = extra[0]
	}
	if err = g.db.Table("cart").Where(sql, binds...).Order(sorts).Find(&resp).Error; err != nil {
		g.logger.Error("cart info error", zap.Error(err))
		return
	}
	return
}

// ListMap 查询map，map遍历的时候是无序的，所以指定sorts参数没有意义
func (g *cart) ListMap(c *gin.Context, conds mysql.Conds) (resp map[int]mysql.Cart, err error) {
	sql, binds := mysql.BuildQuery(conds)

	mysqlSlice := make([]mysql.Cart, 0)
	resp = make(map[int]mysql.Cart, 0)
	if err = g.db.Table("cart").Where(sql, binds...).Find(&mysqlSlice).Error; err != nil {
		g.logger.Error("cart info error", zap.Error(err))
		return
	}
	for _, value := range mysqlSlice {
		resp[value.Id] = value
	}
	return
}

// ListPage 根据分页条件查询list
func (g *cart) ListPage(c *gin.Context, conds mysql.Conds, reqList *trans.ReqPage) (total int, respList []mysql.Cart) {
	if reqList.PageSize == 0 {
		reqList.PageSize = 10
	}
	if reqList.Current == 0 {
		reqList.Current = 1
	}
	sql, binds := mysql.BuildQuery(conds)

	db := g.db.Table("cart").Where(sql, binds...)
	respList = make([]mysql.Cart, 0)
	db.Count(&total)
	db.Order(reqList.Sort).Offset((reqList.Current - 1) * reqList.PageSize).Limit(reqList.PageSize).Find(&respList)
	return
}
