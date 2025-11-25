package service

import (
	"context"

	"github.com/nikitaenmi/AvitoTest/internal/domain"
)

func (s *Service) CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error) {
	exists, err := s.teamRepo.Exists(ctx, domain.TeamFilter{TeamName: &team.TeamName})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.NewTeamExistsError()
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	for i := range team.Members {
		team.Members[i].TeamName = team.TeamName
		if err := s.userRepo.Create(ctx, team.Members[i]); err != nil {
			return nil, err
		}
	}

	return s.teamRepo.FindOne(ctx, domain.TeamFilter{TeamName: &team.TeamName})
}

func (s *Service) GetTeam(ctx context.Context, filter domain.TeamFilter) (*domain.Team, error) {
	if filter.TeamName == nil || *filter.TeamName == "" {
		return nil, domain.NewValidationError("team name cannot be empty")
	}
	team, err := s.teamRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, domain.NewNotFoundError("team")
	}
	return team, nil
}
