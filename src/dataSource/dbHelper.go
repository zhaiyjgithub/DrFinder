package dataSource

import (
	"DrFinder/src/conf"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

var (
	masterEngine *gorm.DB
	lock sync.Mutex
)

func InstanceMaster() *gorm.DB {
	if masterEngine != nil {
		return masterEngine
	}

	lock.Lock()
	defer lock.Unlock()

	if masterEngine != nil {
		return masterEngine
	}

	c:= conf.MasterDBConf
	driveSource:= fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?",
		c.User, c.Pwd, c.Host, c.Port, c.DBName)

	fmt.Println(driveSource)
	engine, err := gorm.Open(conf.DriverName,
		driveSource)

	if err != nil {
		log.Fatal("dbhelper instance error")
	}

	engine.LogMode(true)

	return engine
}