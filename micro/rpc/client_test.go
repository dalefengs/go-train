package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setFuncField(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr error
	}{
		{
			name:    "nil",
			service: nil,
			wantErr: errors.New("不支持 nil"),
		},
		{
			name:    "not pointer",
			service: UserService{},
			wantErr: errors.New("只支持指向结构体的一级指针"),
		},
		{
			name:    "pointer",
			service: &UserService{},
			wantErr: errors.New("只支持指向结构体的一级指针"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := Client{}
			err := setFuncField(tt.service, ctrl)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			resp, err := tt.service.(*UserService).GetById(context.Background(), &GetByIdReq{Id: 123})
			assert.Equal(t, tt.wantErr, err)
			t.Log(resp)
		})
	}
}
