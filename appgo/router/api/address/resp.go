package address

type RespAddressInfo struct {
	Id            int    `json:"id" form:"id" gorm:"primary_key"`       // 地址ID
	Truename      string `json:"truename" form:"truename" `             // 会员姓名
	ProvinceId    int    `json:"province_id" form:"province_id" `       // 省份id
	CityId        int    `json:"city_id" form:"city_id" `               // 市级ID
	AreaId        int    `json:"area_id" form:"area_id" `               // 地区ID
	Address       string `json:"address" form:"address" `               // 地址
	CombineDetail string `json:"combine_detail" form:"combine_detail" ` // 地区内容组合
	MobilePhone   string `json:"mobile_phone" form:"mobile_phone" `     // 手机电话
	IsDefault     int    `json:"is_default" form:"is_default" `         // 1默认收货地址
	Type          string `json:"type" form:"type" `                     // &#39;个人&#39;,&#39;公司&#39;,&#39;其他&#39;....
}
