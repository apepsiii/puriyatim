package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"puriyatim-app/internal/config"
	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	cfg                *config.Config
	anakAsuhService    *services.AnakAsuhService
	keuanganService    *services.KeuanganService
	jumatBerkahService *services.JumatBerkahService
	artikelService     *services.ArtikelService
}

func NewDashboardHandler(
	cfg *config.Config,
	anakAsuhService *services.AnakAsuhService,
	keuanganService *services.KeuanganService,
	jumatBerkahService *services.JumatBerkahService,
	artikelService *services.ArtikelService,
) *DashboardHandler {
	return &DashboardHandler{
		cfg:                cfg,
		anakAsuhService:    anakAsuhService,
		keuanganService:    keuanganService,
		jumatBerkahService: jumatBerkahService,
		artikelService:     artikelService,
	}
}

type DashboardData struct {
	User                    models.Pengurus
	Title                   string
	PageTitle               string
	ActivePage              string
	Stats                   DashboardStats
	KeuanganChart           KeuanganChartData
	PendingJumatBerkah      []PendingJumatBerkahItem
	RecentKas               []RecentKasItem
	NotificationCount       int
	PendingJumatBerkahCount int
}

type DashboardStats struct {
	AnakAsuhCount      string
	KasTersedia        string
	PendingJumatBerkah int
	KuotaJumatBerkah   int
	ArtikelCount       string
}

// KeuanganChartData menyimpan data chart & summary keuangan untuk dashboard.
type KeuanganChartData struct {
	TotalPemasukan      string
	TotalPengeluaran    string
	PemasukanBulanIni   string
	PengeluaranBulanIni string
	PendingCount        int
	KategoriLabelsJSON  string // JSON array string untuk Chart.js
	KategoriValuesJSON  string
	BulanLabelsJSON     string
	BulanPemasukanJSON  string
	BulanPengeluaranJSON string
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
	// Ambil data user dari JWT context yang di-set middleware
	userNama, _ := c.Get("user_nama").(string)
	userEmail, _ := c.Get("user_email").(string)
	userPeran, _ := c.Get("user_role").(models.PeranPengurus)
	userID, _ := c.Get("user_id").(string)
	if userNama == "" {
		userNama = "Admin"
	}
	user := models.Pengurus{
		ID:          userID,
		NamaLengkap: userNama,
		Email:       userEmail,
		Peran:       userPeran,
	}

	anakAsuhCount, _ := h.anakAsuhService.Count()

	// Artikel terbit
	artikelCount, _ := h.artikelService.CountByStatus(models.StatusPublikasiTerbit)

	// Ambil data chart & summary keuangan
	keuanganDash, _ := h.keuanganService.GetDashboardKeuangan()
	kasTersedia := FormatRupiah(keuanganDash.Stats.TotalSaldo)

	// Serialize slice ke JSON untuk diteruskan ke template
	toJSONStr := func(v interface{}) string {
		b, _ := json.Marshal(v)
		return string(b)
	}

	keuanganChart := KeuanganChartData{
		TotalPemasukan:       FormatRupiah(keuanganDash.Stats.TotalPemasukan),
		TotalPengeluaran:     FormatRupiah(keuanganDash.Stats.TotalPengeluaran),
		PemasukanBulanIni:    FormatRupiah(keuanganDash.Stats.PemasukanBulanIni),
		PengeluaranBulanIni:  FormatRupiah(keuanganDash.Stats.PengeluaranBulanIni),
		PendingCount:         keuanganDash.PendingCount,
		KategoriLabelsJSON:   toJSONStr(keuanganDash.KategoriLabels),
		KategoriValuesJSON:   toJSONStr(keuanganDash.KategoriValues),
		BulanLabelsJSON:      toJSONStr(keuanganDash.BulanLabels),
		BulanPemasukanJSON:   toJSONStr(keuanganDash.BulanPemasukan),
		BulanPengeluaranJSON: toJSONStr(keuanganDash.BulanPengeluaran),
	}

	pendingCount := h.jumatBerkahService.GetPendingCount()
	kuota := 20

	kegiatan, _ := h.jumatBerkahService.GetCurrentKegiatan()
	if kegiatan != nil {
		kuota = kegiatan.KuotaMaksimal
	}

	stats := DashboardStats{
		AnakAsuhCount:      fmt.Sprintf("%d", anakAsuhCount),
		KasTersedia:        kasTersedia,
		PendingJumatBerkah: pendingCount,
		KuotaJumatBerkah:   kuota,
		ArtikelCount:       fmt.Sprintf("%d", artikelCount),
	}

	pendingJumatBerkah := []PendingJumatBerkahItem{}

	if kegiatan != nil {
		pendingList, _ := h.jumatBerkahService.GetPendaftarByStatus(kegiatan.ID, models.StatusApprovalMenunggu)
		for i, p := range pendingList {
			if i >= 5 {
				break
			}

			nama := ""
			jenjang := ""
			status := ""
			wilayah := ""

			if p.Anak != nil {
				nama = p.Anak.NamaLengkap
				jenjang = p.Anak.JenjangPendidikan
				status = string(p.Anak.StatusAnak)
				wilayah = fmt.Sprintf("RT %s / RW %s", p.Anak.RT, p.Anak.RW)
			}

			pendingJumatBerkah = append(pendingJumatBerkah, PendingJumatBerkahItem{
				ID:          p.ID,
				Nama:        nama,
				Jenjang:     jenjang,
				Status:      status,
				Wilayah:     wilayah,
				WaktuDaftar: formatTimeAgo(p.WaktuSubmit),
				Initials:    getInitials(nama),
			})
		}
	}

	recentKas := []RecentKasItem{}

	kasTransactions, _ := h.keuanganService.GetBukuKas()
	for i, t := range kasTransactions {
		if i >= 5 {
			break
		}

		kategori := t.Kategori
		if t.Type == "keluar" {
			kategori = "Pengeluaran"
		}

		recentKas = append(recentKas, RecentKasItem{
			Type:      t.Type,
			Deskripsi: t.Deskripsi,
			Kategori:  kategori,
			Tanggal:   t.Tanggal.Format("02 Jan 2006"),
			Jumlah:    FormatRupiah(t.Jumlah),
		})
	}

	data := DashboardData{
		User:                    user,
		Title:                   "Dashboard Admin",
		PageTitle:               "Ringkasan Hari Ini",
		ActivePage:              "dashboard",
		Stats:                   stats,
		KeuanganChart:           keuanganChart,
		PendingJumatBerkah:      pendingJumatBerkah,
		RecentKas:               recentKas,
		NotificationCount:       pendingCount,
		PendingJumatBerkahCount: pendingCount,
	}

	return c.Render(http.StatusOK, "admin/dashboard.html", data)
}

// Shim functions — delegate ke helpers.go (DRY)
func formatTimeAgo(t time.Time) string   { return FormatTimeAgo(t) }
func getInitials(name string) string     { return GetInitials(name) }
func formatRupiah(amount float64) string { return FormatRupiah(amount) }
