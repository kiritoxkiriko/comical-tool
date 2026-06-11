package handler

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kiritoxkiriko/comical-tool/server/internal/biz/middleware"
	"github.com/kiritoxkiriko/comical-tool/server/pkg/apperror"
)

type successEnvelope struct {
	Data any `json:"data"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code      apperror.Code `json:"code"`
	Message   string        `json:"message"`
	RequestID string        `json:"request_id,omitempty"`
}

func bindJSON(c *app.RequestContext, out any) bool {
	if err := json.Unmarshal(c.Request.Body(), out); err != nil {
		writeError(c, apperror.New(apperror.CodeBadRequest, "invalid json"))
		return false
	}
	return true
}

func parseTTL(value string, fallback time.Duration) (time.Duration, error) {
	if value == "" {
		return fallback, nil
	}
	return time.ParseDuration(value)
}

func writeResult(c *app.RequestContext, value any, err error) {
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, successEnvelope{Data: value})
}

func writeError(c *app.RequestContext, err error) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		appErr = apperror.New(apperror.CodeInternal, err.Error())
	}
	c.JSON(statusCode(appErr.Code), errorEnvelope{Error: errorBody{
		Code:      appErr.Code,
		Message:   appErr.Message,
		RequestID: requestID(c),
	}})
}

func requestID(c *app.RequestContext) string {
	value, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		return ""
	}
	id, ok := value.(string)
	if !ok {
		return ""
	}
	return id
}

func statusCode(code apperror.Code) int {
	switch code {
	case apperror.CodeBadRequest:
		return stdhttp.StatusBadRequest
	case apperror.CodeNotFound:
		return stdhttp.StatusNotFound
	case apperror.CodeUnauthorized:
		return stdhttp.StatusUnauthorized
	case apperror.CodeForbidden:
		return stdhttp.StatusForbidden
	case apperror.CodeConflict:
		return stdhttp.StatusConflict
	case apperror.CodeExpired, apperror.CodeRevoked:
		return stdhttp.StatusGone
	default:
		return stdhttp.StatusInternalServerError
	}
}
