package handlers

import (
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
}

func NewDashboardHandler(
	cfg *config.Config,
	anakAsuhService *services.AnakAsuhService,
	keuanganService *services.KeuanganService,
	jumatBerkahService *services.JumatBerkahService,
) *DashboardHandler {
	return &DashboardHandler{
		cfg:                cfg,
		anakAsuhService:    anakAsuhService,
		keuanganService:    keuanganService,
		jumatBerkahService: jumatBerkahService,
	}
}

type DashboardData struct {
	User                    models.Pengurus
	Title                   string
	PageTitle               string
	ActivePage              string
	Stats                   DashboardStats
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
	user := models.Pengurus{
		ID:          "1",
		NamaLengkap: "Admin",
		Email:       "admin@puriyatim.com",
		Peran:       models.PeranSuperadmin,
	}

	anakAsuhCount, _ := h.anakAsuhService.Count()

	keuanganStats, _ := h.keuanganService.GetStatistics()
	kasTersedia := formatRupiah(keuanganStats.TotalSaldo)

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
		ArtikelCount:       "0",
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
			Jumlah:    formatRupiah(t.Jumlah),
		})
	}

	data := DashboardData{
		User:                    user,
		Title:                   "Dashboard Admin",
		PageTitle:               "Ringkasan Hari Ini",
		ActivePage:              "dashboard",
		Stats:                   stats,
		PendingJumatBerkah:      pendingJumatBerkah,
		RecentKas:               recentKas,
		NotificationCount:       pendingCount,
		PendingJumatBerkahCount: pendingCount,
	}

	return c.Render(http.StatusOK, "admin/dashboard.html", data)
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration.Minutes() < 1 {
		return "Baru saja"
	} else if duration.Minutes() < 60 {
		return fmt.Sprintf("%d menit yang lalu", int(duration.Minutes()))
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%d jam yang lalu", int(duration.Hours()))
	} else if duration.Hours() < 48 {
		return "Kemarin"
	} else {
		return t.Format("02 Jan 2006")
	}
}

func getInitials(name string) string {
	if len(name) < 2 {
		return name
	}
	parts := splitWords(name)
	if len(parts) >= 2 {
		return string(parts[0][0]) + string(parts[1][0])
	}
	return name[:2]
}

func splitWords(s string) []string {
	var words []string
	word := ""
	for _, r := range s {
		if r == ' ' {
			if word != "" {
				words = append(words, word)
				word = ""
			}
		} else {
			word += string(r)
		}
	}
	if word != "" {
		words = append(words, word)
	}
	return words
}

func formatRupiah(amount float64) string {
	amountInt := int64(amount)
	str := fmt.Sprintf("%d", amountInt)

	var result []byte
	n := 0
	for i := len(str) - 1; i >= 0; i-- {
		if n > 0 && n%3 == 0 {
			result = append([]byte{'.'}, result...)
		}
		result = append([]byte{str[i]}, result...)
		n++
	}

	return string(result)
}
