package middleware

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/policy"
)

const (
	RequestIDKey    = "request_id"
	RequestIDHeader = "X-Request-ID"
)

func RequestID(ctx context.Context, c *app.RequestContext) {
	requestID := strings.TrimSpace(string(c.Request.Header.Peek(RequestIDHeader)))
	if requestID == "" {
		generated, err := policy.RandomID()
		if err != nil {
			generated = strconv.FormatInt(time.Now().UnixNano(), 36)
		}
		requestID = generated
	}
	c.Set(RequestIDKey, requestID)
	c.Response.Header.Set(RequestIDHeader, requestID)
	c.Next(ctx)
}
