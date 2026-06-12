package main

import (
	"syscall/js"
	"time"

	"github.com/kiritoxkiriko/comical-tool/server/pkg/policy"
)

func main() {
	js.Global().Set("comicalPolicy", map[string]any{
		"expiryUnix":         js.FuncOf(expiryUnix),
		"expiredUnix":        js.FuncOf(expiredUnix),
		"randomSlug":         js.FuncOf(randomSlug),
		"validateSlug":       js.FuncOf(validateSlug),
		"visitLimitExceeded": js.FuncOf(visitLimitExceeded),
	})
	select {}
}

func expiryUnix(_ js.Value, args []js.Value) any {
	ttl := argString(args, 0)
	fallbackSeconds := argInt(args, 1)
	duration, err := policy.ParseTTLDuration(ttl, time.Duration(fallbackSeconds)*time.Second)
	if err != nil {
		return -1
	}
	expiresAt := policy.ExpiryFromDuration(duration)
	if expiresAt == nil {
		return 0
	}
	return expiresAt.Unix()
}

func expiredUnix(_ js.Value, args []js.Value) any {
	timestamp := int64(argInt(args, 0))
	if timestamp <= 0 {
		return false
	}
	expiresAt := time.Unix(timestamp, 0).UTC()
	return policy.IsExpired(&expiresAt)
}

func randomSlug(_ js.Value, _ []js.Value) any {
	slug, err := policy.RandomSlug()
	if err != nil {
		return ""
	}
	return slug
}

func validateSlug(_ js.Value, args []js.Value) any {
	return policy.ValidateSlug(argString(args, 0))
}

func visitLimitExceeded(_ js.Value, args []js.Value) any {
	return policy.VisitLimitExceeded(argInt(args, 0), argInt(args, 1))
}

func argString(args []js.Value, index int) string {
	if len(args) <= index {
		return ""
	}
	return args[index].String()
}

func argInt(args []js.Value, index int) int {
	if len(args) <= index {
		return 0
	}
	return args[index].Int()
}
