package services

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
	"puriyatim-app/pkg/pakasir"
)

type DonasiOnlineService struct {
	repo           *repository.DonasiOnlineRepository
	keuanganRepo   *repository.KeuanganRepository
	pakasirClient  *pakasir.Client
}

func NewDonasiOnlineService(
	repo *repository.DonasiOnlineRepository,
	keuanganRepo *repository.KeuanganRepository,
	pakasirClient *pakasir.Client,
) *DonasiOnlineService {
	return &DonasiOnlineService{
		repo:          repo,
		keuanganRepo:  keuanganRepo,
		pakasirClient: pakasirClient,
	}
}

// CreateTransactionRequest parameter untuk membuat transaksi baru
type CreateTransactionRequest struct {
	Jenis         models.JenisDonasiOnline
	NamaDonatur   string
	Nominal       float64
	PaymentMethod string
	Catatan       string
}

// CreateTransaction membuat transaksi di Pakasir dan menyimpannya ke DB
func (s *DonasiOnlineService) CreateTransaction(req CreateTransactionRequest) (*models.DonasiOnline, error) {
	if req.Nominal <= 0 {
		return nil, fmt.Errorf("nominal harus lebih dari 0")
	}
	if req.PaymentMethod == "" {
		return nil, fmt.Errorf("metode pembayaran wajib dipilih")
	}
	if req.NamaDonatur == "" {
		req.NamaDonatur = "Hamba Allah"
	}

	orderID := generateOrderID(string(req.Jenis))

	// Panggil Pakasir API
	payment, err := s.pakasirClient.CreateTransaction(
		req.PaymentMethod,
		orderID,
		int64(req.Nominal),
	)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat transaksi pembayaran: %w", err)
	}

	// Tentukan apakah VA atau QRIS
	qrString := ""
	vaNumber := ""
	if req.PaymentMethod == pakasir.MethodQRIS {
		qrString = payment.PaymentNumber
	} else {
		vaNumber = payment.PaymentNumber
	}

	donasi := &models.DonasiOnline{
		OrderID:       orderID,
		Jenis:         req.Jenis,
		NamaDonatur:   req.NamaDonatur,
		Nominal:       req.Nominal,
		PaymentMethod: req.PaymentMethod,
		Status:        models.StatusDonasiPending,
		QRString:      qrString,
		VANumber:      vaNumber,
		TotalPayment:  float64(payment.TotalPayment),
		Fee:           float64(payment.Fee),
		ExpiredAt:     &payment.ExpiredAt,
		Catatan:       req.Catatan,
	}

	if err := s.repo.Create(donasi); err != nil {
		return nil, fmt.Errorf("gagal menyimpan data transaksi: %w", err)
	}

	return donasi, nil
}

// GetByOrderID mengambil data transaksi berdasarkan order_id
func (s *DonasiOnlineService) GetByOrderID(orderID string) (*models.DonasiOnline, error) {
	return s.repo.GetByOrderID(orderID)
}

// GetAll mengambil semua transaksi (limit 0 = semua)
func (s *DonasiOnlineService) GetAll(limit int) ([]*models.DonasiOnline, error) {
	return s.repo.GetAll(limit)
}

// HandleWebhook memproses payload webhook dari Pakasir saat pembayaran berhasil.
// Validasi: amount & order_id harus cocok dengan data di DB.
// Jika valid, update status ke completed + otomatis buat record pemasukan di buku kas.
func (s *DonasiOnlineService) HandleWebhook(payload *pakasir.WebhookPayload) error {
	if payload.Status != "completed" {
		log.Printf("pakasir webhook: order %s status=%s (skip)", payload.OrderID, payload.Status)
		return nil
	}

	// Ambil data transaksi dari DB
	donasi, err := s.repo.GetByOrderID(payload.OrderID)
	if err != nil {
		return fmt.Errorf("webhook: order_id %s tidak ditemukan di DB: %w", payload.OrderID, err)
	}

	// Validasi amount agar tidak bisa dipalsukan
	if int64(donasi.Nominal) != payload.Amount {
		return fmt.Errorf("webhook: amount tidak cocok order %s (expected %d, got %d)",
			payload.OrderID, int64(donasi.Nominal), payload.Amount)
	}

	// Kalau sudah completed, skip (idempotent)
	if donasi.Status == models.StatusDonasiCompleted {
		log.Printf("pakasir webhook: order %s sudah completed, skip", payload.OrderID)
		return nil
	}

	// Update status ke completed
	now := time.Now()
	if payload.CompletedAt != nil {
		now = *payload.CompletedAt
	}
	if err := s.repo.UpdateStatus(payload.OrderID, models.StatusDonasiCompleted, &now); err != nil {
		return fmt.Errorf("webhook: gagal update status: %w", err)
	}

	// Otomatis catat ke pemasukan donasi (buku kas)
	if s.keuanganRepo != nil {
		kategori := jenisToKategori(donasi.Jenis)
		catatan := buildCatatan(donasi)

		pemasukan := &models.PemasukanDonasi{
			NamaDonatur:      donasi.NamaDonatur,
			TanggalDonasi:    now,
			Nominal:          donasi.Nominal,
			KategoriDana:     kategori,
			Catatan:          catatan,
			BuktiTransaksi:   "",
			StatusVerifikasi: models.StatusVerifikasiVerified,
		}
		pemasukan.ID = "DON-" + donasi.OrderID
		if err := s.keuanganRepo.CreatePemasukan(pemasukan); err != nil {
			// Log tapi jangan return error agar status tetap tersimpan
			log.Printf("webhook: gagal buat pemasukan untuk order %s: %v", payload.OrderID, err)
		} else {
			log.Printf("webhook: pemasukan dicatat untuk order %s nominal %.0f", payload.OrderID, donasi.Nominal)
		}
	}

	log.Printf("pakasir webhook: order %s COMPLETED, nominal=%.0f", payload.OrderID, donasi.Nominal)
	return nil
}

