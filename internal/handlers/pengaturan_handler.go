package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
)

type PengaturanHandler struct {
	service *services.PengaturanService
}

func NewPengaturanHandler(service *services.PengaturanService) *PengaturanHandler {
	return &PengaturanHandler{service: service}
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
	RekeningBSI     string
	RekeningMandiri string
	NamaRekening    string
}

func (h *PengaturanHandler) Page(c echo.Context) error {
	setting, err := h.service.Get()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Gagal memuat pengaturan")
	}

	data := PengaturanData{
		PageTitle:   "Pengaturan Web - Admin Panel",
		User:        GetUserFromContext(c),
		NamaLembaga: setting.NamaLembaga,
		Deskripsi:   setting.DeskripsiTentangKami,
		Whatsapp:    setting.NomorWA,
		Email:       setting.EmailLembaga,
		Alamat:      setting.AlamatLengkap,
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
	if setting.RekeningBSI != nil {
		data.RekeningBSI = *setting.RekeningBSI
	}
	if setting.RekeningMandiri != nil {
		data.RekeningMandiri = *setting.RekeningMandiri
	}
	if setting.NamaPemilikRekening != nil {
		data.NamaRekening = *setting.NamaPemilikRekening
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
	rekeningBSI := strings.TrimSpace(c.FormValue("rekening_bsi"))
	rekeningMandiri := strings.TrimSpace(c.FormValue("rekening_mandiri"))
	namaRekening := strings.TrimSpace(c.FormValue("nama_rekening"))

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
	if rekeningBSI != "" {
		setting.RekeningBSI = &rekeningBSI
	} else {
		setting.RekeningBSI = nil
	}
	if rekeningMandiri != "" {
		setting.RekeningMandiri = &rekeningMandiri
	} else {
		setting.RekeningMandiri = nil
	}
	if namaRekening != "" {
		setting.NamaPemilikRekening = &namaRekening
	} else {
		setting.NamaPemilikRekening = nil
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
