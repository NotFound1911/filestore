package grpc

import (
	"context"
	userv1 "github.com/NotFound1911/filestore/app/account/api/proto/gen/user/v1"
	"github.com/NotFound1911/filestore/app/account/domain"
	"github.com/NotFound1911/filestore/app/account/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountServiceServer struct {
	userv1.UnimplementedUserServiceServer
	svc service.UserService
}

func NewAccountServiceServer(svc service.UserService) *AccountServiceServer {
	return &AccountServiceServer{svc: svc}
}
func (a *AccountServiceServer) Signup(ctx context.Context, req *userv1.SignupReq) (*userv1.SignupResp, error) {
	err := a.svc.Signup(ctx, a.toDomain(req))
	return &userv1.SignupResp{}, err
}
func (a *AccountServiceServer) toDomain(r *userv1.SignupReq) domain.User {
	creatTimestamp := r.User.CreateAt.AsTime()
	updateTimestamp := r.User.UpdateAt.AsTime()
	return domain.User{
		Name:     r.User.Name,
		Email:    r.User.Email,
		Password: r.User.Password,
		Phone:    r.User.Phone,
		Status:   r.User.Status,
		CreateAt: &creatTimestamp,
		UpdateAt: &updateTimestamp,
	}
}

func (a *AccountServiceServer) Login(ctx context.Context, req *userv1.LoginReq) (*userv1.LoginResp, error) {
	user, err := a.svc.Login(ctx, req.Email, req.Password)
	return &userv1.LoginResp{Id: user.Id}, err
}
func (a *AccountServiceServer) Profile(ctx context.Context, req *userv1.ProfileReq) (*userv1.ProfileResp, error) {
	user, err := a.svc.FindById(ctx, req.Id)
	return &userv1.ProfileResp{
		User: &userv1.User{
			Id:       user.Id,
			Email:    user.Email,
			Name:     user.Name,
			Phone:    user.Phone,
			CreateAt: timestamppb.New(*user.CreateAt),
			UpdateAt: timestamppb.New(*user.UpdateAt),
		},
	}, err
}
