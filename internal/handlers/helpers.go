package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"puriyatim-app/internal/models"

	"github.com/labstack/echo/v4"
)

// NamaBulan adalah nama bulan dalam Bahasa Indonesia.
var NamaBulan = [12]string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// ─── Format Helpers ────────────────────────────────────────────────────────────

// FormatRupiah memformat angka float64 menjadi string angka dengan pemisah titik.
// Contoh: 1500000 → "1.500.000"
func FormatRupiah(amount float64) string {
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

// FormatDate memformat time.Time ke string tanggal Indonesia.
// Contoh: "5 Januari 2024"
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return fmt.Sprintf("%d %s %d", t.Day(), NamaBulan[t.Month()-1], t.Year())
}

// FormatMonthYear memformat time.Time ke "Januari 2024".
func FormatMonthYear(t time.Time) string {
	return NamaBulan[int(t.Month())-1] + " " + fmt.Sprintf("%d", t.Year())
}

// FormatTimeAgo memformat time.Time ke string relatif Bahasa Indonesia.
// Contoh: "5 menit yang lalu", "Kemarin", "02 Jan 2006"
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	switch {
	case duration.Minutes() < 1:
		return "Baru saja"
	case duration.Minutes() < 60:
		return fmt.Sprintf("%d menit yang lalu", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%d jam yang lalu", int(duration.Hours()))
	case duration.Hours() < 48:
		return "Kemarin"
	default:
		return t.Format("02 Jan 2006")
	}
}

// FormatPercentChange memformat float64 menjadi string persentase dengan tanda.
// Contoh: 5.3 → "+5.3%", -2.1 → "-2.1%"
func FormatPercentChange(change float64) string {
	if change == 0 {
		return "0%"
	}
	sign := ""
	if change > 0 {
		sign = "+"
	}
	return fmt.Sprintf("%s%.1f%%", sign, change)
}

// ─── Avatar & Styling Helpers ──────────────────────────────────────────────────

// GetInitials mengambil dua huruf pertama dari nama (huruf kapital nama pertama & kedua).
func GetInitials(name string) string {
	if len(name) == 0 {
		return "??"
	}
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "??"
	}
	if len(parts) == 1 {
		if len(parts[0]) < 2 {
			return strings.ToUpper(parts[0])
		}
		return strings.ToUpper(parts[0][:2])
	}
	return strings.ToUpper(string(parts[0][0])) + strings.ToUpper(string(parts[1][0]))
}

// getAvatarStyle mengembalikan pasangan (bg-color, text-color) Tailwind CSS
// berdasarkan hash sederhana dari nama.
func getAvatarStyle(name string) (bg, text string) {
	colors := [][2]string{
		{"bg-blue-100", "text-blue-700"},
		{"bg-green-100", "text-green-700"},
		{"bg-purple-100", "text-purple-700"},
		{"bg-amber-100", "text-amber-700"},
		{"bg-pink-100", "text-pink-700"},
		{"bg-teal-100", "text-teal-700"},
		{"bg-indigo-100", "text-indigo-700"},
		{"bg-orange-100", "text-orange-700"},
	}
	if name == "" {
		return colors[0][0], colors[0][1]
	}
	idx := 0
	for _, r := range name {
		idx += int(r)
	}
	c := colors[idx%len(colors)]
	return c[0], c[1]
}

// GetChangeCSS mengembalikan class Tailwind untuk warna perubahan persentase.
// isPemasukan=true → hijau jika naik, merah jika turun.
func GetChangeCSS(change float64, isPemasukan bool) string {
	switch {
	case change > 0:
		if isPemasukan {
			return "text-emerald-600"
		}
		return "text-red-500"
	case change < 0:
		if isPemasukan {
			return "text-red-500"
		}
		return "text-emerald-600"
	default:
		return "text-gray-500"
	}
}

