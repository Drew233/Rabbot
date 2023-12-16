package mysql

import (
	"fmt"
	// "time"
	// "strings"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"rabbot/config"
	"rabbot/internal/log"
)

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", 
					   config.RabConfig.MysqlConfig.Username, 
					   config.RabConfig.MysqlConfig.Password, 
					   config.RabConfig.MysqlConfig.Ip, 
					   config.RabConfig.MysqlConfig.Port, 
					   config.RabConfig.MysqlConfig.Dbname, 
					   config.RabConfig.MysqlConfig.Charset)
					   
	log.RabLog.Info(dsn)
	//Open打开一个driverName指定的数据库，dataSourceName指定数据源
	//不会校验用户名和密码是否正确，只会对dsn的格式进行检测
	db, err := sql.Open("mysql", dsn)
	if err != nil { //dsn格式不正确的时候会报错
		log.RabLog.Error("校验失败,err: ", err)
		return
	}
	//尝试与数据库连接，校验dsn是否正确
	err = db.Ping()
	if err != nil {
		log.RabLog.Error("校验失败,err: ", err)
		return
	}
	// 设置最大连接数
	db.SetMaxOpenConns(50)

	log.RabLog.Info("连接数据库成功！")
	return
}