package standard

import (
	"github.com/gin-gonic/gin"
)

// TODO(从配置中读取)
const AccessTokenExpireInterval = 7 * 24 * 60 * 60

type AccessTokenTicket struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

type TokenAccessor interface {
	CreateAccessToken(c *gin.Context, uid int, startTime int64) (resp AccessTokenTicket, err error)
	CheckAccessToken(c *gin.Context, tokenStr string) bool
	RefreshAccessToken(c *gin.Context, tokenStr string, startTime int64) (resp AccessTokenTicket, err error)
	DecodeAccessToken(tokenStr string) (resp map[string]interface{}, err error)
}
