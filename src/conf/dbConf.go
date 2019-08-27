package conf

const DriverName = "mysql"

type DBConf struct {
	Host string
	Port int
	User string
	Pwd string
	DBName string
}

var MasterDBConf = DBConf{
	Host:"127.0.0.1",
	Port: 3306,
	User:"root",
	Pwd: "12345678",
	DBName:"drfinder",
}