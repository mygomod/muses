package redis

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/goecology/muses/pkg/cache/redis"
	"github.com/goecology/muses/pkg/logger"
	standard2 "github.com/goecology/muses/pkg/token/standard"
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
type redisTokenAccessor struct {
	standard2.JwtTokenAccessor
	logger *logger.Client
	client *redis.Client
}

func InitRedisTokenAccessor(logger *logger.Client, client *redis.Client) standard2.TokenAccessor {
	return &redisTokenAccessor{
		JwtTokenAccessor: standard2.JwtTokenAccessor{},
		logger:           logger,
		client:           client,
	}
}

func (accessor *redisTokenAccessor) CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp standard2.AccessTokenTicket, err error) {

	// using the uid as the jwtId
	tokenString, err := accessor.EncodeAccessToken(uid, uid, startTime)
	if err != nil {
		return
	}

	_, err = accessor.client.Set(fmt.Sprintf(tokenKeyPattern, uid), tokenString, standard2.AccessTokenExpireInterval)
	if err != nil {
		return
	}
	resp.AccessToken = tokenString
	resp.ExpiresIn = standard2.AccessTokenExpireInterval
	return
}

func (accessor *redisTokenAccessor) CheckAccessToken(c *gin.Context, tokenStr string) bool {
	sc, err := accessor.DecodeAccessToken(tokenStr)
	if err != nil {
		return false
	}
	uid := sc["jti"]
	_, err = accessor.client.Get(fmt.Sprintf(tokenKeyPattern, uid))
	return err == nil
}

func (accessor *redisTokenAccessor) RefreshAccessToken(c *gin.Context, tokenStr string, startTime int64) (resp standard2.AccessTokenTicket, err error) {
	sc, err := accessor.DecodeAccessToken(tokenStr)
	if err != nil {
		return
	}
	uid := sc["jti"].(int)
	return accessor.CreateAccessToken(c, uid, startTime)
}
