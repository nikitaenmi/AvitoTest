package models

import (
	"time"

	"github.com/nikitaenmi/AvitoTest/internal/domain"
)

type Team struct {
	TeamName string `gorm:"primaryKey" json:"team_name"`
}

type User struct {
	UserID   string `gorm:"primaryKey" json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequest struct {
	PullRequestID     string     `gorm:"primaryKey" json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `gorm:"type:jsonb;serializer:json" json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

func UserToDomain(m User) domain.User {
	return domain.User{
		UserID:   m.UserID,
		Username: m.Username,
		TeamName: m.TeamName,
		IsActive: m.IsActive,
	}
}

func UserFromDomain(d domain.User) User {
	return User{
		UserID:   d.UserID,
		Username: d.Username,
		TeamName: d.TeamName,
		IsActive: d.IsActive,
	}
}

func UsersToDomain(models []User) []domain.User {
	domainUsers := make([]domain.User, len(models))
	for i, model := range models {
		domainUsers[i] = UserToDomain(model)
	}
	return domainUsers
}

func TeamToDomain(m Team, members []domain.User) domain.Team {
	return domain.Team{
		TeamName: m.TeamName,
		Members:  members,
	}
}

func TeamFromDomain(d domain.Team) Team {
	return Team{
		TeamName: d.TeamName,
	}
}

func PullRequestToDomain(m PullRequest) domain.PullRequest {
	return domain.PullRequest{
		PullRequestID:     m.PullRequestID,
		PullRequestName:   m.PullRequestName,
		AuthorID:          m.AuthorID,
		Status:            m.Status,
		AssignedReviewers: m.AssignedReviewers,
		CreatedAt:         m.CreatedAt,
		MergedAt:          m.MergedAt,
	}
}

func PullRequestFromDomain(d domain.PullRequest) PullRequest {
	return PullRequest{
		PullRequestID:     d.PullRequestID,
		PullRequestName:   d.PullRequestName,
		AuthorID:          d.AuthorID,
		Status:            d.Status,
		AssignedReviewers: d.AssignedReviewers,
		CreatedAt:         d.CreatedAt,
		MergedAt:          d.MergedAt,
	}
}

func PullRequestsToDomain(models []PullRequest) []domain.PullRequest {
	domainPRs := make([]domain.PullRequest, len(models))
	for i, model := range models {
		domainPRs[i] = PullRequestToDomain(model)
	}
	return domainPRs
}
