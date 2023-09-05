package core

import (
	"LibraryManagementV1/LM_V4/global"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func InitGorm() {
	if global.Config.Mysql.Host == "" {
		log.Println("没有连接Mysql")
		return
	}
	dsn := global.Config.Mysql.Dsn()
	var mysqlLogger logger.Interface
	if global.Config.System.Env == "debug" {
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		log.Fatalf("dsn:%+v连接失败", dsn)
	}
	fmt.Println("数据库已连接")
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10) // 最大连接数
	sqlDB.SetMaxOpenConns(100)
	global.DB = db
}
