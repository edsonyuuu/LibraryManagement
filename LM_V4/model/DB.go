package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 连接DB配置在core包中，此文件已被废弃
var DB *gorm.DB

func MySql() {
	username := "root"            //账号
	password := "123456"          //密码
	host := "192.168.16.129"      //数据库地址，可以是IP
	port := 3307                  //数据库端口
	Dbname := "librarymanagement" //数据库名
	timeout := "10s"              //连接超时，10s

	var MysqlLogger logger.Interface
	//要显示的日志等级
	MysqlLogger = logger.Default.LogMode(logger.Info)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: MysqlLogger,
	})

	//设置数据库最大连接数
	if err != nil {
		panic("连接数据库失败，err=" + err.Error())
	}
	DB = db
	//err = DB.AutoMigrate(&SendMsg{}, &User{})
	//if err != nil {
	//	fmt.Printf("创建表结构失败err:%+v\n", err.Error())
	//}
	//err2 := DB.AutoMigrate(&User{}, &Book{}, &Category{}, &Librarian{}, &Record{})
	//err2 := DB.AutoMigrate(&Book{})
	//if err2 != nil {
	//	fmt.Println(err2)
	//}
	//连接成功
	//DB.AutoMigrate(Book{})
	//var bookinfo []BookInfo
	//sql := "select * from book_info"
	//err = DB.Raw(sql).Find(&bookinfo).Error
	//if err != nil {
	//	fmt.Println("err", err)
	//}
	//fmt.Println(len(bookinfo))
	//
	//for i := 0; i < len(bookinfo); i++ {
	//	fmt.Println(i)
	//}
	//i := 0
	//k := 100
	//for i < len(bookinfo) {
	//	rand.Seed(time.Now().UnixNano())
	//	// 生成一个0到10之间的随机整数
	//	randomNumber := rand.Intn(10) + 1
	//	tx := DB.Begin()
	//	sql1 := "insert into `book` (bn,name,description,count,category_id,img_url) values (?,?,?,?,?,?)"
	//	err1 := tx.Exec(sql1, bookinfo[i].ISBN, bookinfo[i].BookName, bookinfo[i].BriefIntroduction, k, randomNumber, bookinfo[i].ImgUrl).Error
	//	if err1 != nil {
	//		fmt.Println(err1)
	//		tx.Rollback()
	//	}
	//	tx.Commit()
	//	i++
	//	if i == len(bookinfo) {
	//		break
	//	}
	//
	//}

}

func Close() {
	db, _ := DB.DB()
	_ = db.Close()
}
