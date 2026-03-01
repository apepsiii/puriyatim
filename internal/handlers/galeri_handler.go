package handlers

import (
	"bytes"
	"fmt"
	"image"
	imagedraw "image/draw"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
	xdraw "golang.org/x/image/draw"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type GaleriHandler struct {
	galeriService     *services.GaleriService
	pengaturanService *services.PengaturanService
}

func NewGaleriHandler(galeriService *services.GaleriService, pengaturanService *services.PengaturanService) *GaleriHandler {
	return &GaleriHandler{
		galeriService:     galeriService,
		pengaturanService: pengaturanService,
	}
}

type GaleriAdminItem struct {
	ID               string
	Judul            string
	Deskripsi        string
	GambarURL        string
	GambarAsliURL    string
	GambarOverlayURL string
	CreatedAt        string
}

type GaleriAdminData struct {
	PageTitle  string
	User       *UserInfo
	OverlayURL string
	HasOverlay bool
	FotoList   []GaleriAdminItem
}

func (h *GaleriHandler) AdminPage(c echo.Context) error {
	setting, _ := h.pengaturanService.Get()
	overlayURL := ""
	hasOverlay := false
	if setting != nil && setting.OverlayGaleriURL != nil {
		overlayURL = normalizeStaticURL(*setting.OverlayGaleriURL)
		hasOverlay = overlayURL != ""
	}

	list, _ := h.galeriService.ListAll()
	items := make([]GaleriAdminItem, 0, len(list))
	for _, g := range list {
		items = append(items, GaleriAdminItem{
			ID:               g.ID,
			Judul:            g.Judul,
			Deskripsi:        g.Deskripsi,
			GambarURL:        normalizeStaticURL(g.GambarOverlayURL),
			GambarAsliURL:    normalizeStaticURL(g.GambarAsliURL),
			GambarOverlayURL: normalizeStaticURL(g.GambarOverlayURL),
			CreatedAt:        g.CreatedAt.Format("02 Jan 2006 15:04"),
		})
	}

	data := GaleriAdminData{
		PageTitle:  "Galeri Foto - Admin Panel",
		User:       &UserInfo{NamaLengkap: "Admin", Peran: "Administrator"},
		OverlayURL: overlayURL,
		HasOverlay: hasOverlay,
		FotoList:   items,
	}
	return c.Render(http.StatusOK, "admin/galeri_list.html", data)
}

func (h *GaleriHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "ID galeri tidak valid",
		})
	}

	existing, err := h.galeriService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	judul := strings.TrimSpace(c.FormValue("judul"))
	deskripsi := strings.TrimSpace(c.FormValue("deskripsi"))
	if judul == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Judul foto wajib diisi",
		})
	}

	existing.Judul = judul
	existing.Deskripsi = deskripsi

	file, fileErr := c.FormFile("foto")
	if fileErr == nil && file != nil && file.Filename != "" {
		setting, err := h.pengaturanService.Get()
		if err != nil || setting == nil || setting.OverlayGaleriURL == nil || strings.TrimSpace(*setting.OverlayGaleriURL) == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Overlay galeri belum diatur. Upload dahulu di Pengaturan Web.",
			})
		}

		originalURL, originalPath, err := saveGaleriOriginal(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
		}
		overlayPath := strings.TrimPrefix(normalizeStaticURL(*setting.OverlayGaleriURL), "/")
		processedURL, err := applyOverlayAndSave(originalPath, overlayPath)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Gagal memproses overlay: " + err.Error(),
			})
		}

		oldOriginal := existing.GambarAsliURL
		oldOverlay := existing.GambarOverlayURL
		existing.GambarAsliURL = originalURL
		existing.GambarOverlayURL = processedURL

		if err := h.galeriService.Update(existing); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
		}

		deleteFileIfExists(oldOriginal)
		deleteFileIfExists(oldOverlay)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Foto galeri berhasil diperbarui",
		})
	}

	if err := h.galeriService.Update(existing); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Data galeri berhasil diperbarui",
	})
}

func (h *GaleriHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "ID galeri tidak valid",
		})
	}

	existing, err := h.galeriService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := h.galeriService.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	deleteFileIfExists(existing.GambarAsliURL)
	deleteFileIfExists(existing.GambarOverlayURL)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Foto galeri berhasil dihapus",
	})
}

