package repository

import (
	"context"
	"fmt"

	"github.com/nikitaenmi/AvitoTest/internal/database/models"
	"github.com/nikitaenmi/AvitoTest/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	userModel := models.UserFromDomain(user)
	return r.db.WithContext(ctx).Create(&userModel).Error
}

func (r *UserRepository) FindOne(ctx context.Context, filter domain.UserFilter) (*domain.User, error) {
	var userModel models.User
	q := r.db.WithContext(ctx)
	q = r.buildFilterByParams(q, filter)

	if err := q.First(&userModel).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user := models.UserToDomain(userModel)
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	userModel := models.UserFromDomain(*user)
	return r.db.WithContext(ctx).Save(&userModel).Error
}

func (r *UserRepository) FindAll(ctx context.Context, filter domain.UserFilter) ([]domain.User, error) {
	var userModels []models.User
	q := r.db.WithContext(ctx)
	q = r.buildFilterByParams(q, filter)

	if err := q.Find(&userModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return models.UsersToDomain(userModels), nil
}

func (r *UserRepository) buildFilterByParams(q *gorm.DB, filter domain.UserFilter) *gorm.DB {
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}
	if filter.TeamName != nil {
		q = q.Where("team_name = ?", *filter.TeamName)
	}
	if filter.IsActive != nil {
		q = q.Where("is_active = ?", *filter.IsActive)
	}
	return q
}
