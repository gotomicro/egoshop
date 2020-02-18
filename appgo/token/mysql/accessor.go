package mysql

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goecology/muses/pkg/logger"
	"go.uber.org/zap"

	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/token"
)

type tokenAccessor struct {
	token.BaseTokenAccessor
	logger *logger.Client
}

func InitMysqlTokenAccessor(logger *logger.Client) token.TokenAccessor {
	return &tokenAccessor{
		BaseTokenAccessor: token.BaseTokenAccessor{},
		logger:            logger,
	}
}

func (g *tokenAccessor) CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp token.AccessTokenTicket, err error) {
	AccessTokenData := &mysql.AccessToken{
		Jti:        0,
		Sub:        uid,
		IaTime:     startTime,
		ExpTime:    startTime + token.AccessTokenExpireInterval,
		Ip:         "",
		CreateTime: time.Now().Unix(),
		IsLogout:   0,
		IsInvalid:  0,
		LogoutTime: 0,
	}

	err = dao.AccessToken.Create(c, AccessTokenData)
	if err != nil {
		return
	}
	tokenString, err := g.EncodeAccessToken(AccessTokenData.Jti, uid, startTime)
	if err != nil {
		return
	}
	resp.AccessToken = tokenString
	resp.ExpiresIn = token.AccessTokenExpireInterval
	return
}

func (g *tokenAccessor) CheckAccessToken(c *gin.Context, tokenStr string) bool {
	sc, err := g.DecodeAccessToken(tokenStr)
	if err != nil {
		g.logger.Error("access_token CheckAccessToken error1", zap.String("err", err.Error()))
		return false
	}

	conds := make(mysql.Conds, 5)
	conds["jti"] = sc["jti"]
	conds["sub"] = sc["sub"]
	conds["exp_time"] = mysql.Cond{
		Op:  ">=",
		Val: sc["exp"],
	}
	conds["is_invalid"] = 0
	conds["is_logout"] = 0
	if _, err := dao.AccessToken.InfoX(c, conds); err != nil {
		g.logger.Error("access_token CheckAccessToken error2", zap.String("err", err.Error()))
		return false
	}
	return true
}

func (g *tokenAccessor) RefreshAccessToken(c *gin.Context, tokenStr string, startTime int64) (resp token.AccessTokenTicket, err error) {
	sc, err := g.DecodeAccessToken(tokenStr)
	if err != nil {
		g.logger.Error("access_token CheckAccessToken error1", zap.String("err", err.Error()))
		return
	}

	refreshToken, err := g.EncodeAccessToken(sc["jti"].(int), sc["uid"].(int), startTime)

	if err != nil {
		return
	}

	err = dao.AccessToken.UpdateX(c, mysql.Conds{"jti": sc["jti"].(int)}, map[string]interface{}{
		"exp_time": startTime + token.AccessTokenExpireInterval,
	})
	if err != nil {
		return
	}
	resp.AccessToken = refreshToken
	resp.ExpiresIn = token.AccessTokenExpireInterval
	return
}
