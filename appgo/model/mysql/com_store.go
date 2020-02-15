package mysql

import "time"

type ComStore struct {
	CreatedAt   time.Time `json:"createdAt"` // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"` // 更新时间
	ComId       int       `gorm:"not null;comment:'注释'"json:"comId"`
	PreMarkdown string    `gorm:"not null;type:text;comment:'注释'"json:"preMarkdown"`
	PreHtml     string    `gorm:"not null;type:text;comment:'注释'"json:"preHtml"`
	Markdown    string    `gorm:"not null;type:text;comment:'注释'"json:"markdown"`
	Html        string    `gorm:"not null;type:text;comment:'注释'"json:"webHtml"`
	WechatHtml  string    `gorm:"not null;type:text;comment:'注释'"json:"wechatHtml"`
	CreatedBy   int       `gorm:"not null;comment:'注释'"json:"createdBy"`
	UpdatedBy   int       `gorm:"not null;comment:'注释'"json:"updatedBy"`
	PartId      int       `gorm:"not null"json:"partid"` //部门分类
}

func (*ComStore) TableName() string {
	return "com_store"
}
