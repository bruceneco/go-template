package user

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (Entity, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) CreateUser(ctx context.Context, entity *Entity) error {
	if entity == nil {
		return ErrNilEntity
	}

	err := s.repo.Create(ctx, entity)
	if err != nil && !errors.Is(err, ErrEmailAlreadyExists) {
		log.Error().Err(err).Msg("failed to create user")
	}

	return err
}

func (s *Service) UpdateUser(ctx context.Context, entity *Entity) error {
	if entity == nil || entity.ID == nil {
		return ErrInvalidID
	}

	err := s.repo.Update(ctx, entity)
	if err != nil && !errors.Is(err, ErrUserNotFound) && !errors.Is(err, ErrEmailAlreadyExists) {
		log.Error().Err(err).Msg("failed to update user")
	}

	return err
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if !errors.Is(err, ErrUserNotFound) {
		log.Error().Err(err).Msg("failed to delete user")
	}
	return err
}
