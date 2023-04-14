package logic

import (
	"context"
	"go-cloud-disk/core/helper"
	"go-cloud-disk/core/modules"

	"go-cloud-disk/core/internal/svc"
	"go-cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadPrepareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadPrepareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadPrepareLogic {
	return &FileUploadPrepareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadPrepareLogic) FileUploadPrepare(req *types.FileUploadPrepareRequest) (resp *types.FileUploadPrepareReply, err error) {
	rp := new(modules.RepositoryPool)
	has, err := l.svcCtx.Engine.Where("hash = ?", req.Md5).Get(rp)
	if err != nil {
		return
	}
	resp = new(types.FileUploadPrepareReply)
	if has {
		// 秒传成功
		resp.Identity = rp.Identity
	} else {
		// 获取该文件的UploadID、Key,用来进行文件的分片上传
		key, uploadId, err := helper.CosInitPart(req.Ext)
		if err != nil {
			return nil, err
		}
		resp.Key = key
		resp.UploadId = uploadId
	}

	return
}
