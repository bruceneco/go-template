package postgres

import (
	"context"
	"errors"
	"github.com/bruceneco/go-template/internal/adapters/db/models"
	"github.com/bruceneco/go-template/internal/domain/user"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

var _ user.Repository = (*UserRepository)(nil)

type UserRepository struct {
	db *Connection
}

func NewUserRepository(db *Connection) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (user.Entity, error) {
	model := models.User{
		Base: models.Base{ID: &id},
	}

	if err := u.db.WithContext(ctx).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.Entity{}, user.ErrUserNotFound
		}
		return user.Entity{}, err
	}

	return model.ToEntity(), nil
}

func (u *UserRepository) Create(ctx context.Context, entity *user.Entity) error {
	if entity == nil {
		return models.ErrNilEntity
	}

	model, err := new(models.User).FromEntity(entity)
	if err != nil {
		return err
	}

	if err = u.db.
		WithContext(ctx).
		Create(&model).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return user.ErrEmailAlreadyExists
		}

		return err
	}

	*entity = model.ToEntity()

	return nil
}

func (u *UserRepository) Update(ctx context.Context, entity *user.Entity) error {
	if entity == nil {
		return models.ErrNilEntity
	}
	model, err := new(models.User).FromEntity(entity)
	if err != nil {
		return err
	}
	tx := u.db.
		WithContext(ctx).
		Updates(&model)
	if err = tx.Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return user.ErrEmailAlreadyExists
		}

		return err
	}
	if tx.RowsAffected == 0 {
		return user.ErrUserNotFound
	}
	*entity = model.ToEntity()
	return nil
}

func (u *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	model := models.User{}
	model.ID = &id
	tx := u.db.
		WithContext(ctx).
		Delete(&model)
	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return user.ErrUserNotFound
	}
	return nil
}
