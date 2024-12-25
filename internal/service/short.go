package service

import (
    "context"
	"github.com/kiritoxkiriko/comical-tool/internal/model"
	"github.com/kiritoxkiriko/comical-tool/internal/repository"
)

type ShortService interface {
	GetShort(ctx context.Context, id int64) (*model.Short, error)
}
func NewShortService(
    service *Service,
    shortRepository repository.ShortRepository,
) ShortService {
	return &shortService{
		Service:        service,
		shortRepository: shortRepository,
	}
}

type shortService struct {
	*Service
	shortRepository repository.ShortRepository
}

func (s *shortService) GetShort(ctx context.Context, id int64) (*model.Short, error) {
	return s.shortRepository.GetShort(ctx, id)
}
