package common

import (
	"fmt"

	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"shortlinks/config"
)

var engine *xorm.Engine
var mutex = new(sync.RWMutex)

func RegisterMysql() *xorm.Engine {
	mutex.RLock()
	if engine != nil {
		mutex.RUnlock()
		return engine
	} else {
		mutex.RUnlock()
		mutex.Lock()

		if engine == nil {
			var dns string

			db_host := config.GetConf("db.host", "115.29.5.235")
			db_port := config.GetConf("db.port", "3306")
			db_user := config.GetConf("db.user", "linghitadmin")
			db_pass := config.GetConf("db.pass", "linghitmmclick2011")
			db_name := config.GetConf("db.name", "datacenter_core")

			dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)

			var err error
			engine, err = xorm.NewEngine("mysql", dns)
			if err != nil {
				fmt.Println("连接数据库失败")
			}

			engine.SetMaxIdleConns(10)  //连接池的空闲数大小
			engine.SetMaxOpenConns(100) //最大打开连接数
		}

		mutex.Unlock()
		return engine
	}
}
