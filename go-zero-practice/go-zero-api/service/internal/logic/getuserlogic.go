package logic

import (
	"context"

	"go-zero-api/service/internal/svc"
	"go-zero-api/service/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUserLogic {
	return GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser() (*types.RegisterRequest, error) {
	// todo: add your logic here and delete this line

	return &types.RegisterRequest{}, nil
}
