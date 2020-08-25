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
	Pwd: "123456",
	DBName:"drfinder",
}

var MongoDBConf = DBConf{
	Host: "localhost",
	Port: 27017,
	User: "zack",
	Pwd: "123456",
	DBName: "user",
}

var ElasticSearchConf = DBConf{
	Host:   "http://192.168.11.143",
	Port:   9200,
	User:   "",
	Pwd:    "",
	DBName: "",
}