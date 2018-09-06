package models

import (
	"time"
	"short_links/config"
)

type Application struct {
	Id uint32 `xorm:"pk autoincr"`
	Name string `xorm:"varchar(20)"`
	AppKey string `xorm:"varchar(100)"`
	AppSecret string `xorm:"varchar(100)"`
	Status int8 `xorm:"tinyint"`
	CreatedTime time.Time `xorm:"timestamp"`
	UpdatedTime time.Time `xorm:"timestamp"`
}

func (a *Application) TableName() string {
	return config.GetConf("db.table_prefix", "") + "applications"
}