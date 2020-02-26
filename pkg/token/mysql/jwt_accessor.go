package mysql

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/i2eco/muses/pkg/logger"
	"github.com/i2eco/muses/pkg/token/standard"
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
type mysqlTokenAccessor struct {
	standard.JwtTokenAccessor
	logger *logger.Client
	db     *gorm.DB
}

func InitTokenAccessor(logger *logger.Client, db *gorm.DB) standard.TokenAccessor {
	return &mysqlTokenAccessor{
		JwtTokenAccessor: standard.JwtTokenAccessor{},
		logger:           logger,
		db:               db,
	}
}

func (accessor *mysqlTokenAccessor) CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp standard.AccessTokenTicket, err error) {
	AccessTokenData := &AccessToken{
		Jti:        0,
		Sub:        uid,
		IaTime:     startTime,
		ExpTime:    startTime + standard.AccessTokenExpireInterval,
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
	resp.ExpiresIn = standard.AccessTokenExpireInterval
	return
}

func (accessor *mysqlTokenAccessor) CheckAccessToken(c *gin.Context, tokenStr string) bool {
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

func (accessor *mysqlTokenAccessor) RefreshAccessToken(c *gin.Context, tokenStr string, startTime int64) (resp standard.AccessTokenTicket, err error) {
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
		"exp_time": startTime + standard.AccessTokenExpireInterval,
	}).Error

	if err != nil {
		return
	}

	resp.AccessToken = refreshToken
	resp.ExpiresIn = standard.AccessTokenExpireInterval
	return
}
