package services

import (
	"sort"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
)

type KeuanganService struct {
	repo *repository.KeuanganRepository
}

func NewKeuanganService(repo *repository.KeuanganRepository) *KeuanganService {
	return &KeuanganService{repo: repo}
}

func (s *KeuanganService) CreatePemasukan(p *models.PemasukanDonasi) error {
	if p.Nominal <= 0 {
		return ErrInvalidNominal
	}
	if p.TanggalDonasi.IsZero() {
		p.TanggalDonasi = time.Now()
	}
	if p.NamaDonatur == "" {
		p.NamaDonatur = "Hamba Allah"
	}
	if p.StatusVerifikasi == "" {
		p.StatusVerifikasi = models.StatusVerifikasiPending
	}
	return s.repo.CreatePemasukan(p)
}

func (s *KeuanganService) GetPemasukanByID(id string) (*models.PemasukanDonasi, error) {
	return s.repo.GetPemasukanByID(id)
}

func (s *KeuanganService) GetAllPemasukan() ([]*models.PemasukanDonasi, error) {
	return s.repo.GetAllPemasukan()
}

func (s *KeuanganService) CreatePengeluaran(p *models.Pengeluaran) error {
	if p.Nominal <= 0 {
		return ErrInvalidNominal
	}
	if p.TanggalPengeluaran.IsZero() {
		p.TanggalPengeluaran = time.Now()
	}
	if p.Keterangan == "" {
		return ErrKeteranganRequired
	}
	return s.repo.CreatePengeluaran(p)
}

func (s *KeuanganService) GetPengeluaranByID(id string) (*models.Pengeluaran, error) {
	return s.repo.GetPengeluaranByID(id)
}

func (s *KeuanganService) GetAllPengeluaran() ([]*models.Pengeluaran, error) {
	return s.repo.GetAllPengeluaran()
}

func (s *KeuanganService) GetStatistics() (*models.KeuanganStats, error) {
	totalPemasukan, err := s.repo.GetTotalPemasukan()
	if err != nil {
		return nil, err
	}

	totalPengeluaran, err := s.repo.GetTotalPengeluaran()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	var lastYear, lastMonth int
	if currentMonth == 1 {
		lastYear = currentYear - 1
		lastMonth = 12
	} else {
		lastYear = currentYear
		lastMonth = currentMonth - 1
	}

	pemasukanBulanIni, err := s.repo.GetTotalPemasukanByMonth(currentYear, currentMonth)
	if err != nil {
		return nil, err
	}

	pengeluaranBulanIni, err := s.repo.GetTotalPengeluaranByMonth(currentYear, currentMonth)
	if err != nil {
		return nil, err
	}

	pemasukanBulanLalu, err := s.repo.GetTotalPemasukanByMonth(lastYear, lastMonth)
	if err != nil {
		return nil, err
	}

	pengeluaranBulanLalu, err := s.repo.GetTotalPengeluaranByMonth(lastYear, lastMonth)
	if err != nil {
		return nil, err
	}

	pemasukanChange := calculatePercentChange(pemasukanBulanLalu, pemasukanBulanIni)
	pengeluaranChange := calculatePercentChange(pengeluaranBulanLalu, pengeluaranBulanIni)

	return &models.KeuanganStats{
		TotalSaldo:           totalPemasukan - totalPengeluaran,
		TotalPemasukan:       totalPemasukan,
		TotalPengeluaran:     totalPengeluaran,
		PemasukanBulanIni:    pemasukanBulanIni,
		PengeluaranBulanIni:  pengeluaranBulanIni,
		PemasukanBulanLalu:   pemasukanBulanLalu,
		PengeluaranBulanLalu: pengeluaranBulanLalu,
		PemasukanChange:      pemasukanChange,
		PengeluaranChange:    pengeluaranChange,
	}, nil
}

func calculatePercentChange(old, new float64) float64 {
	if old == 0 {
		if new > 0 {
			return 100
		}
		return 0
	}
	return ((new - old) / old) * 100
}

