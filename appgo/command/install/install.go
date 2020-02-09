package install

import (
	"fmt"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/pkg/util"
	"github.com/goecology/egoshop/appgo/service"
	"github.com/goecology/muses/pkg/cache/redis"
	mmysql "github.com/goecology/muses/pkg/database/mysql"
	"github.com/goecology/muses/pkg/logger"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)
var Models []interface{}

func init()  {
	Models = []interface{}{
		&mysql.Banner{},
		&mysql.Com{},
		&mysql.ComSku{},
		&mysql.ComStore{},
		&mysql.ComCate{},
		&mysql.ComRelateCate{},
		&mysql.Attachment{},
		&mysql.Cart{},
		&mysql.Address{},
		&mysql.AddressType{},
		&mysql.Order{},
		&mysql.OrderExtend{},
		&mysql.OrderLog{},
		&mysql.OrderPay{},
		&mysql.OrderGoods{},
		&mysql.ComImage{},
		&mysql.ComSpec{},
		&mysql.ComSpecValue{},
		&mysql.Freight{},
		&mysql.User{},
		&mysql.UserGoods{},
		&mysql.UserOpen{},
		&mysql.AccessToken{},
		&mysql.Comment{},
		&mysql.AdminUser{},
		&mysql.Signin{},
		&mysql.SigninLog{},
		&mysql.PointLog{},
		&mysql.PointLimit{},
		&mysql.Image{},
	}
}
func Create(isClear bool) error {
	db := mmysql.Caller("egoshop")
	if isClear {
		db.DropTable(Models...)
	}

	db.SingularTable(true)
	if db.Error != nil {
		return db.Error
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(Models...)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func Mock() error  {
	mock(
		mmysql.Caller("egoshop"),
		viper.GetString("oss.endpoint"),
		viper.GetString("oss.accessKeyID"),
		viper.GetString("oss.accessKeySecret"),
		viper.GetString("oss.bucket"),
		redis.Caller("egoshop"),
	)
	return nil
}


func mock(db *gorm.DB, endpoints, accessKeyId, accessKeySecret, bucketName string, rclient *redis.Client) {
	createCommodity(db, endpoints, accessKeyId, accessKeySecret, bucketName)
	createAddressType(db)
	createComCate(db)
	createAdminUser(db)
	createNew(db, rclient)
}

func createCommodity(db *gorm.DB, endpoints, accessKeyId, accessKeySecret, bucketName string) {
	var err error
	service.InitOssCli(endpoints, accessKeyId, accessKeySecret, bucketName, logger.DefaultLogger())

	key1_1 := service.Oss.Key("mock")
	key1_2 := service.Oss.Key("mock")
	key1_3 := service.Oss.Key("mock")
	key2_1 := service.Oss.Key("mock")
	key2_2 := service.Oss.Key("mock")
	key2_3 := service.Oss.Key("mock")
	key2_4 := service.Oss.Key("mock")

	err = service.Oss.PutObj(bucketName, key1_1, "./mockdata/1_1.jpg")
	if err != nil {
		panic(err)
	}
	err = service.Oss.PutObj(bucketName, key1_2, "./mockdata/1_2.jpg")
	if err != nil {
		panic(err)
	}
	err = service.Oss.PutObj(bucketName, key1_3, "./mockdata/1_3.jpg")
	if err != nil {
		panic(err)
	}
	err = service.Oss.PutObj(bucketName, key2_1, "./mockdata/2_1.jpg")
	if err != nil {
		panic(err)
	}
	err = service.Oss.PutObj(bucketName, key2_2, "./mockdata/2_2.jpg")
	if err != nil {
		panic(err)
	}
	err = service.Oss.PutObj(bucketName, key2_3, "./mockdata/2_3.jpg")
	if err != nil {
		panic(err)
	}
	err = service.Oss.PutObj(bucketName, key2_4, "./mockdata/2_4.jpg")
	if err != nil {
		panic(err)
	}
	// 初始化规格
	db.Create(&mysql.ComSpec{
		Id:        1,
		Name:      "颜色",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	db.Create(&mysql.ComSpecValue{
		Id:     1,
		SpecId: 1,
		Name:   "灰色",
	})

	comInfo := &mysql.Com{
		Id:           0,
		CreatedAt:    time.Date(2019, 10, 28, 23, 20, 31, 0, time.Local),
		UpdatedAt:    time.Date(2019, 10, 28, 23, 20, 31, 0, time.Local),
		DeletedAt:    nil,
		CntView:      0,
		CntStar:      0,
		CntCollect:   0,
		CntShare:     0,
		CntComment:   0,
		CntIsPay:     0,
		CreatedBy:    0,
		UpdatedBy:    0,
		Title:        "Ayuko 时髦气质毛呢外套 加棉可脱卸呢大衣",
		SubTitle:     "Ayuko 时髦气质毛呢外套 加棉可脱卸呢大衣",
		Cover:        key1_1,
		Gallery:      mysql.ComGalleryJson{key1_1, key1_2, key1_3},
		Stock:        100,
		SaleNum:      0,
		GroupSaleNum: 0,
		SaleTime:     time.Now(),
		PayType:      0,
		FreightFee:   0,
		FreightId:    0,
		BaseSaleNum:  0,
		IsOnSale:     1,
		ImageSpecId:  0,
		OriginPrice:  20,
		Price:        0.01,
		PreMarkdown:  "",
		PreHtml:      "",
		Markdown:     "",
		Html:         "",
		WechatHtml:   "",
		SkuList:      nil,
		Cids:         []int{1},
		SpecList: []mysql.ComSpecOneInfo{
			mysql.ComSpecOneInfo{
				Id:   1,
				Name: "颜色",
				ValueList: []mysql.CreateSpecValue{
					{
						Id:   1,
						Name: "灰色",
					},
				},
			},
		},
	}

	db.Create(comInfo)

	db.Create(&mysql.ComRelateCate{
		ComId: comInfo.Id,
		Cid:   1,
	})

	db.Create(&mysql.Banner{
		Id:        0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		CreatedBy: 0,
		UpdatedBy: 0,
		ResType:   0,
		Title:     "",
		Link:      "/pages/cominfo/main?id=1",
		Image:     key1_1,
		Content:   "",
		Enable:    0,
		StartTime: 0,
		EndTime:   0,
	})

	db.Create(&mysql.ComSku{
		CreatedBy: 0,
		UpdatedBy: 0,
		ComId:     1,
		Spec: []mysql.ComSkuSpecOneInfo{
			mysql.ComSkuSpecOneInfo{
				ID:        1,
				Name:      "颜色",
				ValueID:   1,
				ValueImg:  "",
				ValueName: "灰色",
			},
		},
		Price:        0.01,
		Stock:        100,
		Code:         "",
		Cover:        key1_1,
		Weight:       0,
		Title:        "灰色",
		SaleNum:      0,
		GroupSaleNum: 0,
	})

	// 初始化规格
	db.Create(&mysql.ComSpec{
		Id:        2,
		Name:      "颜色",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	db.Create(&mysql.ComSpecValue{
		Id:     2,
		SpecId: 2,
		Name:   "蓝色",
	})

	db.Create(&mysql.ComSpecValue{
		Id:     3,
		SpecId: 2,
		Name:   "白色",
	})

	db.Create(&mysql.Com{
		Id:           0,
		CreatedAt:    time.Date(2019, 10, 28, 23, 20, 31, 0, time.Local),
		UpdatedAt:    time.Date(2019, 10, 28, 23, 20, 31, 0, time.Local),
		DeletedAt:    nil,
		CntView:      0,
		CntStar:      0,
		CntCollect:   0,
		CntShare:     0,
		CntComment:   0,
		CntIsPay:     0,
		CreatedBy:    0,
		UpdatedBy:    0,
		Title:        "韩国高腰牛仔短裤女2019夏季新款韩版宽松显瘦a字毛边阔腿热裤女",
		SubTitle:     "韩国高腰牛仔短裤女2019夏季新款韩版宽松显瘦a字毛边阔腿热裤女",
		Cover:        key2_1,
		Gallery:      mysql.ComGalleryJson{key2_1, key2_2, key2_3, key2_4},
		Stock:        100,
		SaleNum:      0,
		GroupSaleNum: 0,
		SaleTime:     time.Now(),
		PayType:      0,
		FreightFee:   0,
		FreightId:    0,
		BaseSaleNum:  0,
		IsOnSale:     1,
		ImageSpecId:  0,
		OriginPrice:  10,
		Price:        0.01,
		PreMarkdown:  "",
		PreHtml:      "",
		Markdown:     "",
		Html:         "",
		WechatHtml:   "",
		Cids:         []int{2},
		SkuList:      nil,
		SpecList: []mysql.ComSpecOneInfo{
			mysql.ComSpecOneInfo{
				Id:   2,
				Name: "颜色",
				ValueList: []mysql.CreateSpecValue{
					{
						Id:   2,
						Name: "蓝色",
					},
					{
						Id:   3,
						Name: "白色",
					},
				},
			},
		},
	})

	db.Create(&mysql.Banner{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     "韩国高腰牛仔短裤女2019夏季新款韩版宽松显瘦a字毛边阔腿热裤女",
		Link:      "/pages/cominfo/main?id=2",
		Image:     key2_1,
	})

	db.Create(&mysql.ComSku{
		ComId: 2,
		Spec: []mysql.ComSkuSpecOneInfo{
			mysql.ComSkuSpecOneInfo{
				ID:        2,
				Name:      "颜色",
				ValueID:   2,
				ValueImg:  "",
				ValueName: "蓝色",
			},
		},
		Price:        0.01,
		Stock:        60,
		Code:         "",
		Cover:        key2_1,
		Weight:       0,
		Title:        "蓝色",
		SaleNum:      0,
		GroupSaleNum: 0,
	})

	db.Create(&mysql.ComSku{
		ComId: 2,
		Spec: []mysql.ComSkuSpecOneInfo{
			mysql.ComSkuSpecOneInfo{
				ID:        2,
				Name:      "颜色",
				ValueID:   2,
				ValueImg:  "",
				ValueName: "白色",
			},
		},
		Price:        0.01,
		Stock:        40,
		Code:         "",
		Cover:        key2_2,
		Weight:       0,
		Title:        "白色",
		SaleNum:      0,
		GroupSaleNum: 0,
	})

}

func createAddressType(db *gorm.DB) {
	db.Create(&mysql.AddressType{
		Name: "家",
	})
	db.Create(&mysql.AddressType{
		Name: "学校",
	})
	db.Create(&mysql.AddressType{
		Name: "公司",
	})
	db.Create(&mysql.AddressType{
		Name: "其他",
	})
}

func createComCate(db *gorm.DB) {
	db.Create(&mysql.ComCate{
		Name:   "女上装",
		Icon:   "/static/images/mall/category/5.png",
		Status: 1,
	})
	db.Create(&mysql.ComCate{
		Name:   "裙子",
		Icon:   "/static/images/mall/category/6.jpg",
		Status: 1,
	})
	db.Create(&mysql.ComCate{
		Name:   "皮包",
		Icon:   "/static/images/mall/category/9.jpg",
		Status: 1,
	})
	db.Create(&mysql.ComCate{
		Name:   "护肤品",
		Icon:   "/static/images/mall/category/8.jpg",
		Status: 1,
	})
}

type Info struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Cover       string  `json:"cover"`
	Price       float64 `json:"price"`
	OriginPrice float64 `json:"originPrice"`
	IsLabel     int     `json:"isLabel"`
	LabelIcon   string  `json:"labelIcon"`
}

func createAdminUser(db *gorm.DB) {
	pwdHash, err := util.Hash("egoshop")
	if err != nil {
		fmt.Println("err", err)
		return
	}
	user := mysql.User{
		Name:          "egoshop",
		Password:      pwdHash,
		Status:        1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastLoginIP:   "127.0.0.1",
		LastLoginTime: time.Now(),
	}
	if err = db.Create(&user).Error; err != nil {
		fmt.Println("err", err)
		return
	}
}

func createNew(db *gorm.DB, rclient *redis.Client) {
	var output []mysql.Com
	redisOutput := make([]Info, 0)
	db.Where("id in (?)", []int{1, 2}).Find(&output)
	for _, value := range output {
		redisOutput = append(redisOutput, Info{
			Id:          value.Id,
			Title:       value.Title,
			Cover:       value.Cover,
			Price:       value.Price,
			OriginPrice: value.OriginPrice,
			IsLabel:     1,
			LabelIcon:   "/static/images/mall/new/new.png",
		})
	}
	rclient.Set("home:new", redisOutput, 0)

}

