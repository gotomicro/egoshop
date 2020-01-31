package constx

const (
	// 会员角色常量
	UserDefaultRole  = 0 // 系统错误角色
	UserCommonRole   = 1 // 普通角色
	UserVipRole      = 2 // vip角色
	UserSuperVipRole = 3 // super vip 角色
	// 积分增加常量
)

var (
	RoleNameMap = map[int]string{
		UserDefaultRole:  "系统错误",
		UserCommonRole:   "普通",
		UserVipRole:      "VIP",
		UserSuperVipRole: "SVIP",
	}

	PointSignin = PointConfig{
		Id:    1,
		Point: 0,
	}

	PointComment = PointConfig{
		Id:    2,
		Point: 1, // 评论积分 +1
	}
	PointStar = PointConfig{
		Id:    3,
		Point: 1, // 点赞积分 +1
	}
	PointShare = PointConfig{
		Id:    4,
		Point: 1, // 分享积分 +1
	}
	PointQuestion = PointConfig{
		Id:    5,
		Point: 1, // 提问积分 +1
	}
)

// todo 后面可以放到数据库
type PointConfig struct {
	Id    int // type id 类型
	Point int // 积分数
}

// todo 后面可以放到数据库
// todo 看要不要支持直接支付
type CostConfig struct {
	Id    int // type id 类型
	Point int // 积分
}