// GetChangeIcon mengembalikan class FontAwesome icon untuk arah perubahan.
func GetChangeIcon(change float64) string {
	switch {
	case change > 0:
		return "fa-arrow-up"
	case change < 0:
		return "fa-arrow-down"
	default:
		return "fa-minus"
	}
}

// ApprovalStatusCSS mengembalikan class CSS Tailwind untuk status approval Jumat Berkah.
func ApprovalStatusCSS(status models.StatusApproval) string {
	switch status {
	case models.StatusApprovalMenunggu:
		return "bg-orange-50 text-orange-700 border border-orange-200"
	case models.StatusApprovalDisetujui:
		return "bg-emerald-50 text-emerald-700 border border-emerald-200"
	case models.StatusApprovalDitolak:
		return "bg-red-50 text-red-700 border border-red-200"
	default:
		return "bg-gray-50 text-gray-700 border border-gray-200"
	}
}

// ApprovalStatusBgText mengembalikan pasangan (bg, text) CSS untuk badge status approval.
func ApprovalStatusBgText(status models.StatusApproval) (bg, text string) {
	switch status {
	case models.StatusApprovalDisetujui:
		return "bg-green-100", "text-green-700"
	case models.StatusApprovalDitolak:
		return "bg-red-100", "text-red-700"
	default:
		return "bg-yellow-100", "text-yellow-700"
	}
}

// StatusVerifikasiLabel mengembalikan label string untuk status verifikasi pemasukan.
func StatusVerifikasiLabel(status models.StatusVerifikasiPemasukan) string {
	switch status {
	case models.StatusVerifikasiPending:
		return "Pending"
	case models.StatusVerifikasiVerified:
		return "Verified"
	default:
		return "Verified"
	}
}

// StatusVerifikasiCSS mengembalikan class CSS Tailwind untuk badge status verifikasi.
func StatusVerifikasiCSS(status models.StatusVerifikasiPemasukan) string {
	switch status {
	case models.StatusVerifikasiPending:
		return "bg-amber-50 text-amber-700 border border-amber-100"
	default:
		return "bg-emerald-50 text-emerald-700 border border-emerald-100"
	}
}

// KategoriDanaCSS mengembalikan class CSS Tailwind berdasarkan kategori dana dan tipe transaksi.
func KategoriDanaCSS(tipe, kategori string) string {
	if tipe != "masuk" {
		return "bg-orange-50 text-orange-700 border border-orange-100"
	}
	switch kategori {
	case "Zakat":
		return "bg-purple-50 text-purple-700 border border-purple-100"
	case "Infaq":
		return "bg-blue-50 text-blue-700 border border-blue-100"
	case "Sedekah":
		return "bg-emerald-50 text-emerald-700 border border-emerald-100"
	case "Wakaf":
		return "bg-amber-50 text-amber-700 border border-amber-100"
	default:
		return "bg-gray-50 text-gray-700 border border-gray-100"
	}
}

// ─── Photo / File Helpers ──────────────────────────────────────────────────────

// NormalizeFotoProfilURL menormalisasi URL foto profil agar selalu diawali "/static/".
func NormalizeFotoProfilURL(raw string) string {
	url := strings.TrimSpace(raw)
	if url == "" {
		return ""
	}
	switch {
	case strings.HasPrefix(url, "http://"), strings.HasPrefix(url, "https://"):
		return url
	case strings.HasPrefix(url, "/static/"):
		return url
	case strings.HasPrefix(url, "/uploads/"):
		return "/static" + url
	case strings.HasPrefix(url, "uploads/"):
		return "/static/" + url
	case strings.HasPrefix(url, "static/"):
		return "/" + url
	default:
		return url
	}
}

// NormalizeFotoProfilURLPtr adalah versi pointer dari NormalizeFotoProfilURL.
func NormalizeFotoProfilURLPtr(raw *string) *string {
	if raw == nil || *raw == "" {
		return nil
	}
	normalized := NormalizeFotoProfilURL(*raw)
	return &normalized
}

