package services

import (
	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
)

type AnakAsuhService struct {
	repo *repository.AnakAsuhRepository
}

func NewAnakAsuhService(repo *repository.AnakAsuhRepository) *AnakAsuhService {
	return &AnakAsuhService{
		repo: repo,
	}
}

func (s *AnakAsuhService) GetAll() ([]*models.AnakAsuh, error) {
	return s.repo.GetAll()
}

func (s *AnakAsuhService) GetByID(id string) (*models.AnakAsuh, error) {
	return s.repo.GetByID(id)
}

func (s *AnakAsuhService) Create(anak *models.AnakAsuh) error {
	return s.repo.Create(anak)
}

func (s *AnakAsuhService) Update(anak *models.AnakAsuh) error {
	return s.repo.Update(anak)
}

func (s *AnakAsuhService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *AnakAsuhService) Count() (int, error) {
	return s.repo.Count()
}

func (s *AnakAsuhService) CountAll() (int, error) {
	return s.repo.CountAll()
}

func (s *AnakAsuhService) GetByRT(rt string) ([]*models.AnakAsuh, error) {
	return s.repo.GetByRT(rt)
}

func (s *AnakAsuhService) GetByRTRW(rt, rw string) ([]*models.AnakAsuh, error) {
	return s.repo.GetByRTRW(rt, rw)
}