func (s *KeuanganService) GetBukuKas() ([]models.KasTransaction, error) {
	pemasukanList, err := s.repo.GetAllPemasukan()
	if err != nil {
		return nil, err
	}

	pengeluaranList, err := s.repo.GetAllPengeluaran()
	if err != nil {
		return nil, err
	}

	var transactions []models.KasTransaction

	for _, p := range pemasukanList {
		transactions = append(transactions, models.KasTransaction{
			ID:        p.ID,
			Tanggal:   p.TanggalDonasi,
			CreatedAt: p.CreatedAt,
			Deskripsi: p.NamaDonatur,
			Kategori:  string(p.KategoriDana),
			Jumlah:    p.Nominal,
			Type:      "masuk",
			Donatur:   p.NamaDonatur,
			Status:    p.StatusVerifikasi,
		})
	}

	for _, p := range pengeluaranList {
		anakName := ""
		if p.Anak != nil {
			anakName = p.Anak.NamaLengkap
		}
		transactions = append(transactions, models.KasTransaction{
			ID:        p.ID,
			Tanggal:   p.TanggalPengeluaran,
			CreatedAt: p.CreatedAt,
			Deskripsi: p.Keterangan,
			Kategori:  "Pengeluaran",
			Jumlah:    p.Nominal,
			Type:      "keluar",
			AnakAsuh:  anakName,
		})
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.After(transactions[j].CreatedAt)
	})

	return transactions, nil
}

func (s *KeuanganService) GetBukuKasFiltered(month, transType string) ([]models.KasTransaction, float64, float64, error) {
	var transactions []models.KasTransaction
	var totalPemasukan, totalPengeluaran float64

	if transType == "" || transType == "in" {
		pemasukanList, err := s.repo.GetPemasukanByMonthStr(month)
		if err != nil {
			return nil, 0, 0, err
		}

		for _, p := range pemasukanList {
			transactions = append(transactions, models.KasTransaction{
				ID:        p.ID,
				Tanggal:   p.TanggalDonasi,
				CreatedAt: p.CreatedAt,
				Deskripsi: p.NamaDonatur,
				Kategori:  string(p.KategoriDana),
				Jumlah:    p.Nominal,
				Type:      "masuk",
				Donatur:   p.NamaDonatur,
				Status:    p.StatusVerifikasi,
			})
			if p.StatusVerifikasi == models.StatusVerifikasiVerified {
				totalPemasukan += p.Nominal
			}
		}
	}

	if transType == "" || transType == "out" {
		pengeluaranList, err := s.repo.GetPengeluaranByMonthStr(month)
		if err != nil {
			return nil, 0, 0, err
		}

		for _, p := range pengeluaranList {
			anakName := ""
			if p.Anak != nil {
				anakName = p.Anak.NamaLengkap
			}
			transactions = append(transactions, models.KasTransaction{
				ID:        p.ID,
				Tanggal:   p.TanggalPengeluaran,
				CreatedAt: p.CreatedAt,
				Deskripsi: p.Keterangan,
				Kategori:  "Pengeluaran",
				Jumlah:    p.Nominal,
				Type:      "keluar",
				AnakAsuh:  anakName,
			})
			totalPengeluaran += p.Nominal
		}
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.After(transactions[j].CreatedAt)
	})

	return transactions, totalPemasukan, totalPengeluaran, nil
}

func (s *KeuanganService) CountPemasukan() (int, error) {
	return s.repo.CountPemasukan()
}

func (s *KeuanganService) CountPengeluaran() (int, error) {
	return s.repo.CountPengeluaran()
}

func (s *KeuanganService) GetAllDonatur() ([]models.Donatur, error) {
	return s.repo.GetAllDonatur()
}

func (s *KeuanganService) CreateDonatur(d *models.Donatur) error {
	if d.NamaDonatur == "" {
		return &ValidationError{Message: "Nama donatur wajib diisi"}
	}
	if d.NoTelepon == "" {
		return &ValidationError{Message: "No. telepon wajib diisi"}
	}
	if d.TipeDonatur == "" {
		d.TipeDonatur = models.TipeDonaturIndividu
	}
	return s.repo.CreateDonatur(d)
}

func (s *KeuanganService) DeletePemasukan(id string) error {
	return s.repo.DeletePemasukan(id)
}

func (s *KeuanganService) DeletePengeluaran(id string) error {
	return s.repo.DeletePengeluaran(id)
}

