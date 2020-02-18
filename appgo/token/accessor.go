package token

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const AccessTokenIss = "github.com/goecology/egoshop"
const AccessTokenExpireInterval = 7 * 24 * 60 * 60
const AccessTokenKey = "ecologysK#xo"

type AccessTokenTicket struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

type TokenAccessor interface {
	CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp AccessTokenTicket, err error)
	CheckAccessToken(c *gin.Context, tokenStr string) bool
	RefreshAccessToken(c *gin.Context, tokenStr string, startTime int64) (resp AccessTokenTicket, err error)
	EncodeAccessToken(jwtId int, uid int, startTime int64) (tokenStr string, err error)
	DecodeAccessToken(tokenStr string) (resp jwt.MapClaims, err error)
}

type BaseTokenAccessor struct {

}

func (g *BaseTokenAccessor) EncodeAccessToken(jwtId int, uid int, startTime int64) (tokenStr string, err error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["jti"] = jwtId
	claims["iss"] = AccessTokenIss
	claims["sub"] = uid
	claims["iat"] = startTime
	claims["exp"] = startTime + AccessTokenExpireInterval
	jwtToken.Claims = claims

	tokenStr, err = jwtToken.SignedString([]byte(AccessTokenKey))
	if err != nil {
		return
	}
	return
}

func (g *BaseTokenAccessor) DecodeAccessToken(tokenStr string) (resp jwt.MapClaims, err error) {
	tokenParse, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
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

var accessor TokenAccessor

// Register 和 GetAccessor方法是暂时为了解决使用何种实现的方案。
// 更加优雅的方案，在系统初始化的时候，设置好这个Accessor
func Register(acc TokenAccessor)  {
	accessor = acc
}

func GetAccessor() TokenAccessor{
	return accessor
}