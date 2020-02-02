package dao

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/util"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
)

const (
	UserInit    = 0
	UserActived = 1
	UserBanned  = 2
)

func (g *user) GetRandNickname() (n string, err error) {
	n = "fs_" + funk.RandomString(8)
	var req = mysql.User{}
	if err = g.db.Table("user").Where("`nickname`=? ", n).First(&req).Error; err != nil {
		g.logger.Error("biz update error", zap.String("err", err.Error()))
		return
	}
	return g.GetRandNickname()
}

func (g *user) AddBiz(nickname string, pwd string, ip string) (err error) {
	var pwdHash string
	pwdHash, err = util.Hash(pwd)
	if err != nil {
		g.logger.Debug("add user hash error", zap.String("err", err.Error()))
		return
	}
	user := mysql.User{
		Name:        nickname,
		Password:    pwdHash,
		Status:      UserActived,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastLoginIP: ip,
	}
	if err = g.db.Create(&user).Error; err != nil {
		g.logger.Debug("add user create error", zap.String("err", err.Error()))
		return
	}
	return nil
}

func (g *user) GetBizByOidPwd(openId int, pwd string) (resp *mysql.User, err error) {
	resp = &mysql.User{}
	if err = g.db.Where("open_id = ?", openId).First(resp).Error; err != nil {
		g.logger.Debug("GetBizByPwd ERROR", zap.String("err", err.Error()))
		return
	}
	err = util.Verify(resp.Password, pwd)
	if err != nil {
		g.logger.Debug("verify error1", zap.String("err", err.Error()))
		return
	}
	return
}

func (g *user) UpdatePwd(openId int, pwd string) (err error) {
	var pwdHash string
	pwdHash, err = util.Hash(pwd)
	if err != nil {
		g.logger.Debug("update user hash error", zap.String("err", err.Error()))
		return
	}
	if err = g.db.Table("biz").Where("open_id = ?", openId).Updates(gin.H{
		"password":   pwdHash,
		"updated_at": time.Now().Unix(),
	}).Error; err != nil {
		g.logger.Debug("update user create error", zap.String("err", err.Error()))
		return
	}
	return nil
}

// Login 用户登录.
func (g *user) Login(name string, password string) (*mysql.User, error) {
	member := &mysql.User{}
	var err error
	err = g.db.Where("name = ? AND status = ?", name, UserActived).Find(member).Error
	if err != nil {
		return member, err
	}
	err = util.Verify(member.Password, password)
	if err == nil {
		return member, nil
	}
	return member, ErrorMemberPasswordError
}

func (g *user) GetBizByPwd(nickname string, email string, pwd string, clientIp string) (resp *mysql.User, err error) {
	query := "1=1"
	data := make([]interface{}, 0)
	if nickname != "" {
		query += " and name = ? AND status = ?"
		data = append(data, nickname, UserActived)
	}
	if email != "" {
		query += " and email = ? AND status = ? "
		data = append(data, email, UserActived)
	}
	resp = &mysql.User{}
	if err = g.db.Where(query, data...).First(resp).Error; err != nil {
		g.logger.Debug("GetBizByPwd ERROR", zap.String("err", err.Error()))
		return
	}
	err = util.Verify(resp.Password, pwd)
	if err != nil {
		g.logger.Debug("verify error1", zap.String("err", err.Error()))
		return
	}
	if err = g.db.Table("user").Where("id = ?", resp.Id).Updates(gin.H{
		"updated_at":    time.Now(),
		"last_login_ip": clientIp,
	}).Error; err != nil {
		g.logger.Debug("update user create error", zap.String("err", err.Error()))
		return
	}
	return
}

//分页查找用户.
func (m *user) FindToPager(c *gin.Context, pageIndex, pageSize int) ([]*mysql.User, int) {
	members := make([]*mysql.User, 0)

	total, resp := m.ListPage(c, mysql.Conds{}, &trans.ReqPage{
		Current:  pageIndex,
		PageSize: pageSize,
		Sort:     "id desc",
	})

	for _, m := range resp {
		tmp := &m
		members = append(members, tmp)
	}
	return members, total
}

//根据账号查找用户.
func (g *user) FindByAccount(account string) (*mysql.User, error) {
	member := &mysql.User{}
	err := g.db.Where("account = ?", account).Find(member).Error
	if err == nil {
	}
	return member, err
}

func (g *user) Find(id int) (*mysql.User, error) {
	member := &mysql.User{}
	err := g.db.Where("id = ?", id).Find(member).Error
	if err != nil {
		return member, err
	}
	return member, nil
}

//根据指定字段查找用户.
func (g *user) FindByFieldFirst(field string, value interface{}) (*mysql.User, error) {
	member := &mysql.User{}
	var err error
	err = g.db.Where(field+" = ?", value).Order("id desc").Find(member).Error
	if err != nil {
		return member, err
	}
	return member, err
}

//获取用户名
func (g *user) GetUsernameByUid(id interface{}) string {
	var user mysql.User
	g.db.Where("id = ?", id).Find(&user)
	return user.Account
}

//获取昵称
func (g *user) GetNicknameByUid(id interface{}) string {
	var user mysql.User
	g.db.Where("id = ?", id).Find(&user)
	return user.Account
}

//根据用户id获取二维码
func (g *user) GetQrcodeByUid(uid interface{}) (qrcode map[string]string) {
	var member mysql.User
	g.db.Select("alipay,wxpay").Where("id = ?", uid).Find(&member)
	qrcode = make(map[string]string)
	qrcode["Alipay"] = member.Alipay
	qrcode["Wxpay"] = member.Wxpay
	return qrcode
}
