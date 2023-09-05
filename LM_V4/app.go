package LM_V4

import (
	"LibraryManagementV1/GRPC_Server"
	core2 "LibraryManagementV1/LM_V4/core"
	"LibraryManagementV1/LM_V4/global"
	"LibraryManagementV1/LM_V4/model"
	"LibraryManagementV1/LM_V4/router"
	"LibraryManagementV1/LM_V4/tools"
)

func Start() {
	//读取配置
	core2.InitYaml(core2.ConfigFile)
	//白名单手机号
	//user_service.ReadWhitePhones()
	//初始化全局日志
	core2.InitLogger()
	//连接数据库
	core2.InitGorm()
	//连接Redis
	core2.InitRedis()
	//初始化session
	model.InitStore()
	//初始化token的密钥
	tools.NewToken("")
	// GRPC服务端
	// "LibraryManagementV1/GRPC_Server"
	go GRPC_Server.StartGRPCServer()
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
