package service

import (
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/mus"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/jinzhu/gorm"
)

type queueView struct {
	ch chan ViewTask
}

// Task defines a unit of work in the queue.
type ViewTask struct {
	// ID identifies this task.
	GoodsId int `json:"id,omitempty"`
	Uid     int
	TypeId  int
	Name    string
}

func InitQueueView() *queueView {
	g := &queueView{
		ch: make(chan ViewTask, 1000),
	}
	go func() {
		for {
			if err := g.run(); err != nil {
				// todo log
			}
		}
	}()
	return g
}

func (w *queueView) Push(t ViewTask) {
	w.ch <- t
}

// todo 注意防双击问题，应用层处理，这里不做处理
func (q *queueView) run() (err error) {
	task := <-q.ch
	tx := mus.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	var goodsData mysql.UserGoods
	err = tx.Select("*").Where("uid = ? and goods_id = ? and type_id = ?", task.Uid, task.GoodsId, task.TypeId).Find(&goodsData).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	// 如果找不到就创建
	if err == gorm.ErrRecordNotFound {
		err = tx.Create(&mysql.UserGoods{
			Uid:     task.Uid,
			GoodsId: task.GoodsId,
			TypeId:  task.TypeId,
			Name:    task.Name,
			IsRead:  1,
		}).Error
		return
	}
	// 如果找到了就更新
	err = tx.Model(mysql.UserGoods{}).Where("uid = ? and goods_id = ? and type_id = ?", task.Uid, task.GoodsId, task.TypeId).
		Updates(map[string]interface{}{
			"is_read": 1,
			"name":    task.Name,
		}).Error

	return
}