func (s *KeuanganService) UpdatePemasukan(p *models.PemasukanDonasi) error {
	if p.Nominal <= 0 {
		return ErrInvalidNominal
	}
	if p.TanggalDonasi.IsZero() {
		p.TanggalDonasi = time.Now()
	}
	if p.NamaDonatur == "" {
		p.NamaDonatur = "Hamba Allah"
	}
	if p.StatusVerifikasi == "" {
		existing, err := s.repo.GetPemasukanByID(p.ID)
		if err == nil {
			p.StatusVerifikasi = existing.StatusVerifikasi
		} else {
			p.StatusVerifikasi = models.StatusVerifikasiPending
		}
	}
	return s.repo.UpdatePemasukan(p)
}

func (s *KeuanganService) UpdatePengeluaran(p *models.Pengeluaran) error {
	if p.Nominal <= 0 {
		return ErrInvalidNominal
	}
	if p.TanggalPengeluaran.IsZero() {
		p.TanggalPengeluaran = time.Now()
	}
	if p.Keterangan == "" {
		return ErrKeteranganRequired
	}
	return s.repo.UpdatePengeluaran(p)
}

func (s *KeuanganService) GetPengeluaranByAnakID(anakID string) ([]*models.Pengeluaran, error) {
	return s.repo.GetPengeluaranByAnakID(anakID)
}

func (s *KeuanganService) VerifyPemasukan(id string) error {
	return s.repo.VerifyPemasukan(id)
}

func (s *KeuanganService) RejectPemasukan(id string) error {
	return s.repo.RejectPemasukan(id)
}

func (s *KeuanganService) GetPemasukanByNomorHP(nomorHP string) ([]*models.PemasukanDonasi, error) {
	return s.repo.GetPemasukanByNomorHP(nomorHP)
}

var ErrInvalidNominal = &ValidationError{Message: "Nominal harus lebih dari 0"}
var ErrKeteranganRequired = &ValidationError{Message: "Keterangan wajib diisi"}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// KeuanganDashboardData berisi semua data keuangan yang dibutuhkan oleh halaman dashboard.
type KeuanganDashboardData struct {
	Stats           *models.KeuanganStats
	KategoriLabels  []string
	KategoriValues  []float64
	BulanLabels     []string
	BulanPemasukan  []float64
	BulanPengeluaran []float64
	PendingCount    int
}

// urutan tampil kategori di chart
var kategoriUrutan = []string{"Zakat", "Infaq", "Sedekah", "Wakaf", "Lainnya"}

// GetDashboardKeuangan mengumpulkan semua data statistik keuangan untuk dashboard.
func (s *KeuanganService) GetDashboardKeuangan() (*KeuanganDashboardData, error) {
	stats, _ := s.GetStatistics()

	// Pemasukan by kategori — pakai urutan tetap
	rawKategori, _ := s.repo.GetPemasukanByKategori()
	// Buat map untuk lookup cepat
	kategoriMap := make(map[string]float64, len(rawKategori))
	for _, kt := range rawKategori {
		kategoriMap[kt.Kategori] = kt.Total
	}

	var kategoriLabels []string
	var kategoriValues []float64
	for _, k := range kategoriUrutan {
		if v, ok := kategoriMap[k]; ok && v > 0 {
			kategoriLabels = append(kategoriLabels, k)
			kategoriValues = append(kategoriValues, v)
		}
	}
	// Tambahkan kategori lain di luar daftar urutan (custom)
	for _, kt := range rawKategori {
		found := false
		for _, k := range kategoriUrutan {
			if kt.Kategori == k {
				found = true
				break
			}
		}
		if !found && kt.Total > 0 {
			kategoriLabels = append(kategoriLabels, kt.Kategori)
			kategoriValues = append(kategoriValues, kt.Total)
		}
	}

	// Tren bulanan 6 bulan terakhir
	bulanan, _ := s.repo.GetPemasukanPengeluaranBulanan(6)
	bulanLabels := make([]string, len(bulanan))
	bulanPemasukan := make([]float64, len(bulanan))
	bulanPengeluaran := make([]float64, len(bulanan))
	for i, b := range bulanan {
		bulanLabels[i] = b.Label
		bulanPemasukan[i] = b.Pemasukan
		bulanPengeluaran[i] = b.Pengeluaran
	}

	pendingCount, _ := s.repo.GetCountPemasukanPending()

	return &KeuanganDashboardData{
		Stats:            stats,
		KategoriLabels:   kategoriLabels,
		KategoriValues:   kategoriValues,
		BulanLabels:      bulanLabels,
		BulanPemasukan:   bulanPemasukan,
		BulanPengeluaran: bulanPengeluaran,
		PendingCount:     pendingCount,
	}, nil
}
