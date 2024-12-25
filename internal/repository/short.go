package repository

import (
    "context"
	"github.com/kiritoxkiriko/comical-tool/internal/model"
)

type ShortRepository interface {
	GetShort(ctx context.Context, id int64) (*model.Short, error)
}

func NewShortRepository(
	repository *Repository,
) ShortRepository {
	return &shortRepository{
		Repository: repository,
	}
}

type shortRepository struct {
	*Repository
}

func (r *shortRepository) GetShort(ctx context.Context, id int64) (*model.Short, error) {
	var short model.Short

	return &short, nil
}
