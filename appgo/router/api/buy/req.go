package buy

type ReqBuyCalculate struct {
	AddressId int   `json:"addressId"`
	CartIds   []int `json:"cartIds"`
}

type Freight struct {
	ComSkuId   int     `json:"comSkuId"`
	FreightFee float64 `json:"freightFee"`
	FreightWay string  `json:"freightWay"`
}

type RespBuyCalculate struct {
	ComAmount          float64   `json:"comAmount"`
	PayAmount          float64   `json:"payAmount"`
	ComFreightList     []Freight `json:"comFreightList"`
	FreightUnifiedFee  float64   `json:"freightUnifiedFee"`
	FreightTemplateFee float64   `json:"freightTemplateFee"`
	PayFreightFee      float64   `json:"payFreightFee"`
	SubTotal           float64   `json:"subTotal"`
	ComNum             int       `json:"comNum"`
}
