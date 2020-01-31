package code

const (
	MsgOk       = 0
	MsgErr      = 1
	MsgParamErr = 3

	DefaultErr = 10000

	LoginWechatErr  = 10001
	LoginWechatErr2 = 10002
	LoginWechatErr3 = 10003
	LoginWechatErr4 = 10004
	LoginWechatErr5 = 10005

	PayParamErr                   = 10101
	PayGoodsInfoErr               = 10102
	PayFreeTypeErr                = 10103
	PayUserGoodsErr               = 10104
	PayPrePayedErr                = 10105
	PayPayedErr                   = 10106
	PayOrderPayCreateErr          = 10107
	PayOrderCreateErr             = 10108
	PayUserGoodsCreateOrUpdateErr = 10109
	PayUserOpenErr                = 10110
	PayWechatPayErr               = 10111
	PayGoonPayedErr               = 10112
	PayGoonPayOrderInfoErr        = 10113

	DownloadGoodsNotExistErr    = 10201
	DownloadGoodsMysqlErr       = 10202
	DownloadNeedShareErr        = 10203
	DownloadNeedPointErr        = 10204
	DownloadNeedMoneyErr        = 10205
	DownloadNeedPointOrMoneyErr = 10206
	DownloadTypeErr             = 10207
	DownloadGoodsRedisErr       = 10208
	DownloadUserPrePayErr       = 10209

	DownloadGoodsRedis2Err = 10251

	PostGoodsEmptyFiles   = 10351
	PostGoodsInvalidFiles = 10352

	// 后台审核
	AdminAuditParamErr    = 10901
	AdminGoodsNotExistErr = 10902

	// 10000 ~ 10199
	UserAccountLengthErr   = 10001
	UserAccountNicknameErr = 10002
	UserAccountErr3        = 10003
	UserAccountErr4        = 10004
	UserAccountErr5        = 10005
	UserAccountErr6        = 10006
	UserAccountErr7        = 10007

	UserUpdateErr1 = 10100
	UserUpdateErr2 = 10101
	UserUpdateErr3 = 10102
	UserUpdateErr4 = 10103
	UserUpdateErr5 = 10104
	UserUpdateErr6 = 10105
	UserUpdateErr7 = 10106

	BookCreateErr0  = 10200
	BookCreateErr1  = 10201
	BookCreateErr2  = 10202
	BookCreateErr3  = 10203
	BookCreateErr4  = 10204
	BookCreateErr5  = 10205
	BookCreateErr6  = 10206
	BookCreateErr7  = 10207
	BookCreateErr8  = 10208
	BookCreateErr9  = 10209
	BookCreateErr10 = 10210
	BookCreateErr11 = 10211
	BookCreateErr12 = 10212

	DocumentContentAuto   = 10301
	DocumentContentTrue   = 10302
	DocumentContentPost3  = 10303
	DocumentContentPost4  = 10304
	DocumentContentPost5  = 10305
	DocumentContentPost6  = 10306
	DocumentContentPost7  = 10307
	DocumentContentPost8  = 10308
	DocumentContentPost9  = 10309
	DocumentContentPost10 = 10310

	UploadCover1 = 10401
	UploadCover2 = 10402
	UploadCover3 = 10403
	UploadCover4 = 10404
	UploadCover5 = 10405
	UploadCover6 = 10406
	UploadCover7 = 10407
	UploadCover8 = 10408
	UploadCover9 = 10409

	// feed upload error
	FeedUploadErr1 = 10501
	FeedUploadErr2 = 10502
	FeedUploadErr3 = 10503
	FeedUploadErr4 = 10504
	FeedUploadErr5 = 10505
	FeedUploadErr6 = 10506
	FeedUploadErr7 = 10507
	FeedUploadErr8 = 10508

	WechatOpenCallbackErr1  = 10601
	WechatOpenCallbackErr2  = 10602
	WechatOpenCallbackErr3  = 10603
	WechatOpenCallbackErr4  = 10604
	WechatOpenCallbackErr5  = 10605
	WechatOpenCallbackErr6  = 10606
	WechatOpenCallbackErr7  = 10607
	WechatOpenCallbackErr8  = 10608
	WechatOpenCallbackErr9  = 10609
	WechatOpenCallbackErr10 = 10610

	AdminLoginParamsErr   = 20001
	AdminLoginEmailErr    = 20002
	AdminLoginNicknameErr = 20003
	AdminLoginUpdateErr   = 20004

	AdminLogoutErr   = 20005
	AdminUserInfoErr = 20101

	PayCalculateNotFoundAddress = 30001
	PayCalculate2               = 30002
	PayCalculate3               = 30003

	FeedCartExist           = 40001
	FeedInfoSystemError     = 40002
	FeedInfoNotExist        = 40003
	FeedInfoStockNotEnough  = 40004
	FeedInfoCartCreateError = 40005
	ComSkuInfoSystemError   = 40006
	ComInfoSystemError      = 40007
	ComInfoNotExist         = 40008

	OssErr2 = 50001
	OssErr3 = 50002
)

var CodeMap = map[int]string{
	0:                     "成功",
	MsgParamErr:           "参数错误",
	10001:                 "登录微信失败",
	10112:                 "不是待支付商品",
	10201:                 "商品不存在",
	10202:                 "系统错误",
	10203:                 "分享后,可以下载",
	10204:                 "积分购买后,可以下载",
	10205:                 "现金购买后,可以下载",
	10206:                 "积分或者现金购买后,可以下载",
	10207:                 "类型错误",
	10208:                 "系统错误",
	10209:                 "已经有未支付订单,请及时支付",
	10251:                 "系统错误",
	PostGoodsEmptyFiles:   "上传资源为空",
	PostGoodsInvalidFiles: "上传资源格式错误",
	AdminLoginParamsErr:   "参数错误",
	AdminLoginEmailErr:    "邮箱错误",
	AdminLoginNicknameErr: "名称错误",
	AdminLoginUpdateErr:   "更新失败",
	AdminLogoutErr:        "登出失败",
	AdminUserInfoErr:      "获取用户信息失败",

	PayCalculateNotFoundAddress: "请填写收货地址",

	FeedCartExist:           "购物车已经存在该商品",
	FeedInfoSystemError:     "商品宠物信息系统错误",
	FeedInfoNotExist:        "不存在该商品信息",
	FeedInfoStockNotEnough:  "宠物库存不足",
	FeedInfoCartCreateError: "宠物购物车创建失败",
	ComSkuInfoSystemError:   "商品sku信息系统错误",
	ComInfoSystemError:      "商品信息系统错误",
	ComInfoNotExist:         "商品信息不存在",
}
