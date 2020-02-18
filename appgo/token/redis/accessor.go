package redis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/goecology/muses/pkg/cache/mixcache/standard"
	"github.com/goecology/muses/pkg/logger"

	"github.com/goecology/egoshop/appgo/token"
)

const tokenKeyPattern = "/egoshop/token/%d"

type tokenAccessor struct {
	token.BaseTokenAccessor
	logger *logger.Client
	cache  standard.MixCache
}

func InitTokenAccessor(logger *logger.Client, cache standard.MixCache) token.TokenAccessor {
	return &tokenAccessor{
		BaseTokenAccessor: token.BaseTokenAccessor{},
		logger:            logger,
		cache:             cache,
	}
}

func (accessor *tokenAccessor) CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp token.AccessTokenTicket, err error) {

	// using the uid as the jwtId
	tokenString, err := accessor.EncodeAccessToken(uid, uid, startTime)
	if err != nil {
		return
	}

	_, err = accessor.cache.Set(fmt.Sprintf(tokenKeyPattern, uid), tokenString, token.AccessTokenExpireInterval)
	if err != nil {
		return
	}
	resp.AccessToken = tokenString
	resp.ExpiresIn = token.AccessTokenExpireInterval
	return
}

func (accessor *tokenAccessor) CheckAccessToken(c *gin.Context, tokenStr string) bool {
	sc, err := accessor.DecodeAccessToken(tokenStr)
	if err != nil {
		return false
	}
	uid := sc["jti"]
	_, err = accessor.cache.Get(fmt.Sprintf(tokenKeyPattern, uid))
	return err == nil
}

func (accessor *tokenAccessor) RefreshAccessToken(c *gin.Context, tokenStr string, startTime int64) (resp token.AccessTokenTicket, err error) {
	sc, err := accessor.DecodeAccessToken(tokenStr)
	if err != nil {
		return
	}
	uid := sc["jti"].(int)
	return accessor.CreateAccessToken(c, uid, startTime)
}