func (h *GaleriHandler) Upload(c echo.Context) error {
	setting, err := h.pengaturanService.Get()
	if err != nil || setting == nil || setting.OverlayGaleriURL == nil || strings.TrimSpace(*setting.OverlayGaleriURL) == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Overlay galeri belum diatur. Upload dahulu di Pengaturan Web.",
		})
	}

	judulPrefix := strings.TrimSpace(c.FormValue("judul"))
	deskripsi := strings.TrimSpace(c.FormValue("deskripsi"))
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Format upload tidak valid",
		})
	}
	files := form.File["foto"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Minimal upload 1 file foto",
		})
	}
	if len(files) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Maksimal 20 foto per upload",
		})
	}

	overlayPath := strings.TrimPrefix(normalizeStaticURL(*setting.OverlayGaleriURL), "/")
	var createdCount int
	var failedFiles []string

	for i, file := range files {
		originalURL, originalPath, err := saveGaleriOriginal(file)
		if err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}

		processedURL, err := applyOverlayAndSave(originalPath, overlayPath)
		if err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}

		judul := judulPrefix
		if judul == "" {
			judul = defaultGaleriTitle(file.Filename)
		} else if len(files) > 1 {
			judul = fmt.Sprintf("%s (%d)", judulPrefix, i+1)
		}

		item := &models.GaleriFoto{
			Judul:            judul,
			Deskripsi:        deskripsi,
			GambarAsliURL:    originalURL,
			GambarOverlayURL: processedURL,
		}
		if err := h.galeriService.Create(item); err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}
		createdCount++
	}

	if createdCount == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Semua file gagal diproses. Cek format gambar dan overlay.",
		})
	}

	message := fmt.Sprintf("%d foto berhasil diupload dan diproses overlay", createdCount)
	if len(failedFiles) > 0 {
		message += fmt.Sprintf(", %d file gagal", len(failedFiles))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": message,
		"failed":  failedFiles,
	})
}

func (h *GaleriHandler) PublicPage(c echo.Context) error {
	list, _ := h.galeriService.ListAll()
	items := make([]map[string]string, 0, len(list))
	for _, g := range list {
		items = append(items, map[string]string{
			"judul":      g.Judul,
			"deskripsi":  g.Deskripsi,
			"gambar_url": normalizeStaticURL(g.GambarOverlayURL),
		})
	}

	data := map[string]interface{}{
		"Title":      "Galeri",
		"ActivePage": "gallery",
		"Year":       time.Now().Year(),
		"Items":      items,
	}
	return c.Render(http.StatusOK, "public/galeri.html", data)
}

func saveGaleriOriginal(file *multipart.FileHeader) (string, string, error) {
	if file.Size > 10*1024*1024 {
		return "", "", fmt.Errorf("ukuran file maksimal 10MB")
	}

	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("gagal membuka file upload")
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Join("static", "uploads", "gallery", "original"), 0755); err != nil {
		return "", "", fmt.Errorf("gagal membuat direktori upload")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join("static", "uploads", "gallery", "original", filename)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", "", fmt.Errorf("gagal menyimpan file upload")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", "", fmt.Errorf("gagal menulis file upload")
	}

	return "/" + filepath.ToSlash(savePath), savePath, nil
}

func applyOverlayAndSave(originalPath, overlayPath string) (string, error) {
	originalFile, err := os.Open(originalPath)
	if err != nil {
		return "", err
	}
	defer originalFile.Close()

	originalImg, _, err := image.Decode(originalFile)
	if err != nil {
		return "", err
	}

	overlayFile, err := os.Open(overlayPath)
	if err != nil {
		return "", err
	}
	defer overlayFile.Close()

	overlayImg, err := png.Decode(overlayFile)
	if err != nil {
		return "", fmt.Errorf("overlay harus format PNG")
	}

	baseBounds := originalImg.Bounds()
	canvas := image.NewRGBA(baseBounds)
	imagedraw.Draw(canvas, baseBounds, originalImg, baseBounds.Min, imagedraw.Src)

	scaledOverlay := image.NewRGBA(baseBounds)
	xdraw.CatmullRom.Scale(scaledOverlay, baseBounds, overlayImg, overlayImg.Bounds(), imagedraw.Over, nil)
	imagedraw.Draw(canvas, baseBounds, scaledOverlay, baseBounds.Min, imagedraw.Over)

	if err := os.MkdirAll(filepath.Join("static", "uploads", "gallery", "processed"), 0755); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%d.png", time.Now().UnixNano())
	savePath := filepath.Join("static", "uploads", "gallery", "processed", filename)

	var buf bytes.Buffer
	if err := png.Encode(&buf, canvas); err != nil {
		return "", err
	}
	if err := os.WriteFile(savePath, buf.Bytes(), 0644); err != nil {
		return "", err
	}

	return "/" + filepath.ToSlash(savePath), nil
}

func defaultGaleriTitle(filename string) string {
	name := strings.TrimSpace(strings.TrimSuffix(filename, filepath.Ext(filename)))
	if name == "" {
		return "Dokumentasi Kegiatan"
	}
	return strings.ReplaceAll(name, "_", " ")
}

func deleteFileIfExists(urlPath string) {
	clean := strings.TrimSpace(urlPath)
	if clean == "" || !strings.HasPrefix(clean, "/static/uploads/") {
		return
	}
	fsPath := strings.TrimPrefix(clean, "/")
	if _, err := os.Stat(fsPath); err == nil {
		_ = os.Remove(fsPath)
	}
}
