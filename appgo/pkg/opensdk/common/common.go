package common

const (
	// UnifiedOrderURL 微信统一下单
	UnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419317853&token=&lang=zh_CN
const (
	//https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
	// AccessTokenURL code获取access_token
	AccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token"

	// RefreshTokenURL 重新获取access_token
	RefreshTokenURL = "https://api.weixin.qq.com/sns/oauth2/refresh_token"

	// UserInfoURL 通过access_token获取userInfo
	UserInfoURL = "https://api.weixin.qq.com/sns/userinfo"

	// CheckAccessTokenURL 检验授权凭证（access_token）是否有效
	CheckAccessTokenURL = "https://api.weixin.qq.com/sns/auth"

	// JsCode2SessionURL 临时登录凭证校验接口
	JsCode2SessionURL = "https://api.weixin.qq.com/sns/jscode2session"

	// SendRedPackURL 发送现金红包
	SendRedPackURL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack"
)

// 错误的信息
const (
	ErrAccessTokenEmpty  = "access token is empty"
	ErrAppIDEmpty        = "appid empty"
	ErrRefreshTokenEmpty = "refresh token is empty"
	ErrOpenIDEmpty       = "openid is empty"
	ErrCertCertEmpty     = "cert path is empty "
)
