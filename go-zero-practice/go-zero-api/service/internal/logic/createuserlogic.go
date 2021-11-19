package logic

import (
	"context"

	m "go-zero-api/service/internal/model"
	"go-zero-api/service/internal/svc"
	"go-zero-api/service/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateUserLogic {
	return CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req types.RegisterRequest) (*types.RegisterResponse, error) {
	// todo: add your logic here and delete this line

	

	userName := req.Username
	passWord := req.Password

	m.CreateUser(userName, passWord)

	return &types.RegisterResponse{
		Status:  "200",
		Message: "create user successful",
		Data:    userName,
	}, nil
}
