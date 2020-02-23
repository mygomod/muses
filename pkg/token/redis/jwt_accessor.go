package redis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/i2eco/muses/pkg/cache/mixcache/standard"
	"github.com/i2eco/muses/pkg/logger"
	"github.com/i2eco/muses/pkg/token"
)

// TODO(配置中读取)
const tokenKeyPattern = "/egoshop/token/%d"

// 如果你希望使用这个实现来作为token的实现，那么需要在配置文件里面设置：
// [muses.logger.system]
//    ...logger的配置
// [muses.redis.default]
//    ...mysql的配置
// [muses.token.jwt.redis]
//    logger = "system"
//    client = "default"
// 而后将Register()方法注册进去muses.Container(...)中
type tokenAccessor struct {
	token.JwtTokenAccessor
	logger *logger.Client
	cache  standard.MixCache
}

func initTokenAccessor(logger *logger.Client, cache standard.MixCache) token.TokenAccessor {
	return &tokenAccessor{
		JwtTokenAccessor: token.JwtTokenAccessor{},
		logger:           logger,
		cache:            cache,
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


