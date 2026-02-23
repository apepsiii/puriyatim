package services

import (
	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
)

type ArtikelService struct {
	repo *repository.ArtikelRepository
}

func NewArtikelService(repo *repository.ArtikelRepository) *ArtikelService {
	return &ArtikelService{
		repo: repo,
	}
}

func (s *ArtikelService) GetAll() ([]*models.Artikel, error) {
	return s.repo.GetAll()
}

func (s *ArtikelService) GetByID(id string) (*models.Artikel, error) {
	return s.repo.GetByID(id)
}

func (s *ArtikelService) GetBySlug(slug string) (*models.Artikel, error) {
	return s.repo.GetBySlug(slug)
}

func (s *ArtikelService) GetPublished(limit int) ([]*models.Artikel, error) {
	return s.repo.GetPublished(limit)
}

func (s *ArtikelService) Create(artikel *models.Artikel) error {
	return s.repo.Create(artikel)
}

func (s *ArtikelService) Update(artikel *models.Artikel) error {
	return s.repo.Update(artikel)
}

func (s *ArtikelService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *ArtikelService) Count() (int, error) {
	return s.repo.Count()
}

func (s *ArtikelService) CountByStatus(status models.StatusPublikasi) (int, error) {
	return s.repo.CountByStatus(status)
}
