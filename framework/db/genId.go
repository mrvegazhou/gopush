package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopush/conf"
	"gopush/framework/config"
	"time"
)

type Lid struct {
	db         *gorm.DB    // 数据库连接
	businessId string     // 业务id
	ch         chan int64 // id缓冲池
	min, max   int64      // id段最小值，最大值
}

type GenId struct {
	Id		int64     `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessId	string 	`json:"business_id"`
	MaxId	int64	`json:"max_id"`
	Step	int		`json:"step"`
	Description	string	`json:"description"`
	CreateTime	int64	`json:"create_time"`
	UpdateTime    int64 	 `json:"update_time"`
}

func (GenId) TableName() string {
	return "t_gen_id"
}

func NewLid(db *gorm.DB, businessId string, len int) (*Lid, error) {
	lid := Lid{
		db:         db,
		businessId: businessId,
		ch:         make(chan int64, len),
	}
	go lid.productId()
	return &lid, nil
}

func (l *Lid) Get() int64 {
	return <-l.ch
}

func (l *Lid) productId() {
	l.reset()
	for {
		if l.min >= l.max {
			l.reset()
		}
		l.min++
		l.ch <- l.min
	}
}

func (l *Lid) reset() {
	for {
		err := l.getFromDB()
		if err == nil {
			return
		}
		time.Sleep(time.Second)
		continue
	}
}

func (l *Lid) getFromDB() error {
	var (
		maxId int64
		step  int64
	)
	tx := l.db.Begin()
	var genId GenId
	err := l.db.Raw("SELECT id,business_id,max_id,step,description,create_time,update_time FROM t_gen_id WHERE business_id = ? for update", l.businessId).Scan(&genId).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}
	err = l.db.Model(genId).Where("business_id = ?", l.businessId).Update("max_id", maxId+step).Error
	if err!=nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}
	tx.Commit()
	return nil
}

var LidGenId *Lid

func init() {
	conf := new(conf.MainConfig)
	err := Config.Load(&conf, false, "../../conf/app.yml")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	session := ConnectDB(conf)
	LidGenId, err = NewLid(session.DB, "message_id", 1000)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}