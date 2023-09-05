package LM_V3

import (
	core2 "LibraryManagementV1/LM_V3/core"
	"LibraryManagementV1/LM_V3/global"
	"LibraryManagementV1/LM_V3/model"
	"LibraryManagementV1/LM_V3/router"
	"LibraryManagementV1/LM_V3/tools"
	"LibraryManagementV1/service/user_service"
)

func Start() {
	//读取配置
	core2.InitYaml(core2.ConfigFile)
	//白名单手机号
	user_service.ReadWhitePhones()
	//初始化全局日志
	core2.InitLogger()
	//连接数据库
	core2.InitGorm()
	//连接Redis
	core2.InitRedis()
	//
	model.InitStore()
	//
	tools.NewToken("")
	//初始化路由
	r := router.New()
	//
	addr := global.Config.System.Addr()
	global.Log.Println("项目运行在", addr)
	if err := r.Run(addr); err != nil {
		global.Log.Error(err)
	}

	//
	//model.MySql()
	//
	//defer func() {
	//	model.Close()
	//}()
	//
	//tools.NewToken("")
	//r := router.New()
	//_ = r.Run(":8083")
}
