package mysql

// TODO(从配置里面读取)
const TableName = "access_token"

type AccessToken struct {
	Jti        int    `gorm:"not null" json:"jti" form:"jti" gorm:"primary_key"`
	Sub        int    `gorm:"not null" json:"sub" form:"sub"`
	IaTime     int64  `gorm:"not null" json:"ia_time" form:"ia_time"`
	ExpTime    int64  `gorm:"not null;" json:"exp_time" form:"exp_time"`
	Ip         string `gorm:"not null" json:"ip" form:"ip"`
	CreateTime int64  `gorm:"not null" json:"create_time" form:"create_time"`
	IsLogout   int    `gorm:"not null;" json:"is_logout" form:"is_logout"`
	IsInvalid  int    `gorm:"not null;" json:"is_invalid" form:"is_invalid"`
	LogoutTime int64  `gorm:"not null" json:"logout_time" form:"logout_time"`
}

func (AccessToken) TableName() string {
	return TableName
}