// CheckAndUpdateStatus mengecek status ke Pakasir API secara aktif (polling dari browser)
func (s *DonasiOnlineService) CheckAndUpdateStatus(orderID string) (*models.DonasiOnline, error) {
	donasi, err := s.repo.GetByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	// Kalau sudah final, tidak perlu cek ulang
	if donasi.Status == models.StatusDonasiCompleted ||
		donasi.Status == models.StatusDonasiCancelled {
		return donasi, nil
	}

	// Kalau sudah expired di sisi kita, tandai
	if donasi.IsExpired() && donasi.Status == models.StatusDonasiPending {
		_ = s.repo.UpdateStatus(orderID, models.StatusDonasiExpired, nil)
		donasi.Status = models.StatusDonasiExpired
		return donasi, nil
	}

	// Cek ke Pakasir
	detail, err := s.pakasirClient.GetTransactionDetail(orderID, int64(donasi.Nominal))
	if err != nil {
		// Jangan fail hard, return data DB yang ada
		log.Printf("check status pakasir error (order %s): %v", orderID, err)
		return donasi, nil
	}

	if detail.Status == "completed" && donasi.Status != models.StatusDonasiCompleted {
		now := time.Now()
		if detail.CompletedAt != nil {
			now = *detail.CompletedAt
		}
		_ = s.repo.UpdateStatus(orderID, models.StatusDonasiCompleted, &now)
		donasi.Status = models.StatusDonasiCompleted
		donasi.CompletedAt = &now

		// Buat pemasukan di buku kas
		if s.keuanganRepo != nil {
			kategori := jenisToKategori(donasi.Jenis)
			pemasukan := &models.PemasukanDonasi{
				NamaDonatur:      donasi.NamaDonatur,
				TanggalDonasi:    now,
				Nominal:          donasi.Nominal,
				KategoriDana:     kategori,
				Catatan:          buildCatatan(donasi),
				StatusVerifikasi: models.StatusVerifikasiVerified,
			}
			pemasukan.ID = "DON-" + donasi.OrderID
			if err := s.keuanganRepo.CreatePemasukan(pemasukan); err != nil {
				log.Printf("check_status: gagal buat pemasukan order %s: %v", orderID, err)
			}
		}
	}

	return donasi, nil
}

// CancelTransaction membatalkan transaksi pending
func (s *DonasiOnlineService) CancelTransaction(orderID string) error {
	donasi, err := s.repo.GetByOrderID(orderID)
	if err != nil {
		return err
	}
	if donasi.Status != models.StatusDonasiPending {
		return fmt.Errorf("hanya transaksi pending yang dapat dibatalkan")
	}

	if err := s.pakasirClient.CancelTransaction(orderID, int64(donasi.Nominal)); err != nil {
		log.Printf("cancel pakasir error (order %s): %v", orderID, err)
		// Tetap update DB meski Pakasir gagal
	}

	return s.repo.UpdateStatus(orderID, models.StatusDonasiCancelled, nil)
}

// MarkExpiredOld menandai semua transaksi pending lama sebagai expired
func (s *DonasiOnlineService) MarkExpiredOld() (int64, error) {
	return s.repo.MarkExpiredOld()
}

// IsPakasirConfigured apakah Pakasir sudah dikonfigurasi
func (s *DonasiOnlineService) IsPakasirConfigured() bool {
	return s.pakasirClient.IsConfigured()
}

// ---- helpers ----

func generateOrderID(jenis string) string {
	prefix := map[string]string{
		"donasi":       "DON",
		"zakat":        "ZKT",
		"jumat_berkah": "JMT",
	}
	p, ok := prefix[jenis]
	if !ok {
		p = "TRX"
	}
	ts := time.Now().Format("060102150405") // YYMMDDHHmmss
	rnd := rand.Intn(9000) + 1000           // 4 digit random
	return fmt.Sprintf("%s-%s-%d", p, ts, rnd)
}

func jenisToKategori(jenis models.JenisDonasiOnline) models.KategoriDana {
	switch jenis {
	case models.JenisDonasiZakat:
		return models.KategoriDanaZakat
	case models.JenisDonasiJumat:
		return models.KategoriDanaInfaq
	default:
		return models.KategoriDanaSedekah
	}
}

func buildCatatan(d *models.DonasiOnline) string {
	parts := []string{
		"Donasi Online via Pakasir",
		"Metode: " + pakasir.PaymentMethodLabels[d.PaymentMethod],
		"Order: " + d.OrderID,
	}
	if d.Catatan != "" {
		parts = append(parts, d.Catatan)
	}
	return strings.Join(parts, " | ")
}
