package token

import (
	"github.com/gin-gonic/gin"
)

const AccessTokenExpireInterval = 7 * 24 * 60 * 60

type AccessTokenTicket struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

// TokenAccessor目前来说，单纯是指基于JWT的Token。当我们需要扩展到其它的实现的时候
type TokenAccessor interface {
	CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp AccessTokenTicket, err error)
	CheckAccessToken(c *gin.Context, tokenStr string) bool
	RefreshAccessToken(c *gin.Context, tokenStr string, startTime int64) (resp AccessTokenTicket, err error)
}

// 它将在实现类里面的Register里面被赋值。对于用户来说，它应该永远只使用该变量并且确保在系统初始化的过程中，正确指定了token的实现
var Accessor TokenAccessor
