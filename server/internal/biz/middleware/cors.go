package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func CORS(ctx context.Context, c *app.RequestContext) {
	c.Response.Header.Set("Access-Control-Allow-Origin", "*")
	c.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	c.Next(ctx)
}

func Options(_ context.Context, c *app.RequestContext) {
	c.SetStatusCode(consts.StatusNoContent)
}
