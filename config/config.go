package config

import "flag"

const (
	db_IP       = "192.168.80.200"
	db_port     = 3306
	db_database = "manage_system"
	db_user     = "root"
	db_password = "root"
	log_path    = "manage_system.log"
	app_port    = 80
)

var (
	DB_IP       string
	DB_PORT     int
	DB_DATABASE string
	DB_USER     string
	DB_PASSWORD string
	LOG_PATH    string
	APP_PORT    int
)

func init() {
	flag.StringVar(&DB_IP, "h", db_IP, "数据库 ip 地址")
	flag.IntVar(&DB_PORT, "P", db_port, "数据库端口号")
	flag.StringVar(&DB_DATABASE, "d", db_database, "数据库名称")
	flag.StringVar(&DB_USER, "u", db_user, "数据库用户名")
	flag.StringVar(&DB_PASSWORD, "p", db_password, "数据库密码")
	flag.StringVar(&LOG_PATH, "l", log_path, "日志路径")
	flag.IntVar(&APP_PORT, "a", app_port, "监听端口")
	flag.Parse()
}
