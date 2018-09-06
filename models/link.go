package models

import (
	"time"
	"short_links/config"
)

type Link struct {
	Id uint64 `xorm:"pk autoincr"`
	LongLink string `xorm:"tinytext"`
	AppId uint32
	CreatedTime time.Time `xorm:"timestamp"`
	UpdatedTime time.Time `xorm:"timestamp"`
}

func (a *Link) TableName() string {
	return config.GetConf("db.table_prefix", "") + "links"
}