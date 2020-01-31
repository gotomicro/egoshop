package service

import (
	"time"

	"github.com/goecology/egoshop/appgo/model/constx"

	"errors"

	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/mus"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/jinzhu/gorm"
)

type queueSignin struct {
	ch chan Task
}

// Task defines a unit of work in the queue.
type Task struct {
	// ID identifies this task.
	ID  int `json:"id,omitempty"`
	Uid int
}

func InitQueueSignin() *queueSignin {
	g := &queueSignin{
		ch: make(chan Task, 1000),
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

func (w *queueSignin) Push(t Task) {
	w.ch <- t
}

// todo 注意防双击问题，应用层处理，这里不做处理
func (q *queueSignin) run() (err error) {
	task := <-q.ch
	tx := mus.Db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	var signinData mysql.Signin
	err = tx.Select("*").Where("uid = ?", task.Uid).Find(&signinData).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	// 获取今天0点时间
	todayZeroStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", todayZeroStr)
	todayZero := t.Unix()

	// 上次更新时间
	updatedUnix := signinData.UpdatedAt.Unix()

	// 如果今天已经签到了，就直接返回，但这种不应该出现，记录日志
	if updatedUnix > todayZero {
		err = errors.New("already signin")
		return
	}

	// 说明没有签到过
	if signinData.Id == 0 {
		err = q.CreateSignin(tx, task.Uid, 1, 1)
		return
	}

	// 如果签到过，要判断是否连续签到过
	// 如果在这个区间里说明连续签到了
	if (todayZero-24*60*60) < updatedUnix && todayZero >= updatedUnix {
		addPoint := signinData.Point + 1
		// 如果超过5，封顶为5的积分
		if addPoint > 5 {
			addPoint = 5
		}
		err = q.UpdateSignin(tx, task.Uid, addPoint, signinData.SigninCnt+1)
		return
	}

	// 普通签到，从1开始
	err = q.UpdateSignin(tx, task.Uid, 1, 1)
	return nil
}

func (w *queueSignin) CreateSignin(tx *gorm.DB, uid int, point int, signinCnt int) (err error) {
	err = tx.Create(&mysql.Signin{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Uid:       uid,
		Point:     point,
		SigninCnt: signinCnt,
	}).Error
	if err != nil {
		return
	}
	// 说明是第一天签到
	err = tx.Create(&mysql.SigninLog{
		CreatedAt: time.Now(),
		Uid:       uid,
		Point:     point,
		SigninCnt: signinCnt,
	}).Error
	if err != nil {
		return
	}

	err = tx.Create(&mysql.PointLog{
		TypeId:    constx.PointSignin.Id,
		CreatedAt: time.Now(),
		Uid:       uid,
		Point:     point,
	}).Error
	if err != nil {
		return
	}
	err = tx.Model(mysql.User{}).Where("id = ?", uid).Updates(mysql.Ups{
		"point": gorm.Expr("point+?", point),
	}).Error
	return
}

func (w *queueSignin) UpdateSignin(tx *gorm.DB, uid int, point int, signinCnt int) (err error) {
	err = tx.Model(mysql.Signin{}).Where("uid = ?", uid).Updates(mysql.Ups{
		"updated_at": time.Now(),
		"point":      point,
		"signin_cnt": signinCnt,
	}).Error
	if err != nil {
		return
	}
	// 说明是第一天签到
	err = tx.Create(&mysql.SigninLog{
		CreatedAt: time.Now(),
		Uid:       uid,
		Point:     point,
		SigninCnt: signinCnt,
	}).Error
	if err != nil {
		return
	}

	err = tx.Create(&mysql.PointLog{
		CreatedAt: time.Now(),
		TypeId:    constx.PointSignin.Id,
		Uid:       uid,
		Point:     point,
	}).Error
	if err != nil {
		return
	}

	err = tx.Model(mysql.User{}).Where("id = ?", uid).Updates(mysql.Ups{
		"point": gorm.Expr("point+?", point),
	}).Error
	return
}
