package mysql

import (
	"time"
)

type Doc struct {
	Id            int        `json:"id" form:"id" gorm:"primary_key"`           // 主键ID
	CreatedAt     time.Time  `json:"created_at" form:"created_at" `             // 创建时间
	UpdatedAt     time.Time  `json:"updated_at" form:"updated_at" `             // 更新时间
	DeletedAt     *time.Time `json:"deleted_at" form:"deleted_at" gorm:"index"` // 删除时间
	CntView       int        `json:"cnt_view" form:"cnt_view" `                 // 阅读量
	CntStar       int        `json:"cnt_star" form:"cnt_star" `                 // 点赞量
	CntShare      int        `json:"cnt_share" form:"cnt_share" `               // 转发量
	CntComment    int        `json:"cnt_comment" form:"cnt_comment" `           // 评论量
	CreatedBy     int        `json:"created_by" form:"created_by" `             // 创建者
	UpdatedBy     int        `json:"updated_by" form:"updated_by" `             // 更新者
	Name          string     `json:"name" form:"name" `                         // 名称
	Desc          string     `json:"desc" form:"desc" `                         // 描述
	Cover         string     `json:"cover" form:"cover" `                       // 封面
	CommentStatus string     `json:"comment_status" form:"comment_status" `     //
	Nickname      string     `json:"nickname" form:"nickname" `                 // 昵称，冗余字段
	Avatar        string     `json:"avatar" form:"avatar" `                     // 头像，冗余字段
	Content       string     `json:"content" form:"content" `                   // 资料完整内容，要么是Img数组，要么是长度为1的Pdf、Doc数组
	Hot           int        `json:"hot" form:"hot" `                           // 热度
	Score         int        `json:"score" form:"score" `                       // 评分
	CntPay        int        `json:"cnt_pay" form:"cnt_pay" `                   // 付款量
	FreeType      int        `json:"free_type" form:"free_type" `               // 免费类型[0为初始值不应该存在的值 1为完全免费 2为分享免费 3只能积分购买，4只能现金购买，5现金，积分都可以购买]
	FreePage      int        `json:"free_page" form:"free_page" `               // 前多少页免费看，后多少页则需分享或付费
	NeedPoint     int        `json:"need_point" form:"need_point" `             // 所需要积分
	NeedAmount    int        `json:"need_amount" form:"need_amount" `           // 所需要金额，单位为分
	DownloadUrl   string     `json:"download_url" form:"download_url" `         // 下载url
	Status        int        `json:"status" form:"status" `                     // 状态[0为审核中 1为审核成功 2位审核失败]
	Cid1          int        `json:"cid1" form:"cid1" `                         // 1级栏目id
	Cid2          int        `json:"cid2" form:"cid2" `                         // 2级栏目id
	CateId        int        `json:"cate_id" form:"cate_id" `                   // 分类id[1纯文档 2试卷 3习题答案 4学霸笔记 5PPT 6视频 7实体]
	MainResType   int        `json:"main_res_type" form:"main_res_type" `       // 主体资源类型[1pdf 2img 3ppt 4word 5txt 6markdown 7excel]
	Document      string     `json:"document" form:"document" `                 //

}

func (*Doc) TableName() string {
	return "doc"
}
