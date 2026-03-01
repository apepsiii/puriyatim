package services

import (
	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
)

type PengaturanService struct {
	repo *repository.PengaturanRepository
}

func NewPengaturanService(repo *repository.PengaturanRepository) *PengaturanService {
	return &PengaturanService{repo: repo}
}

func (s *PengaturanService) Get() (*models.PengaturanWeb, error) {
	return s.repo.Get()
}

func (s *PengaturanService) Save(p *models.PengaturanWeb) error {
	return s.repo.Save(p)
}
