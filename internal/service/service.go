package service

import (
	"github.com/kiritoxkiriko/comical-tool/internal/repository"
	"github.com/kiritoxkiriko/comical-tool/pkg/jwt"
	"github.com/kiritoxkiriko/comical-tool/pkg/log"
	"github.com/kiritoxkiriko/comical-tool/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}
