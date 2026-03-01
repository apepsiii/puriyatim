package services

import (
	"errors"
	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
	"strings"
)

type RekeningDonasiService struct {
	repo *repository.RekeningDonasiRepository
}

func NewRekeningDonasiService(repo *repository.RekeningDonasiRepository) *RekeningDonasiService {
	return &RekeningDonasiService{repo: repo}
}

func (s *RekeningDonasiService) GetAll() ([]*models.RekeningDonasi, error) {
	if s.repo == nil {
		return []*models.RekeningDonasi{}, nil
	}
	return s.repo.GetAll()
}

func (s *RekeningDonasiService) Create(item *models.RekeningDonasi) error {
	if s.repo == nil {
		return errors.New("layanan rekening tidak tersedia")
	}
	item.NamaBank = strings.TrimSpace(item.NamaBank)
	item.LogoBank = strings.TrimSpace(item.LogoBank)
	item.NomorRekening = strings.TrimSpace(item.NomorRekening)
	item.AtasNama = strings.TrimSpace(item.AtasNama)

	if item.NamaBank == "" || item.NomorRekening == "" || item.AtasNama == "" {
		return errors.New("nama bank, nomor rekening, dan atas nama wajib diisi")
	}
	item.Aktif = true
	return s.repo.Create(item)
}

func (s *RekeningDonasiService) Update(item *models.RekeningDonasi) error {
	if s.repo == nil {
		return errors.New("layanan rekening tidak tersedia")
	}
	return s.repo.Update(item)
}

func (s *RekeningDonasiService) Delete(id int) error {
	if s.repo == nil {
		return errors.New("layanan rekening tidak tersedia")
	}
	if id <= 0 {
		return errors.New("id tidak valid")
	}
	return s.repo.Delete(id)
}
