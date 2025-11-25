package domain

import (
	"context"
	"time"
)

type Team struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type UserRepository interface {
	Create(ctx context.Context, user User) error
	FindOne(ctx context.Context, filter UserFilter) (*User, error)
	Update(ctx context.Context, user *User) error
	FindAll(ctx context.Context, filter UserFilter) ([]User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team Team) error
	FindOne(ctx context.Context, filter TeamFilter) (*Team, error)
	Exists(ctx context.Context, filter TeamFilter) (bool, error)
}

type PRRepository interface {
	Create(ctx context.Context, pr PullRequest) error
	FindOne(ctx context.Context, filter PRFilter) (*PullRequest, error)
	Update(ctx context.Context, pr *PullRequest) error
	FindAll(ctx context.Context, filter PRFilter) ([]PullRequest, error)
	Exists(ctx context.Context, filter PRFilter) (bool, error)
	FindByReviewer(ctx context.Context, userID string) ([]PullRequest, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, team Team) (*Team, error)
	GetTeam(ctx context.Context, filter TeamFilter) (*Team, error)
}

type UserService interface {
	SetUserActive(ctx context.Context, filter UserFilter, isActive bool) error
	GetUserReviewPRs(ctx context.Context, filter UserFilter) ([]PullRequest, error)
	GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserIDs ...string) ([]User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
}

type PRService interface {
	CreatePR(ctx context.Context, pr PullRequest) (*PullRequest, error)
	MergePR(ctx context.Context, filter PRFilter) error
	ReassignReviewer(ctx context.Context, filter PRFilter, oldReviewerID string) (string, error)
	GetPR(ctx context.Context, filter PRFilter) (*PullRequest, error)
	HealthCheck(ctx context.Context) error
}

type Service interface {
	TeamService
	UserService
	PRService
}

type UserFilter struct {
	UserID   *string
	TeamName *string
	IsActive *bool
}

type TeamFilter struct {
	TeamName *string
}

type PRFilter struct {
	PullRequestID *string
	AuthorID      *string
	Status        *string
}

const (
	PRStatusOpen   = "OPEN"
	PRStatusMerged = "MERGED"
)
