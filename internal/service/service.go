package service

import "github.com/nikitaenmi/AvitoTest/internal/domain"

type Service struct {
	userRepo domain.UserRepository
	teamRepo domain.TeamRepository
	prRepo   domain.PRRepository
}

func NewService(
	userRepo domain.UserRepository, teamRepo domain.TeamRepository, prRepo domain.PRRepository,
) *Service {
	return &Service{
		userRepo: userRepo,
		teamRepo: teamRepo,
		prRepo:   prRepo,
	}
}
