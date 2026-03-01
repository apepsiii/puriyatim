package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
)

type PengaturanHandler struct {
	service         *services.PengaturanService
	rekeningService *services.RekeningDonasiService
}

func NewPengaturanHandler(service *services.PengaturanService, rekeningService *services.RekeningDonasiService) *PengaturanHandler {
	return &PengaturanHandler{service: service, rekeningService: rekeningService}
}

type PengaturanData struct {
	PageTitle       string
	User            *UserInfo
	NamaLembaga     string
	Deskripsi       string
	Whatsapp        string
	Email           string
	Alamat          string
	Instagram       string
	YouTube         string
	OverlayGaleri   string
	RekeningList    []*models.RekeningDonasi
}

func (h *PengaturanHandler) Page(c echo.Context) error {
	setting, err := h.service.Get()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Gagal memuat pengaturan")
	}

	rekeningList, _ := h.rekeningService.GetAll()
	if rekeningList == nil {
		rekeningList = []*models.RekeningDonasi{}
	}

	data := PengaturanData{
		PageTitle:    "Pengaturan Web - Admin Panel",
		User:         GetUserFromContext(c),
		NamaLembaga:  setting.NamaLembaga,
		Deskripsi:    setting.DeskripsiTentangKami,
		Whatsapp:     setting.NomorWA,
		Email:        setting.EmailLembaga,
		Alamat:       setting.AlamatLengkap,
		RekeningList: rekeningList,
	}
	if setting.LinkInstagram != nil {
		data.Instagram = *setting.LinkInstagram
	}
	if setting.LinkYouTube != nil {
		data.YouTube = *setting.LinkYouTube
	}
	if setting.OverlayGaleriURL != nil {
		data.OverlayGaleri = normalizeStaticURL(*setting.OverlayGaleriURL)
	}

	return c.Render(http.StatusOK, "admin/pengaturan.html", data)
}

func (h *PengaturanHandler) Save(c echo.Context) error {
	namaLembaga := strings.TrimSpace(c.FormValue("nama_lembaga"))
	deskripsi := strings.TrimSpace(c.FormValue("deskripsi"))
	whatsapp := strings.TrimSpace(c.FormValue("whatsapp"))
	email := strings.TrimSpace(c.FormValue("email"))
	alamat := strings.TrimSpace(c.FormValue("alamat"))
	instagram := strings.TrimSpace(c.FormValue("instagram"))
	youtube := strings.TrimSpace(c.FormValue("youtube"))

	setting, err := h.service.Get()
	if err != nil {
		return JSONInternalError(c, "Gagal memuat data pengaturan")
	}

	if namaLembaga != "" {
		setting.NamaLembaga = namaLembaga
	}
	if deskripsi != "" {
		setting.DeskripsiTentangKami = deskripsi
	}
	if whatsapp != "" {
		setting.NomorWA = whatsapp
	}
	setting.EmailLembaga = email
	if alamat != "" {
		setting.AlamatLengkap = alamat
	}

	if instagram != "" {
		setting.LinkInstagram = &instagram
	} else {
		setting.LinkInstagram = nil
	}
	if youtube != "" {
		setting.LinkYouTube = &youtube
	} else {
		setting.LinkYouTube = nil
	}

	overlayFile, err := c.FormFile("overlay_galeri")
	if err == nil && overlayFile != nil && overlayFile.Filename != "" {
		overlayURL, upErr := saveUploadFile(overlayFile, filepath.Join("static", "uploads", "overlays"), "overlay-galeri")
		if upErr != nil {
			return JSONBadRequest(c, upErr.Error())
		}
		setting.OverlayGaleriURL = &overlayURL
	}

	if err := h.service.Save(setting); err != nil {
		return JSONInternalError(c, "Gagal menyimpan pengaturan")
	}

	return JSONOk(c, "Pengaturan berhasil disimpan")
}

func saveUploadFile(fileHeader *multipart.FileHeader, dir string, filenamePrefix string) (string, error) {
	if fileHeader.Size > 5*1024*1024 {
		return "", fmt.Errorf("ukuran file maksimal 5MB")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file upload")
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".png"
	}

	filename := filenamePrefix + ext
	targetPath := filepath.Join(dir, filename)
	dst, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan file upload")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("gagal menulis file upload")
	}

	return "/" + filepath.ToSlash(targetPath), nil
}

func normalizeStaticURL(raw string) string {
	url := strings.TrimSpace(raw)
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, "/static/") {
		return url
	}
	if strings.HasPrefix(url, "static/") {
		return "/" + url
	}
	return url
}

// --- API Rekening Donasi ---

// ListRekening mengembalikan semua rekening dalam format JSON.
func (h *PengaturanHandler) ListRekening(c echo.Context) error {
	list, err := h.rekeningService.GetAll()
	if err != nil {
		return JSONInternalError(c, "Gagal memuat data rekening")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    list,
	})
}

// CreateRekening menambah rekening baru beserta upload logo bank.
// Menerima multipart/form-data: nama_bank, nomor_rekening, atas_nama, logo_bank (file).
func (h *PengaturanHandler) CreateRekening(c echo.Context) error {
	namaBank := strings.TrimSpace(c.FormValue("nama_bank"))
	nomorRekening := strings.TrimSpace(c.FormValue("nomor_rekening"))
	atasNama := strings.TrimSpace(c.FormValue("atas_nama"))

	if namaBank == "" || nomorRekening == "" || atasNama == "" {
		return JSONBadRequest(c, "Nama bank, nomor rekening, dan atas nama wajib diisi")
	}

	// Upload logo (opsional)
	logoURL := ""
	logoFile, err := c.FormFile("logo_bank")
	if err == nil && logoFile != nil && logoFile.Filename != "" {
		dir := filepath.Join("static", "uploads", "banks")
		// Gunakan nama unik agar tidak bentrok antar bank
		safeName := strings.ToLower(strings.ReplaceAll(namaBank, " ", "-"))
		url, upErr := saveUploadFile(logoFile, dir, "bank-"+safeName)
		if upErr != nil {
			return JSONBadRequest(c, upErr.Error())
		}
		logoURL = url
	}

	item := &models.RekeningDonasi{
		NamaBank:      namaBank,
		LogoURL:       logoURL,
		NomorRekening: nomorRekening,
		AtasNama:      atasNama,
		Urutan:        99,
		Aktif:         true,
	}
	if err := h.rekeningService.Create(item); err != nil {
		// Bersihkan file logo jika DB gagal
		if logoURL != "" {
			os.Remove(strings.TrimPrefix(logoURL, "/"))
		}
		return JSONBadRequest(c, err.Error())
	}
	return JSONOk(c, "Rekening berhasil ditambahkan")
}

// DeleteRekening menghapus rekening beserta file logo dari disk.
func (h *PengaturanHandler) DeleteRekening(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return JSONBadRequest(c, "ID rekening tidak valid")
	}

	// Ambil dulu data rekening untuk mendapatkan path logo
	existing, _ := h.rekeningService.GetByID(id)

	if err := h.rekeningService.Delete(id); err != nil {
		return JSONInternalError(c, "Gagal menghapus rekening")
	}

	// Hapus file logo dari disk jika ada
	if existing != nil && existing.LogoURL != "" {
		localPath := strings.TrimPrefix(existing.LogoURL, "/")
		os.Remove(localPath) // abaikan error jika file sudah tidak ada
	}

	return JSONOk(c, "Rekening berhasil dihapus")
}
