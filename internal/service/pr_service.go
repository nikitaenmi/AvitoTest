package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/nikitaenmi/AvitoTest/internal/domain"
)

func (s *Service) CreatePR(ctx context.Context, pr domain.PullRequest) (*domain.PullRequest, error) {
	if pr.PullRequestID == "" {
		return nil, domain.NewValidationError("pull request ID cannot be empty")
	}

	exists, err := s.prRepo.Exists(ctx, domain.PRFilter{PullRequestID: &pr.PullRequestID})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.NewPRExistsError()
	}

	author, err := s.userRepo.FindOne(ctx, domain.UserFilter{UserID: &pr.AuthorID})
	if err != nil {
		return nil, domain.NewNotFoundError("author")
	}

	teamMembers, err := s.GetActiveTeamMembers(ctx, author.TeamName, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewers := []string{}
	for i := 0; i < len(teamMembers) && i < 2; i++ {
		reviewers = append(reviewers, teamMembers[i].UserID)
	}

	now := time.Now()
	pr.Status = domain.PRStatusOpen
	pr.AssignedReviewers = reviewers
	pr.CreatedAt = &now
	pr.MergedAt = nil

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (s *Service) MergePR(ctx context.Context, filter domain.PRFilter) error {
	if filter.PullRequestID == nil || *filter.PullRequestID == "" {
		return domain.NewValidationError("pull request ID cannot be empty")
	}

	pr, err := s.prRepo.FindOne(ctx, filter)
	if err != nil {
		return domain.NewNotFoundError("pull request")
	}

	if pr.Status == domain.PRStatusMerged {
		return nil
	}

	pr.Status = domain.PRStatusMerged
	now := time.Now()
	pr.MergedAt = &now

	return s.prRepo.Update(ctx, pr)
}

func (s *Service) ReassignReviewer(ctx context.Context, filter domain.PRFilter, oldReviewerID string) (string, error) {
	if filter.PullRequestID == nil || *filter.PullRequestID == "" {
		return "", domain.NewValidationError("pull request ID cannot be empty")
	}
	if oldReviewerID == "" {
		return "", domain.NewValidationError("old reviewer ID cannot be empty")
	}

	pr, err := s.prRepo.FindOne(ctx, filter)
	if err != nil {
		return "", domain.NewNotFoundError("pull request")
	}

	if pr.Status == domain.PRStatusMerged {
		return "", domain.NewPRMergedError()
	}

	found := false
	reviewerIndex := -1
	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldReviewerID {
			found = true
			reviewerIndex = i
			break
		}
	}

	if !found {
		return "", domain.NewNotAssignedError()
	}

	oldReviewer, err := s.userRepo.FindOne(ctx, domain.UserFilter{UserID: &oldReviewerID})
	if err != nil {
		return "", domain.NewNotFoundError("old reviewer")
	}

	excludedUsers := []string{oldReviewerID, pr.AuthorID}
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer != oldReviewerID {
			excludedUsers = append(excludedUsers, reviewer)
		}
	}

	availableReviewers, err := s.GetActiveTeamMembers(ctx, oldReviewer.TeamName, excludedUsers...)
	if err != nil {
		return "", err
	}

	if len(availableReviewers) == 0 {
		return "", domain.NewNoCandidateError()
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	newReviewer := availableReviewers[r.Intn(len(availableReviewers))]
	newReviewerID := newReviewer.UserID

	pr.AssignedReviewers[reviewerIndex] = newReviewerID

	if err := s.prRepo.Update(ctx, pr); err != nil {
		return "", err
	}

	return newReviewerID, nil
}

func (s *Service) GetPR(ctx context.Context, filter domain.PRFilter) (*domain.PullRequest, error) {
	if filter.PullRequestID == nil || *filter.PullRequestID == "" {
		return nil, domain.NewValidationError("pull request ID cannot be empty")
	}
	return s.prRepo.FindOne(ctx, filter)
}

func (s *Service) HealthCheck(ctx context.Context) error {
	testTeam := "health_check_test_team_12345"
	_, err := s.teamRepo.Exists(ctx, domain.TeamFilter{TeamName: &testTeam})
	return err
}
