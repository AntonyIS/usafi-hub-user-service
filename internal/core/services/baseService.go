package services

import "github.com/AntonyIS/usafi-hub-user-service/internal/core/ports"

type baseService struct {
	repo ports.BaseRepository
}

func NewBaseService(repo ports.BaseRepository) *baseService {
	service := baseService{
		repo: repo,
	}
	return &service
}

func (svc baseService) DropTables() error {
	return svc.repo.DropTables()
}
