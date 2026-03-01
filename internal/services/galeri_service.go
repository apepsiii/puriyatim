package services

import (
	"fmt"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
)

type GaleriService struct {
	repo *repository.GaleriRepository
}

func NewGaleriService(repo *repository.GaleriRepository) *GaleriService {
	return &GaleriService{repo: repo}
}

func (s *GaleriService) Create(item *models.GaleriFoto) error {
	if item.Judul == "" {
		return fmt.Errorf("judul foto wajib diisi")
	}
	if item.GambarAsliURL == "" || item.GambarOverlayURL == "" {
		return fmt.Errorf("file gambar tidak valid")
	}
	return s.repo.Create(item)
}

func (s *GaleriService) ListAll() ([]*models.GaleriFoto, error) {
	return s.repo.ListAll()
}

func (s *GaleriService) GetByID(id string) (*models.GaleriFoto, error) {
	return s.repo.GetByID(id)
}

func (s *GaleriService) Update(item *models.GaleriFoto) error {
	if item.Judul == "" {
		return fmt.Errorf("judul foto wajib diisi")
	}
	if item.GambarAsliURL == "" || item.GambarOverlayURL == "" {
		return fmt.Errorf("file gambar tidak valid")
	}
	return s.repo.Update(item)
}

func (s *GaleriService) Delete(id string) error {
	return s.repo.Delete(id)
}
