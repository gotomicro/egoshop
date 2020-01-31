package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type OrderExtend struct {
	Id                 int                         `gorm:"not null;primary_key"json:"id"`         // 订单索引id
	TrackingTime       int64                       `gorm:"not null;"json:"trackingTime"`          // 配送时间
	TrackingNo         string                      `gorm:"not null;"json:"trackingNo"`            // 物流单号
	ShipperId          int                         `gorm:"not null;"json:"shipperId"`             // 商家物流地址id
	ExpressId          int                         `gorm:"not null;"json:"expressId"`             // 物流公司id，默认为0 代表不需要物流
	Message            string                      `gorm:"not null;"json:"message"`               // 买家留言
	VoucherPrice       int                         `gorm:"not null;"json:"voucherPrice"`          // 代金券面额
	VoucherId          int                         `gorm:"not null;"json:"voucherId"`             // 代金券id
	VoucherCode        string                      `gorm:"not null;"json:"voucherCode"`           // 代金券编码
	Remark             string                      `gorm:"not null;"json:"remark"`                // 发货备注
	ReceiverName       string                      `gorm:"not null;"json:"receiverName"`          // 收货人姓名
	ReceiverPhone      string                      `gorm:"not null;"json:"receiverPhone"`         // 收货人电话
	ReceiverInfo       OrderExtendReceiverInfoJson `gorm:"not null;type:json"json:"receiverInfo"` // 收货人其它信息 json
	ReceiverProvinceId int                         `gorm:"not null;"json:"receiverProvinceId"`
	ReceiverCityId     int                         `gorm:"not null;"json:"receiverCityId"`
	ReceiverAreaId     int                         `gorm:"not null;"json:"receiverAreaId"`
	InvoiceInfo        OrderExtendInvoiceInfoJson  `gorm:"not null;type:json"json:"invoiceInfo"` // 发票信息 json
	PromotionInfo      string                      `gorm:"not null;"json:"promotionInfo"`        // 促销信息备注
	EvaluateTime       int64                       `gorm:"not null;"json:"evaluateTime"`         // 评价时间
	ServiceRemarks     string                      `gorm:"not null;"json:"serviceRemarks"`       // 后台客服对此订单做出的备注
	DeletedAt          *time.Time                  `gorm:"index"json:"deleteTime"`               // 软删除时间
	DeliverName        string                      `gorm:"not null;"json:"deliverName"`          // 配送人名字
	DeliverPhone       string                      `gorm:"not null;"json:"deliverPhone"`         // 配送人电话
	DeliverAddress     string                      `gorm:"not null;"json:"deliverAddress"`       // 配送地址
	FreightRule        int                         `gorm:"not null;"json:"freightRule"`          // 运费规则1按商品累加运费2组合运费
	NeedExpress        int                         `gorm:"not null;"json:"needExpress"`          // 是否需要物流1需要0不需要
	IsRefund           int                         `gorm:"-"json:"isRefund"`
	RefundStatus       int                         `gorm:"-"json:"refundStatus"`
}

func (*OrderExtend) TableName() string {
	return "order_extend"
}

type OrderExtendReceiverInfoJson struct {
	Name    string `json:"name"`
	Detail  string `json:"detail"`
	Phone   string `json:"phone"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

func (c OrderExtendReceiverInfoJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *OrderExtendReceiverInfoJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

type OrderExtendInvoiceInfoJson []string

func (c OrderExtendInvoiceInfoJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *OrderExtendInvoiceInfoJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
