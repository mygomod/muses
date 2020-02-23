package mysql

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/i2eco/muses/pkg/logger"
	"github.com/i2eco/muses/pkg/token"
	"go.uber.org/zap"
)

// 如果你希望使用这个实现来作为token的实现，那么需要在配置文件里面设置：
// [muses.logger.system]
//    ...logger的配置
// [muses.mysql.default]
//    ...mysql的配置
// [muses.token.jwt.mysql]
//    logger = "system"
//    client = "default"
// 而后将Register()方法注册进去muses.Container(...)中
type tokenAccessor struct {
	token.JwtTokenAccessor
	logger *logger.Client
	db     *gorm.DB
}

func initTokenAccessor(logger *logger.Client, db *gorm.DB) token.TokenAccessor {
	return &tokenAccessor{
		JwtTokenAccessor: token.JwtTokenAccessor{},
		logger:           logger,
		db:               db,
	}
}

func (accessor *tokenAccessor) CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp token.AccessTokenTicket, err error) {
	AccessTokenData := &AccessToken{
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

	tokenString, err := accessor.EncodeAccessToken(AccessTokenData.Jti, uid, startTime)
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
		accessor.logger.Error("access_token CheckAccessToken error1", zap.String("err", err.Error()))
		return false
	}
	var resp AccessToken
	if err = accessor.db.Table(TableName).Where("`jti`=? AND `sub`=? AND `exp_time`>=? AND `is_invalid`=? AND `is_logout`=?", sc["jti"], sc["sub"], sc["exp"], 0, 0).Find(&resp).Error; err != nil {
		accessor.logger.Error("access_token CheckAccessToken error2", zap.String("err", err.Error()))
		return false
	}
	return true
}

func (accessor *tokenAccessor) RefreshAccessToken(c *gin.Context, tokenStr string, startTime int64) (resp token.AccessTokenTicket, err error) {
	sc, err := accessor.DecodeAccessToken(tokenStr)
	if err != nil {
		accessor.logger.Error("access_token CheckAccessToken error1", zap.String("err", err.Error()))
		return
	}

	jti := sc["jti"].(int)
	refreshToken, err := accessor.EncodeAccessToken(jti, sc["uid"].(int), startTime)

	if err != nil {
		return
	}

	err = accessor.db.Table(TableName).Where("`jti`=?", jti).Updates(map[string]interface{}{
		"exp_time": startTime + token.AccessTokenExpireInterval,
	}).Error

	if err != nil {
		return
	}

	resp.AccessToken = refreshToken
	resp.ExpiresIn = token.AccessTokenExpireInterval
	return
}
