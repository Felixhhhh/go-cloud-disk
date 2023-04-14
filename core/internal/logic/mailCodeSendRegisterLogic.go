package logic

import (
	"context"
	"errors"
	"go-cloud-disk/core/define"
	"go-cloud-disk/core/helper"
	"go-cloud-disk/core/modules"
	"time"

	"go-cloud-disk/core/internal/svc"
	"go-cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MailCodeSendRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMailCodeSendRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MailCodeSendRegisterLogic {
	return &MailCodeSendRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MailCodeSendRegisterLogic) MailCodeSendRegister(req *types.MailCodeSendRequest) (resp *types.MailCodeSendReply, err error) {
	// 1. 从数据库中查询当前用户
	cnt, err := l.svcCtx.Engine.Where("email = ?", req.Email).Count(new(modules.UserBasic))

	if err != nil {
		return
	}
	// 用户不存在
	if cnt > 0 {
		err = errors.New("该邮箱已注册")
		return
	}

	// 生成验证码
	code := helper.RandCode()
	// 存储验证码到 Redis（设置过期时间）
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))
	// 发送验证码
	err = helper.MailSendCode(req.Email, code)
	if err != nil {
		return nil, err
	}
	return
}
