package logic

import (
	"context"
	"errors"
	"go-cloud-disk/core/internal/svc"
	"go-cloud-disk/core/internal/types"
	"go-cloud-disk/core/modules"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileNameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileNameUpdateLogic {
	return &UserFileNameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileNameUpdateLogic) UserFileNameUpdate(req *types.UserFileNameUpdatedRequest, userIdentity string) (resp *types.UserFieNameUpdateReply, err error) {

		// 判断当前名称在该层级下是否存在
		cnt, err := l.svcCtx.Engine.Where("name = ? AND parent_id = (SELECT parent_id FROM user_repository ur WHERE ur.identity = ?)",
			req.Name, req.Identity).Count(new(modules.UserRepository))
		if err != nil {
			return
		}
		if cnt > 0 {
			return nil, errors.New("该名称已存在")
		}
		// 文件名称修改
		data := &modules.UserRepository{Name: req.Name}
		_, err = l.svcCtx.Engine.Where("identity = ? AND user_identity = ? ", req.Identity, userIdentity).Update(data)
		return

}
