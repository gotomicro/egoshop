package service

import (
	"time"

	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/model/constx"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/jinzhu/gorm"
)

type queuePoint struct {
	ch     chan PointTask
	limit1 int //todo 限制积分，写死
}

// Task defines a unit of work in the queue.
type PointTask struct {
	TypeData constx.PointConfig
	Uid      int
}

func InitQueuePoint() *queuePoint {
	g := &queuePoint{
		ch: make(chan PointTask, 1000),
		// todo 写死积分
		limit1: 10,
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

func (q *queuePoint) Push(t PointTask) {
	q.ch <- t
}

// todo 注意防双击问题，应用层处理，这里不做处理
func (q *queuePoint) run() (err error) {
	task := <-q.ch
	tx := mus.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	isLimit := q.isLimit(task)
	// 如果被限制，那么需要查下数据库，看下封顶值是多少
	var limitData mysql.PointLimit
	// 获取今天0点时间
	todayZeroStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", todayZeroStr)
	todayZero := t.Unix()

	// 受限制类型
	if isLimit {
		err = tx.Select("id,limit1,updated_at").Where("uid = ?", task.Uid).Find(&limitData).Error
		// 如果系统错误
		if err != nil && err != gorm.ErrRecordNotFound {
			return
		}

		// 如果查询为空
		if err != nil && err == gorm.ErrRecordNotFound {
			err = q.CreateQueuePoint(tx, task.Uid, task, isLimit)
			return
		}
		// 上次更新时间
		updatedUnix := limitData.Limit1UpdatedAt.Unix()
		// 如果存在有数据，如果时间不是今天，那么直接可以加积分，不受限制
		if updatedUnix < todayZero {
			// 因为不是当天，所以需要把限制重置为1
			err = q.UpdateQueuePoint(tx, task.Uid, task, isLimit, true)
			return
		}
		// 说明已经超过上限值，不在增加这个数据
		if limitData.Limit1 >= q.limit1 {
			return
		}
		// 当天时间，那么需要控制上限值
		err = q.UpdateQueuePoint(tx, task.Uid, task, isLimit, false)
		return

	}

	// 如果不是受限制类型
	err = q.UpdateQueuePoint(tx, task.Uid, task, isLimit, false)
	return
}

func (q *queuePoint) isLimit(data PointTask) bool {
	switch data.TypeData.Id {
	case constx.PointComment.Id:
		return true
	case constx.PointStar.Id:
		return true
	case constx.PointShare.Id:
		return true
	case constx.PointQuestion.Id:
		return true
	}
	return false
}

func (w *queuePoint) CreateQueuePoint(tx *gorm.DB, uid int, data PointTask, isLimit bool) (err error) {
	// 如果存在限制
	if isLimit {
		err = tx.Create(&mysql.PointLimit{
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			Uid:             uid,
			Limit1:          1,
			Limit1UpdatedAt: time.Now(),
		}).Error
		if err != nil {
			return
		}
	}

	err = tx.Create(&mysql.PointLog{
		TypeId:    constx.PointSignin.Id,
		CreatedAt: time.Now(),
		Uid:       uid,
		Point:     data.TypeData.Point,
	}).Error
	if err != nil {
		return
	}
	err = tx.Model(mysql.User{}).Where("id = ?", uid).Updates(mysql.Ups{
		"point": gorm.Expr("point+?", data.TypeData.Point),
	}).Error
	return
}

func (w *queuePoint) UpdateQueuePoint(tx *gorm.DB, uid int, data PointTask, isLimit bool, isReset bool) (err error) {
	// 如果存在限制
	if isLimit {
		if isReset {
			err = tx.Model(mysql.PointLimit{}).Where("uid = ?", uid).Updates(mysql.Ups{
				"updated_at":        time.Now(),
				"limit1_updated_at": time.Now(),
				"limit1":            1,
			}).Error
			if err != nil {
				return
			}
		} else {
			err = tx.Model(mysql.PointLimit{}).Where("uid = ?", uid).Updates(mysql.Ups{
				"updated_at":        time.Now(),
				"limit1_updated_at": time.Now(),
				"limit1":            gorm.Expr("limit1 + ?", 1),
			}).Error
			if err != nil {
				return
			}
		}
	}

	err = tx.Create(&mysql.PointLog{
		CreatedAt: time.Now(),
		TypeId:    constx.PointSignin.Id,
		Uid:       uid,
		Point:     data.TypeData.Point,
	}).Error
	if err != nil {
		return
	}

	err = tx.Model(mysql.User{}).Where("id = ?", uid).Updates(mysql.Ups{
		"point": gorm.Expr("point+?", data.TypeData.Point),
	}).Error
	return
}
