package models

import (
	"shortlinks/config"
	"time"
)

type Log struct {
	Id          uint64 `xorm:"pk autoincr"`
	AppId       uint32
	LongLink    string    `xorm:"tinytext"`
	Status      int8      `xorm:"tinyint"`
	StatusText  string    `xorm:"varchar(255)"`
	Persist     int8      `xorm:"tinyint"`
	CreatedTime time.Time `xorm:"timestamp"`
}

func (a *Log) TableName() string {
	return config.GetConf("db.table_prefix", "") + "logs"
}
