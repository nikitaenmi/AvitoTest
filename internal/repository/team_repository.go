package repository

import (
	"context"
	"fmt"

	"github.com/nikitaenmi/AvitoTest/internal/database/models"
	"github.com/nikitaenmi/AvitoTest/internal/domain"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db       *gorm.DB
	userRepo *UserRepository
}

func NewTeamRepository(db *gorm.DB, userRepo *UserRepository) *TeamRepository {
	return &TeamRepository{
		db:       db,
		userRepo: userRepo,
	}
}

func (r *TeamRepository) Create(ctx context.Context, team domain.Team) error {
	teamModel := models.TeamFromDomain(team)
	return r.db.WithContext(ctx).Create(&teamModel).Error
}

func (r *TeamRepository) FindOne(ctx context.Context, filter domain.TeamFilter) (*domain.Team, error) {
	var teamModel models.Team
	q := r.db.WithContext(ctx)
	q = r.buildFilterByParams(q, filter)

	if err := q.First(&teamModel).Error; err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	members, err := r.userRepo.FindAll(ctx, domain.UserFilter{TeamName: &teamModel.TeamName})
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	team := models.TeamToDomain(teamModel, members)
	return &team, nil
}

func (r *TeamRepository) Exists(ctx context.Context, filter domain.TeamFilter) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&models.Team{})
	q = r.buildFilterByParams(q, filter)

	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}
	return count > 0, nil
}

func (r *TeamRepository) buildFilterByParams(q *gorm.DB, filter domain.TeamFilter) *gorm.DB {
	if filter.TeamName != nil {
		q = q.Where("team_name = ?", *filter.TeamName)
	}
	return q
}
