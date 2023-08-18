package model

import (
	"fmt"
)

// GetUser 得到用户信息
func GetUser(name, password string) *User {
	//user := make([]*User, 0)
	user := &User{}
	//var user *User
	err := DB.Where("user_name = ? and password = ?", name, password).First(user).Error
	if err != nil {
		fmt.Printf("user:%+v\n", err.Error())
	}
	fmt.Printf("测试user数据：%+v\n", user)
	return user
}

func AddUser() {

}

//func UpdateUser(UserName string, password string, phone string, name string, userId int64) *User {
//	var user *User
//
//	DB.Find(&user, "name = ? and password = ?", name, password).Update("user_name", UserName)
//	DB.Find(&user, "name = ? and password = ?", name, password).Update("password", UserName)
//	return user
//}

func LookSelfMsg(userId int64) []*User {
	user := make([]*User, 0)
	sql := "select * from user where  id = ?"
	err := DB.Raw(sql, userId).Scan(&user).Error
	//err := DB.Table("user").Where("id = ?", userId).Find(&user).Error
	if err != nil {
		fmt.Printf("查看用户自己信息失败！err:%+v\n", err.Error())
	}
	return user
}

func GetUserBorrowStatus(userId int64, status int) *Record {
	//
	var userStatus *Record
	sql := "select * from record where status = ? and user_id = ?"
	err := DB.Raw(sql, status, userId).Scan(&userStatus).Error
	if err != nil {
		fmt.Printf("查询用户借阅状态失败！err:%+v\n", err.Error())
	}
	return userStatus
}

func UpdateUserMsg(userId int64, userName string, password string, phone string) []*User {
	user := make([]*User, 0)
	tx := DB.Begin()
	sql := "update user set password = ? ,phone = ?,user_name = ? where id = ?"
	err := tx.Exec(sql, password, phone, userName, userId).Error
	if err != nil {
		fmt.Printf("更新用户信息失败！err:%+v\n", err.Error())
		tx.Rollback()
	}
	tx.Commit()
	return user
}

func GetUsers() []*User {
	var user []*User
	sql := "select * from user"
	err := DB.Raw(sql).Scan(&user).Error
	if err != nil {
		fmt.Printf("管理员查询所有用户信息失败！err：%+v\n", err.Error())
	}
	return user
}
