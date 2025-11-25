package service

import (
	"context"

	"github.com/nikitaenmi/AvitoTest/internal/domain"
)

func (s *Service) SetUserActive(ctx context.Context, filter domain.UserFilter, isActive bool) error {
	if filter.UserID == nil || *filter.UserID == "" {
		return domain.NewValidationError("user ID cannot be empty")
	}

	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		return domain.NewNotFoundError("user")
	}

	user.IsActive = isActive
	return s.userRepo.Update(ctx, user)
}

func (s *Service) GetUserReviewPRs(ctx context.Context, filter domain.UserFilter) ([]domain.PullRequest, error) {
	if filter.UserID == nil || *filter.UserID == "" {
		return nil, domain.NewValidationError("user ID cannot be empty")
	}

	if _, err := s.userRepo.FindOne(ctx, filter); err != nil {
		return nil, domain.NewNotFoundError("user")
	}

	return s.prRepo.FindByReviewer(ctx, *filter.UserID)
}

func (s *Service) GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserIDs ...string) ([]domain.User, error) {
	isActive := true
	users, err := s.userRepo.FindAll(ctx, domain.UserFilter{
		TeamName: &teamName,
		IsActive: &isActive,
	})
	if err != nil {
		return nil, err
	}

	var filtered []domain.User
	for _, user := range users {
		exclude := false
		for _, excludeID := range excludeUserIDs {
			if user.UserID == excludeID {
				exclude = true
				break
			}
		}
		if !exclude {
			filtered = append(filtered, user)
		}
	}

	return filtered, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, domain.NewValidationError("user ID cannot be empty")
	}
	return s.userRepo.FindOne(ctx, domain.UserFilter{UserID: &userID})
}
