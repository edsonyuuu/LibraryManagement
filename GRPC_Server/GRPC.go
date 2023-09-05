package GRPC_Server

import (
	pb "LibraryManagementV1/GRPC_Server/user_proto"
	"LibraryManagementV1/LM_V4/global"
	"LibraryManagementV1/LM_V4/logic"
	"LibraryManagementV1/LM_V4/model"
	"LibraryManagementV1/LM_V4/tools"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type userServiceServer struct{}

func (s *userServiceServer) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	phone := req.Phone
	// 其他逻辑，与原始函数相同
	// ...
	ans := logic.IsPhoneNumber(phone)
	if ans == false {

		return &pb.UserLoginResponse{
			Code:    tools.Failed,
			Message: "输入手机号为非法手机号，请重新输入~",
		}, nil

	}

	code, err := global.RedisConn.Get(ctx, phone).Result()
	if err != nil {
		fmt.Println("Failed to store verification code in Redis:", err.Error())
		return nil, err
	}

	// 获取用户输入的验证码
	userCode := req.Code
	// 验证验证码
	fmt.Println("Redis code:", code)
	fmt.Println("form code:", userCode)
	if userCode == code {
		var user model.User
		// 从 req 中提取用户名和密码等信息
		user.UserName = req.UserName
		user.Password = req.Password

		DbUser := model.GetUser(user.UserName, user.Password)
		fmt.Printf("打印user登录信息:%+v\n", DbUser)

		if DbUser.Id <= 0 {
			return &pb.UserLoginResponse{
				Code:    tools.NotFound,
				Message: "没有找到用户信息",
			}, nil
		}

		// 下发 Token
		a, r, err := tools.Token.GetToken(DbUser.Id, DbUser.UserName)
		fmt.Printf("atoken:%s\n", a)
		fmt.Printf("rtoken:%s\n", r)
		if err != nil {
			return &pb.UserLoginResponse{
				Code:    tools.UserError,
				Message: "Token生效失败！错误信息：" + err.Error(),
			}, err
		}
		// 生成响应
		return &pb.UserLoginResponse{
			Code:         tools.OK,
			Message:      "登录成功,发送token,正在跳转~",
			AccessToken:  a,
			RefreshToken: r,
		}, nil
	}

	return &pb.UserLoginResponse{
		Code:    tools.Failed,
		Message: "验证码错误！",
	}, nil
}

// StartGRPCServer
func StartGRPCServer() {
	listen, err := net.Listen("tcp", ":8083")
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		return
	}
	// 创建grpc服务
	s := grpc.NewServer()
	// 在grpc服务里注册服务
	pb.RegisterUserServiceServer(s, &userServiceServer{})
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		fmt.Printf("Failed to serve: %v", err)
		return
	}
	fmt.Println("grpc server running :9090")
}
