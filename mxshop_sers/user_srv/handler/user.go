package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop_sers/user_srv/global"
	"mxshop_sers/user_srv/model"
	pb "mxshop_sers/user_srv/proto"
	"strings"
	"time"
)

type UserServer struct {
	pb.UnimplementedUserServer
}

// 将modle转换成pb的Response
func ModelToResponse(user model.User) pb.UserInfoResponse {
	userInfoRsp := pb.UserInfoResponse{
		Id:       int32(user.ID),
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	//不能赋nil
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

// 将数据分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// 返回用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *pb.PageInfo) (*pb.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &pb.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

// 通过手机号查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *pb.MobileRequest) (*pb.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("mobile = ?", req.Mobile).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// 通过ID查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *pb.IdRequest) (*pb.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id) //通常会对主键或者具有唯一约束的字段执行这样的操作
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// 新建一个用户，管理员，注册均可调用
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserInfo) (*pb.UserInfoResponse, error) {
	user := model.User{}
	result := global.DB.Where("mobile = ?", req.Mobile).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	//加密密码
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// 更新用户信息
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserInfo) (*pb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthday := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthday
	user.Gender = req.Gender

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &pb.Empty{}, nil
}

// 检查密码
func (s *UserServer) CheckPassWord(ctx context.Context, req *pb.PasswordCheckInfo) (*pb.CheckResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &pb.CheckResponse{
		Success: check,
	}, nil
}
