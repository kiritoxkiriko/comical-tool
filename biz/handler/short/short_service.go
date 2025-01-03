// Code generated by hertz generator.

package short

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	common "github.com/kiritoxkiriko/comical-tool/biz/model/common"
	short "github.com/kiritoxkiriko/comical-tool/biz/model/short"
	"github.com/kiritoxkiriko/comical-tool/pkg/cerror"
	"github.com/kiritoxkiriko/comical-tool/pkg/util"
)

// Short .
// @router /s/short [POST]
func Short(ctx context.Context, c *app.RequestContext) {
	var err error
	var req short.ShortReq
	err = c.BindAndValidate(&req)
	if err != nil {
		util.GenBaseResp(c, nil, cerror.ParamErr)
		return
	}

	resp := new(common.BaseResp)

	c.JSON(consts.StatusOK, resp)
}

// Revoke .
// @router /s/revoke [POST]
func Revoke(ctx context.Context, c *app.RequestContext) {
	var err error
	var req short.RevokeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		util.GenBaseResp(c, nil, cerror.ParamErr)
		return
	}

	resp := new(common.BaseResp)

	c.JSON(consts.StatusOK, resp)
}

// GetShort .
// @router /s/:code [GET]
func GetShort(ctx context.Context, c *app.RequestContext) {
	var err error
	var req short.GetShortReq
	err = c.BindAndValidate(&req)
	if err != nil {
		util.GenBaseResp(c, nil, cerror.ParamErr)
		return
	}

	resp := new(common.BaseResp)

	c.JSON(consts.StatusOK, resp)
}
