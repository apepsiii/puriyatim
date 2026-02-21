package handlers

import (
	"net/http"
	"puriyatim-app/internal/config"
	"puriyatim-app/internal/models"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	cfg *config.Config
}

func NewDashboardHandler(cfg *config.Config) *DashboardHandler {
	return &DashboardHandler{cfg: cfg}
}

type DashboardData struct {
	User                 models.Pengurus
	Title                string
	PageTitle            string
	ActivePage           string
	Stats                DashboardStats
	PendingJumatBerkah   []PendingJumatBerkahItem
	RecentKas            []RecentKasItem
	NotificationCount    int
	PendingJumatBerkahCount int
}

type DashboardStats struct {
	AnakAsuhCount        string
	KasTersedia          string
	PendingJumatBerkah   int
	KuotaJumatBerkah     int
	ArtikelCount         string
}

type PendingJumatBerkahItem struct {
	ID          string
	Nama        string
	Jenjang     string
	Status      string
	Wilayah     string
	WaktuDaftar string
	Initials    string
}

type RecentKasItem struct {
	Type      string
	Deskripsi string
	Kategori  string
	Tanggal   string
	Jumlah    string
}

func (h *DashboardHandler) Dashboard(c echo.Context) error {
	// Get user from session (mock data for now)
	user := models.Pengurus{
		ID:          "1",
		NamaLengkap: "Budi Admin",
		Email:       "admin@puriyatim.com",
		Peran:       models.PeranSuperadmin,
	}

	// Mock dashboard stats
	stats := DashboardStats{
		AnakAsuhCount:      "142",
		KasTersedia:        "12,5 Jt",
		PendingJumatBerkah: 5,
		KuotaJumatBerkah:   20,
		ArtikelCount:       "28",
	}

	// Mock pending Jumat Berkah data
	pendingJumatBerkah := []PendingJumatBerkahItem{
		{
			ID:          "1",
			Nama:        "Budi Santoso",
			Jenjang:     "SD",
			Status:      "Yatim",
			Wilayah:     "RT 01 / RW 01",
			WaktuDaftar: "Hari ini, 08:30 WIB",
			Initials:    "BS",
		},
		{
			ID:          "2",
			Nama:        "Siti Aminah",
			Jenjang:     "SMP",
			Status:      "Dhuafa",
			Wilayah:     "RT 02 / RW 01",
			WaktuDaftar: "Hari ini, 09:15 WIB",
			Initials:    "SA",
		},
	}

	// Mock recent kas data
	recentKas := []RecentKasItem{
		{
			Type:      "masuk",
			Deskripsi: "Hamba Allah (Transfer)",
			Kategori:  "Donasi Sedekah Umum",
			Tanggal:   "10 Feb 2026",
			Jumlah:    "500.000",
		},
		{
			Type:      "masuk",
			Deskripsi: "PT. Maju Makmur",
			Kategori:  "Zakat Perusahaan",
			Tanggal:   "09 Feb 2026",
			Jumlah:    "2.000.000",
		},
		{
			Type:      "keluar",
			Deskripsi: "Bayar SPP SMK (2 Anak)",
			Kategori:  "Pengeluaran - Pendidikan",
			Tanggal:   "08 Feb 2026",
			Jumlah:    "600.000",
		},
	}

	data := DashboardData{
		User:                    user,
		Title:                   "Dashboard Admin",
		PageTitle:               "Ringkasan Hari Ini",
		ActivePage:              "dashboard",
		Stats:                   stats,
		PendingJumatBerkah:      pendingJumatBerkah,
		RecentKas:               recentKas,
		NotificationCount:       1,
		PendingJumatBerkahCount: stats.PendingJumatBerkah,
	}

	return c.Render(http.StatusOK, "admin/dashboard.html", data)
}

func (h *DashboardHandler) ApproveJumatBerkah(c echo.Context) error {
	id := c.Param("id")
	
	// TODO: Implement actual approval logic in database
	// For now, just return success response
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Pendaftar berhasil disetujui",
		"id":      id,
	})
}

func (h *DashboardHandler) RejectJumatBerkah(c echo.Context) error {
	id := c.Param("id")
	
	// TODO: Implement actual rejection logic in database
	// For now, just return success response
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Pendaftar berhasil ditolak",
		"id":      id,
	})
}

func (h *DashboardHandler) ApproveAllJumatBerkah(c echo.Context) error {
	// TODO: Implement actual bulk approval logic in database
	// For now, just return success response
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Semua pendaftar berhasil disetujui",
	})
}