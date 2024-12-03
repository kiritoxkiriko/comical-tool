package cerror

import (
	"errors"
	"fmt"
)

const (
	// Base error code
	SuccessCode    = 0
	ServiceErrCode = iota + 10000
	ParamErrCode
	AuthorizationFailedErrCode

	// Custom error code
	URLExpiredErrCode
	URLNotExistsErrCode
	URLPasswordErrCode
)

const (
	SuccessMsg             = "Success"
	ServerErrMsg           = "Service is unable to start successfully"
	ParamErrMsg            = "Wrong Parameter has been given"
	AuthorizationFailedMsg = "Authorization failed"

	URLExpiredErrMsg   = "URL has expired"
	URLNotExistsErrMsg = "URL does not exist"
	URLPasswordErrMsg  = "URL password is incorrect"
)

var (
	Success                = Cerror{SuccessCode, SuccessMsg}
	ServiceErr             = Cerror{ServiceErrCode, ServerErrMsg}
	ParamErr               = Cerror{ParamErrCode, ParamErrMsg}
	AuthorizationFailedErr = Cerror{AuthorizationFailedErrCode, AuthorizationFailedMsg}

	URLExpiredErr   = Cerror{URLExpiredErrCode, URLExpiredErrMsg}
	URLNotExistsErr = Cerror{URLNotExistsErrCode, URLNotExistsErrMsg}
	URLPasswordErr  = Cerror{URLPasswordErrCode, URLPasswordErrMsg}
)

type Cerror struct {
	ErrCode int32
	ErrMsg  string
}

func (e Cerror) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int32, msg string) Cerror {
	return Cerror{code, msg}
}

func (e Cerror) WithMessage(msg string) Cerror {
	e.ErrMsg = msg
	return e
}

func (e Cerror) WithError(err error) Cerror {
	e.ErrMsg = fmt.Sprintf("%s: %s", e.ErrMsg, err.Error())
	return e
}

// ConvertErr convert error to Errno
func ConvertErr(err error) Cerror {
	Err := Cerror{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
