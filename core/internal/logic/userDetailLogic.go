package logic

import (
	"context"
	"errors"
	"go-cloud-disk/core/modules"

	"go-cloud-disk/core/internal/svc"
	"go-cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(req *types.UserDetailRequest) (resp *types.UserDetailReply, err error) {
	// 1. 从数据库中查询当前用户
	user := new(modules.UserBasic)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.Identity).Get(user)
	if err != nil {
		return
	}
	// 用户不存在
	if !has {
		err = errors.New("用户不存在")
		return
	}
	return &types.UserDetailReply{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
