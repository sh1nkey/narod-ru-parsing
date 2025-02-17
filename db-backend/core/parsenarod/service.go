package parsenarod

import (
	"context"
	"data-sender/core"

	"github.com/rs/zerolog/log"
)

type Repository interface {
	Create(ctx context.Context, url *UrlReqDTO, tx ...core.UpdateOptions) error
	GetAll(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]Url, error)
	MarkAsEmpty(ctx context.Context, url string, options ...core.UpdateOptions) error
	SetDescription(ctx context.Context, url string, description string, options ...core.UpdateOptions) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	log.Info().Msg("creating user service...")

	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, url *UrlReqDTO) error {
	err := s.repo.Create(ctx, url)
	if err != nil {
		return err
	}
	return err
}

func (s *service) GetAll(ctx context.Context, limit, offset int, tx ...core.QueryOptions) ([]Url, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *service) MarkAsEmpty(ctx context.Context, url string, options ...core.UpdateOptions) error {
	return s.repo.MarkAsEmpty(ctx, url, options...)
}

func (s *service) SetDescription(ctx context.Context, url string, description string, options ...core.UpdateOptions) error {
	return s.repo.SetDescription(ctx, url, description, options...)
}
