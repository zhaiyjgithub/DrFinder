package dataSource

import (
	"DrFinder/src/conf"
	"fmt"
	"github.com/go-xorm/xorm"
	"log"
	"sync"
	_ "github.com/go-sql-driver/mysql"
)

var (
	masterEngine *xorm.Engine
	slaveEngine *xorm.Engine
	lock sync.Mutex
)

func InstanceMaster() *xorm.Engine {
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
	engine, err := xorm.NewEngine(conf.DriverName,
		driveSource)

	fmt.Print(err)
	if err != nil {
		log.Fatal("dbhelper instance error")
	}

	engine.ShowSQL(true)

	return engine
}