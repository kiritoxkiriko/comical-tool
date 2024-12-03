package util

import (
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/kiritoxkiriko/comical-tool/pkg/cerror"
)

type BaseResp[T any] struct {
	StatusCode int32  `json:"code"`
	StatusMsg  string `json:"message"`
	Data       T      `json:"data,omitempty"`
}

func GenBaseResp(c *app.RequestContext, data any, err error) {
	resp := BuildBaseRespWithData[any](data, err)
	c.JSON(consts.StatusOK, resp)
}

// BuildBaseRespWithData convert data, error and build BaseResp
func BuildBaseRespWithData[T any](data any, err error) *BaseResp[T] {
	if err == nil {
		resp := baseResp[T](cerror.Success)
		resp.Data = data.(T)
	}

	e := cerror.Cerror{}
	if errors.As(err, &e) {
		return baseResp[T](e)
	}

	s := cerror.ServiceErr.WithMessage(err.Error())
	return baseResp[T](s)
}

// BuildBaseResp convert error and build BaseResp
func BuildBaseResp(err error) *BaseResp[string] {
	return BuildBaseRespWithData[string]("", err)
}

// baseResp build BaseResp from error
func baseResp[T any](err cerror.Cerror) *BaseResp[T] {
	return &BaseResp[T]{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
}
