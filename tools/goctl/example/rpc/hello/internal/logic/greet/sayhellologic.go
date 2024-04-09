package greetlogic

import (
	"context"

	"github.com/jialequ/linux-sdk/tools/goctl/example/rpc/hello/internal/svc"
	"github.com/jialequ/linux-sdk/tools/goctl/example/rpc/hello/pb/hello"

	"github.com/jialequ/linux-sdk/core/logx"
)

type SayHelloLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSayHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHelloLogic {
	return &SayHelloLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SayHelloLogic) SayHello(in *hello.HelloReq) (*hello.HelloResp, error) {
	// : add your logic here and delete this line

	return &hello.HelloResp{}, nil
}
