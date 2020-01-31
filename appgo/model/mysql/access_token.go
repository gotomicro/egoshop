package mysql

type AccessToken struct {
	Jti        int    `gorm:"not null" json:"jti" form:"jti" gorm:"primary_key"`
	Sub        int    `gorm:"not null" json:"sub" form:"sub"`
	IaTime     int64  `gorm:"not null" json:"ia_time" form:"ia_time"`
	ExpTime    int64  `gorm:"not null;comment:'过期时间'" json:"exp_time" form:"exp_time"`
	Ip         string `gorm:"not null" json:"ip" form:"ip"`
	CreateTime int64  `gorm:"not null" json:"create_time" form:"create_time"`
	IsLogout   int    `gorm:"not null;comment:'是否主动退出'" json:"is_logout" form:"is_logout"`
	IsInvalid  int    `gorm:"not null;comment:'是否作废，当修改密码后会作废[1作废 0没有作废]'" json:"is_invalid" form:"is_invalid"`
	LogoutTime int64  `gorm:"not null" json:"logout_time" form:"logout_time"`
}

func (*AccessToken) TableName() string {
	return "access_token"
}
