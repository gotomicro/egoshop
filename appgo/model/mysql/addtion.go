package mysql

type FeedImagesJson []string

type ComGalleryJson []string

type ComSpecListJson []ComSpecOneInfo

type ComSpecOneInfo struct {
	Id        int               `json:"id"`
	Name      string            `json:"name"`
	ValueList []CreateSpecValue `json:"valueList"`
}

type CreateSpecValue struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ComCidsJson []int

type ComImageSpecImagesJson StringsJson

type StringsJson []string

type ComBodyJson []OneComBody

type OneComBody struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type FreightAreasJson []FreightAreas

type FreightAreas struct {
	AreaIds          []int
	FirstAmount      float64
	FirstFee         float64
	AdditionalAmount float64
	AdditionalFee    float64
}
