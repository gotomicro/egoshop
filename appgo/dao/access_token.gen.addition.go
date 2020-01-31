package dao

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"go.uber.org/zap"
)

const AccessTokenIss = "github.com/goecology/egoshop"
const AccessTokenExpireInterval = 7 * 24 * 60 * 60
const AccessTokenKey = "ecologysK#xo"

type AccessTokenTicket struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

func (g *accessToken) CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp AccessTokenTicket, err error) {
	AccessTokenData := &mysql.AccessToken{
		Jti:        0,
		Sub:        uid,
		IaTime:     startTime,
		ExpTime:    startTime + AccessTokenExpireInterval,
		Ip:         "",
		CreateTime: time.Now().Unix(),
		IsLogout:   0,
		IsInvalid:  0,
		LogoutTime: 0,
	}

	err = g.Create(c, g.db, AccessTokenData)
	if err != nil {
		return
	}
	tokenString, err := encodeAccessToken(AccessTokenData.Jti, uid, startTime)
	if err != nil {
		return
	}
	resp.AccessToken = tokenString
	resp.ExpiresIn = AccessTokenExpireInterval
	return
}

func (g *accessToken) CheckAccessToken(token string) bool {
	sc, err := g.DecodeAccessToken(token)
	if err != nil {
		g.logger.Error("access_token CheckAccessToken error1", zap.String("err", err.Error()))
		return false
	}

	var resp mysql.AccessToken
	if err = g.db.Table("access_token").Where("`jti`=? AND `sub`=? AND `exp_time`>=? AND `is_invalid`=? AND `is_logout`=?", sc["jti"], sc["sub"], sc["exp"], 0, 0).Find(&resp).Error; err != nil {
		g.logger.Error("access_token CheckAccessToken error2", zap.String("err", err.Error()))
		return false
	}
	return true
}

func (g *accessToken) RefreshAccessToken(c *gin.Context, token string, startTime int64) (resp AccessTokenTicket, err error) {
	sc, err := g.DecodeAccessToken(token)
	if err != nil {
		g.logger.Error("access_token CheckAccessToken error1", zap.String("err", err.Error()))
		return
	}

	refreshToken, err := encodeAccessToken(sc["jti"].(int), sc["uid"].(int), startTime)

	if err != nil {
		return
	}

	err = g.UpdateX(c, g.db, mysql.Conds{"jti": sc["jti"].(int)}, map[string]interface{}{
		"exp_time": startTime + AccessTokenExpireInterval,
	})
	if err != nil {
		return
	}
	resp.AccessToken = refreshToken
	resp.ExpiresIn = AccessTokenExpireInterval
	return
}

func encodeAccessToken(jwtId int, uid int, startTime int64) (tokenStr string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["jti"] = jwtId
	claims["iss"] = AccessTokenIss
	claims["sub"] = uid
	claims["iat"] = startTime
	claims["exp"] = startTime + AccessTokenExpireInterval
	token.Claims = claims

	tokenStr, err = token.SignedString([]byte(AccessTokenKey))
	if err != nil {
		return
	}
	return
}

func (g *accessToken) DecodeAccessToken(token string) (resp jwt.MapClaims, err error) {
	tokenParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(AccessTokenKey), nil
	})
	if err != nil {
		return
	}
	var flag bool
	resp, flag = tokenParse.Claims.(jwt.MapClaims)
	if !flag {
		err = errors.New("assert error")
		return
	}
	return
}
