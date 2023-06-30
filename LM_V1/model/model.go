package model

import (
	"time"
)

type Book struct {
	Id          int64
	BN          string `gorm:"type:varchar(20)" json:"bn"`
	Name        string `gorm:"type:varchar(200)" json:"name"`
	Description string `gorm:"type:varchar(15000)"`
	Count       int    `json:"count"`
	CategoryId  int64  `json:"categoryId"`
	ImgUrl      string `json:"img_url" gorm:"varchar(200)"`
}
type Category struct {
	Id   int64
	Name string `gorm:"type:varchar(100)"`
	//Book []*Book `gorm:"foreignKey=CategoryId"`
}
type User struct {
	Id       int64  `json:"id" form:"id"`
	UserName string `json:"user_name" form:"user_name" gorm:"type:varchar(100)"`
	Password string `json:"password" form:"password" gorm:"type:varchar(100)"`
	Name     string `json:"name" form:"name" gorm:"type:varchar(100)"`
	Sex      string `json:"sex" form:"sex" gorm:"type:varchar(100)"`
	Phone    string `json:"phone" form:"phone" gorm:"type:varchar(100)"`
	Status   int    `json:"status" form:"status"` //`json:""默认正常0 封禁1
}
type Librarian struct {
	Id       int64
	UserName string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100)"`
	Name     string `gorm:"type:varchar(100)"`
	Sex      string `gorm:"type:varchar(100)"`
	Phone    string `gorm:"type:varchar(100)"`
}
type Record struct {
	Id         int64
	UserId     int64
	BookId     int64
	Status     int //已归还1 未归还0
	StartTime  time.Time
	OverTime   time.Time
	ReturnTime time.Time
}

type BookInfo struct {
	Id                 int       `form:"id"`
	BookName           string    `gorm:"type:varchar(200)" form:"book_name"`
	Author             string    `gorm:"type:varchar(50)" form:"author"`
	PublishingHouse    string    `gorm:"type:varchar(50)" form:"publishing_house"`
	Translator         string    `gorm:"type:varchar(50)" form:"translator"`
	PublishDate        time.Time `json:"publish_date" form:"publish_date"`
	Pages              int       `form:"pages"`
	ISBN               string    `gorm:"type:varchar(20)" form:"ISBN"`
	Price              float64   `form:"price"`
	BriefIntroduction  string    `gorm:"type:varchar(15000)" form:"brief_introduction"`
	AuthorIntroduction string    `gorm:"type:varchar(5000)" form:"author_introduction"`
	ImgUrl             string    `gorm:"type:varchar(200)" form:"img_url"`
	DelFlg             int       `json:"del_flg" form:"del_flg"`
}

type SendMsg struct {
	Id     int64  `form:"id"`
	UserId int64  `json:"user_id" form:"user_id"`
	SendId int64  `json:"send_id" form:"send_id"`
	Msg    string `json:"msg" form:"msg"`
	Status int    `json:"status"`
}

func (BookInfo) TableName() string {
	return "book_info"
}

func (User) TableName() string {
	return "user"
}

func (Book) TableName() string {
	return "book"
}

func (Category) TableName() string {
	return "category"
}
func (Librarian) TableName() string {
	return "librarian"
}
func (Record) TableName() string {
	return "record"
}
func (SendMsg) TableName() string {
	return "send_msg"
}
