package services

import (
	"fmt"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
)

type JumatBerkahService struct {
	repo         *repository.JumatBerkahRepository
	anakAsuhRepo *repository.AnakAsuhRepository
}

func NewJumatBerkahService(repo *repository.JumatBerkahRepository, anakAsuhRepo *repository.AnakAsuhRepository) *JumatBerkahService {
	s := &JumatBerkahService{
		repo:         repo,
		anakAsuhRepo: anakAsuhRepo,
	}

	s.ensureCurrentKegiatan()
	return s
}

func (s *JumatBerkahService) ensureCurrentKegiatan() {
	kegiatan, _ := s.repo.GetCurrentKegiatan()
	if kegiatan == nil {
		nextFriday := s.getNextFriday(time.Now())
		newKegiatan := &models.KegiatanJumatBerkah{
			TanggalKegiatan: nextFriday,
			KuotaMaksimal:   20,
			TotalTerdaftar:  0,
			StatusKegiatan:  models.StatusKegiatanDibuka,
		}
		s.repo.CreateKegiatan(newKegiatan)
	}
}

func (s *JumatBerkahService) getNextFriday(from time.Time) time.Time {
	daysUntilFriday := (5 - int(from.Weekday()) + 7) % 7
	if daysUntilFriday == 0 {
		daysUntilFriday = 7
	}
	return from.AddDate(0, 0, daysUntilFriday)
}

func (s *JumatBerkahService) GetCurrentKegiatan() (*models.KegiatanJumatBerkah, error) {
	return s.repo.GetCurrentKegiatan()
}

func (s *JumatBerkahService) GetKegiatanByID(id string) (*models.KegiatanJumatBerkah, error) {
	return s.repo.GetKegiatanByID(id)
}

func (s *JumatBerkahService) GetAllPendaftar() ([]*models.PendaftarJumatBerkah, error) {
	return s.repo.GetAllPendaftar()
}

func (s *JumatBerkahService) GetPendaftarByKegiatan(kegiatanID string) ([]*models.PendaftarJumatBerkah, error) {
	return s.repo.GetPendaftarByKegiatan(kegiatanID)
}

func (s *JumatBerkahService) GetPendaftarByStatus(kegiatanID string, status models.StatusApproval) ([]*models.PendaftarJumatBerkah, error) {
	return s.repo.GetPendaftarByStatus(kegiatanID, status)
}

func (s *JumatBerkahService) GetPendingCount() int {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return 0
	}
	count, _ := s.repo.CountPendingByKegiatan(kegiatan.ID)
	return count
}

func (s *JumatBerkahService) GetApprovedCount() int {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return 0
	}
	count, _ := s.repo.CountApprovedByKegiatan(kegiatan.ID)
	return count
}

func (s *JumatBerkahService) CountApprovedByKegiatan(kegiatanID string) (int, error) {
	return s.repo.CountApprovedByKegiatan(kegiatanID)
}

func (s *JumatBerkahService) CreateRegistration(anakAsuhID, namaAnak, jenjang, statusAnak, rw, rt string) (*models.PendaftarJumatBerkah, error) {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return nil, fmt.Errorf("no active kegiatan")
	}

	isRegistered, _ := s.repo.IsAnakRegistered(kegiatan.ID, anakAsuhID)
	if isRegistered {
		return nil, fmt.Errorf("anak already registered")
	}

	pendaftar := &models.PendaftarJumatBerkah{
		IDKegiatan: kegiatan.ID,
		IDAnak:     anakAsuhID,
	}

	err = s.repo.CreatePendaftar(pendaftar)
	if err != nil {
		return nil, err
	}

	pendaftar.Anak = &models.AnakAsuh{
		ID:                anakAsuhID,
		NamaLengkap:       namaAnak,
		JenjangPendidikan: jenjang,
		StatusAnak:        models.StatusAnak(statusAnak),
		RT:                rt,
		RW:                rw,
	}

	return pendaftar, nil
}

func (s *JumatBerkahService) ApproveRegistration(id string) error {
	pendaftar, err := s.repo.GetPendaftarByID(id)
	if err != nil {
		return err
	}

	if pendaftar.StatusApproval != models.StatusApprovalMenunggu {
		return fmt.Errorf("registration is not pending")
	}

	return s.repo.UpdatePendaftarStatus(id, models.StatusApprovalDisetujui)
}

func (s *JumatBerkahService) RejectRegistration(id string) error {
	pendaftar, err := s.repo.GetPendaftarByID(id)
	if err != nil {
		return err
	}

	if pendaftar.StatusApproval != models.StatusApprovalMenunggu {
		return fmt.Errorf("registration is not pending")
	}

	return s.repo.UpdatePendaftarStatus(id, models.StatusApprovalDitolak)
}

func (s *JumatBerkahService) ApproveMultiple(ids []string) int {
	count, _ := s.repo.UpdateMultiplePendaftarStatus(ids, models.StatusApprovalDisetujui)
	return count
}

func (s *JumatBerkahService) RejectMultiple(ids []string) int {
	count, _ := s.repo.UpdateMultiplePendaftarStatus(ids, models.StatusApprovalDitolak)
	return count
}

func (s *JumatBerkahService) ApproveAllPending() int {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return 0
	}

	pendingList, err := s.repo.GetPendaftarByStatus(kegiatan.ID, models.StatusApprovalMenunggu)
	if err != nil {
		return 0
	}

	var ids []string
	for _, p := range pendingList {
		ids = append(ids, p.ID)
	}

	return s.ApproveMultiple(ids)
}

func (s *JumatBerkahService) UpdateQuota(quota int) error {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return fmt.Errorf("no active kegiatan")
	}

	kegiatan.KuotaMaksimal = quota
	return s.repo.UpdateKegiatan(kegiatan)
}

func (s *JumatBerkahService) ToggleForm(open bool) error {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return fmt.Errorf("no active kegiatan")
	}

	if open {
		kegiatan.StatusKegiatan = models.StatusKegiatanDibuka
	} else {
		kegiatan.StatusKegiatan = models.StatusKegiatanDitutup
	}

	return s.repo.UpdateKegiatan(kegiatan)
}

func (s *JumatBerkahService) GetRemainingQuota() int {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return 0
	}

	approved, _ := s.repo.CountApprovedByKegiatan(kegiatan.ID)
	remaining := kegiatan.KuotaMaksimal - approved
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s *JumatBerkahService) IsFormOpen() bool {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return false
	}
	return kegiatan.StatusKegiatan == models.StatusKegiatanDibuka
}

func (s *JumatBerkahService) IsAnakAsuhRegistered(anakAsuhID string) bool {
	kegiatan, err := s.repo.GetCurrentKegiatan()
	if err != nil || kegiatan == nil {
		return false
	}

	isRegistered, _ := s.repo.IsAnakRegistered(kegiatan.ID, anakAsuhID)
	return isRegistered
}

func (s *JumatBerkahService) CountKegiatan() (int, error) {
	return s.repo.CountAllKegiatan()
}

func (s *JumatBerkahService) CountPendaftar() (int, error) {
	return s.repo.CountAllPendaftar()
}

func (s *JumatBerkahService) GetPendaftarByAnakID(anakID string) ([]*models.PendaftarJumatBerkah, error) {
	return s.repo.GetPendaftarByAnakID(anakID)
}
