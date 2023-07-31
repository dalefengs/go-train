package rpc

import "context"

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
}

func (u UserService) Name() string {
	return "user-service"
}