// FotoProfilFilePath mengkonversi URL foto ke path sistem file lokal.
func FotoProfilFilePath(url string) string {
	normalized := NormalizeFotoProfilURL(url)
	if !strings.HasPrefix(normalized, "/static/") {
		return ""
	}
	relativePath := strings.TrimPrefix(normalized, "/static/")
	return filepath.Join("static", filepath.FromSlash(relativePath))
}

// NormalizeArtikelThumbnailURL menormalisasi URL thumbnail artikel.
func NormalizeArtikelThumbnailURL(raw string) string {
	return NormalizeFotoProfilURL(raw)
}

// ─── Flash Message Helpers ─────────────────────────────────────────────────────

// SetFlash menyimpan pesan flash ke cookies (type, title, message).
func SetFlash(c echo.Context, flashType, title, message string) {
	for _, item := range []struct{ name, value string }{
		{"flash_type", flashType},
		{"flash_title", title},
		{"flash_message", message},
	} {
		c.SetCookie(&http.Cookie{
			Name:     item.name,
			Value:    item.value,
			Path:     "/",
			MaxAge:   60,
			HttpOnly: true,
		})
	}
}

// GetFlash membaca dan menghapus pesan flash dari cookies.
// Mengembalikan nil jika tidak ada flash.
func GetFlash(c echo.Context) *FlashMessage {
	flashType, err := c.Cookie("flash_type")
	if err != nil || flashType.Value == "" {
		return nil
	}
	flashTitle, _ := c.Cookie("flash_title")
	flashMessage, _ := c.Cookie("flash_message")

	// Hapus flash cookies
	for _, name := range []string{"flash_type", "flash_title", "flash_message"} {
		c.SetCookie(&http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
	}

	title := ""
	if flashTitle != nil {
		title = flashTitle.Value
	}
	msg := ""
	if flashMessage != nil {
		msg = flashMessage.Value
	}

	return &FlashMessage{
		Type:    flashType.Value,
		Title:   title,
		Message: msg,
	}
}

// ─── User Context Helpers ──────────────────────────────────────────────────────

// GetUserFromContext membaca informasi user yang sudah di-set oleh middleware JWT.
// Mengembalikan UserInfo dengan fallback jika tidak ada.
func GetUserFromContext(c echo.Context) *UserInfo {
	nama, _ := c.Get("user_nama").(string)
	peran, _ := c.Get("user_role_str").(string)
	if nama == "" {
		nama = "Admin"
	}
	if peran == "" {
		peran = "Administrator"
	}
	return &UserInfo{
		NamaLengkap: nama,
		Peran:       peran,
	}
}

// GenerateMonthOptions menghasilkan daftar 12 bulan terakhir untuk dropdown filter.
func GenerateMonthOptions() []MonthOption {
	now := time.Now()
	options := make([]MonthOption, 0, 12)
	for i := 0; i < 12; i++ {
		date := now.AddDate(0, -i, 0)
		value := date.Format("01-2006")
		label := FormatMonthYear(date)
		options = append(options, MonthOption{
			Value:   value,
			Label:   label,
			Current: date.Month() == now.Month() && date.Year() == now.Year(),
		})
	}
	return options
}

// IsSecureRequest mendeteksi apakah request dilakukan melalui HTTPS.
func IsSecureRequest(c echo.Context) bool {
	proto := c.Request().Header.Get("X-Forwarded-Proto")
	return proto == "https" || c.Request().TLS != nil
}

// SetAuthCookie menyimpan JWT token ke cookie HTTP-Only dengan durasi tertentu.
func SetAuthCookie(c echo.Context, token string, duration time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   IsSecureRequest(c),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(duration),
		MaxAge:   int(duration.Seconds()),
	})
}

// ClearAuthCookie menghapus JWT cookie (logout).
func ClearAuthCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   IsSecureRequest(c),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}
