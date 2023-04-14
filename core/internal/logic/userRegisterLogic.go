package logic

import (
	"context"
	"errors"
	"go-cloud-disk/core/helper"
	"go-cloud-disk/core/modules"
	"log"

	"go-cloud-disk/core/internal/svc"
	"go-cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserRegister
//判断 code 是否和 Redis 中存的一致
//判断用户是否已经存在
//往数据库插入用户数据（生成 UUID，密码加密）
func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterReply, err error) {
	// 判断 code 是否一致
	code, err := l.svcCtx.RDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		return nil, errors.New("该邮箱的验证码为空，请重新发送验证码")
	}

	if code != req.Code {
		return nil, errors.New("验证码错误")
	}
	// 判断用户是否已经存在
	cnt, err := l.svcCtx.Engine.Where("name = ?", req.Name).Count(new(modules.UserBasic))
	if err != nil {
		return
	}
	if cnt > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 数据入库
	user := &modules.UserBasic{
		Identity: helper.UUID(),
		Name:     req.Name,
		Password: helper.Md5(req.Password),
		Email:    req.Email,
	}

	n, err := l.svcCtx.Engine.Insert(user)
	if err != nil {
		return
	}
	log.Println("insert user row:", n)

	return
}
