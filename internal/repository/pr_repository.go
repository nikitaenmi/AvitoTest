package repository

import (
	"context"
	"fmt"

	"github.com/nikitaenmi/AvitoTest/internal/database/models"
	"github.com/nikitaenmi/AvitoTest/internal/domain"
	"gorm.io/gorm"
)

type PRRepository struct {
	db *gorm.DB
}

func NewPRRepository(db *gorm.DB) *PRRepository {
	return &PRRepository{db: db}
}

func (r *PRRepository) Create(ctx context.Context, pr domain.PullRequest) error {
	prModel := models.PullRequestFromDomain(pr)
	return r.db.WithContext(ctx).Create(&prModel).Error
}

func (r *PRRepository) FindOne(ctx context.Context, filter domain.PRFilter) (*domain.PullRequest, error) {
	var prModel models.PullRequest
	q := r.db.WithContext(ctx)
	q = r.buildFilterByParams(q, filter)

	if err := q.First(&prModel).Error; err != nil {
		return nil, fmt.Errorf("pull request not found: %w", err)
	}

	pr := models.PullRequestToDomain(prModel)
	return &pr, nil
}

func (r *PRRepository) Update(ctx context.Context, pr *domain.PullRequest) error {
	prModel := models.PullRequestFromDomain(*pr)
	return r.db.WithContext(ctx).Save(&prModel).Error
}

func (r *PRRepository) FindAll(ctx context.Context, filter domain.PRFilter) ([]domain.PullRequest, error) {
	var prModels []models.PullRequest
	q := r.db.WithContext(ctx)
	q = r.buildFilterByParams(q, filter)

	if err := q.Find(&prModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find pull requests: %w", err)
	}

	return models.PullRequestsToDomain(prModels), nil
}

func (r *PRRepository) Exists(ctx context.Context, filter domain.PRFilter) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&models.PullRequest{})
	q = r.buildFilterByParams(q, filter)

	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check PR existence: %w", err)
	}
	return count > 0, nil
}

func (r *PRRepository) FindByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	var prModels []models.PullRequest
	err := r.db.WithContext(ctx).
		Where("assigned_reviewers::jsonb @> ?", fmt.Sprintf(`["%s"]`, userID)).
		Find(&prModels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find PRs by reviewer: %w", err)
	}

	return models.PullRequestsToDomain(prModels), nil
}

func (r *PRRepository) buildFilterByParams(q *gorm.DB, filter domain.PRFilter) *gorm.DB {
	if filter.PullRequestID != nil {
		q = q.Where("pull_request_id = ?", *filter.PullRequestID)
	}
	if filter.AuthorID != nil {
		q = q.Where("author_id = ?", *filter.AuthorID)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	return q
}
