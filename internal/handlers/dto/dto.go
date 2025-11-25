package dto

import "github.com/nikitaenmi/AvitoTest/internal/domain"

type CreateTeamRequest struct {
	TeamName string        `json:"team_name"`
	Members  []UserRequest `json:"members"`
}

type UserRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

func (r CreateTeamRequest) ToDomain() domain.Team {
	members := make([]domain.User, len(r.Members))
	for i, member := range r.Members {
		members[i] = domain.User{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return domain.Team{
		TeamName: r.TeamName,
		Members:  members,
	}
}

func (r SetUserActiveRequest) ToUserFilter() domain.UserFilter {
	return domain.UserFilter{UserID: &r.UserID}
}

func (r CreatePRRequest) ToDomain() domain.PullRequest {
	return domain.PullRequest{
		PullRequestID:   r.PullRequestID,
		PullRequestName: r.PullRequestName,
		AuthorID:        r.AuthorID,
	}
}

func (r MergePRRequest) ToPRFilter() domain.PRFilter {
	return domain.PRFilter{PullRequestID: &r.PullRequestID}
}

func (r ReassignReviewerRequest) ToPRFilter() domain.PRFilter {
	return domain.PRFilter{PullRequestID: &r.PullRequestID}
}

func TeamFilterFromQuery(teamName string) domain.TeamFilter {
	return domain.TeamFilter{TeamName: &teamName}
}

func UserFilterFromQuery(userID string) domain.UserFilter {
	return domain.UserFilter{UserID: &userID}
}

func PRFilterFromQuery(prID string) domain.PRFilter {
	return domain.PRFilter{PullRequestID: &prID}
}
