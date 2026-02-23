package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
)

type ArtikelHandler struct {
	service *services.ArtikelService
}

func NewArtikelHandler(service *services.ArtikelService) *ArtikelHandler {
	return &ArtikelHandler{
		service: service,
	}
}

type ArtikelListItem struct {
	ID           string
	Judul        string
	Thumbnail    string
	HasThumbnail bool
	Author       string
	Kategori     string
	Tanggal      string
	Waktu        string
	Status       string
	StatusCSS    string
	Slug         string
	CreatedAt    string
	UpdatedAt    string
}

type ArtikelListData struct {
	PageTitle   string
	User        *UserInfo
	TotalKonten int
	TelahTerbit int
	Draft       int
	ArtikelList []ArtikelListItem
	Success     string
	Error       string
}

func (h *ArtikelHandler) List(c echo.Context) error {
	success := c.QueryParam("success")
	errMsg := c.QueryParam("error")

	artikels, err := h.service.GetAll()
	if err != nil {
		artikels = []*models.Artikel{}
	}

	totalKonten, _ := h.service.Count()
	telahTerbit, _ := h.service.CountByStatus(models.StatusPublikasiTerbit)
	draft, _ := h.service.CountByStatus(models.StatusPublikasiDraft)

	artikelList := make([]ArtikelListItem, 0)
	for _, a := range artikels {
		tanggal := "-"
		waktu := "-"
		if a.TanggalTerbit != nil && !a.TanggalTerbit.IsZero() {
			tanggal = a.TanggalTerbit.Format("02 Jan 2006")
			waktu = a.TanggalTerbit.Format("15:04 WIB")
		} else if !a.CreatedAt.IsZero() {
			tanggal = a.CreatedAt.Format("02 Jan 2006")
			waktu = a.CreatedAt.Format("15:04 WIB")
		}

		thumbnail := ""
		hasThumbnail := false
		if a.GambarThumbnail != nil && *a.GambarThumbnail != "" {
			thumbnail = *a.GambarThumbnail
			hasThumbnail = true
			log.Printf("Artikel %s has thumbnail, length: %d", a.ID, len(thumbnail))
		} else {
			log.Printf("Artikel %s has NO thumbnail", a.ID)
		}

		kategori := "Umum"
		if a.Kategori != nil && a.Kategori.NamaKategori != "" {
			kategori = a.Kategori.NamaKategori
		}

		statusCSS := "bg-gray-100 text-gray-700 border border-gray-200"
		if a.StatusPublikasi == models.StatusPublikasiTerbit {
			statusCSS = "bg-emerald-50 text-emerald-700 border border-emerald-100"
		}

		author := "Admin"
		if a.Pengurus != nil {
			author = a.Pengurus.NamaLengkap
		}

		createdAt := "-"
		if !a.CreatedAt.IsZero() {
			createdAt = a.CreatedAt.Format("02 Jan 2006 15:04")
		}

		updatedAt := "-"
		if !a.UpdatedAt.IsZero() {
			updatedAt = a.UpdatedAt.Format("02 Jan 2006 15:04")
		}

		artikelList = append(artikelList, ArtikelListItem{
			ID:           a.ID,
			Judul:        a.Judul,
			Thumbnail:    thumbnail,
			HasThumbnail: hasThumbnail,
			Author:       author,
			Kategori:     kategori,
			Tanggal:      tanggal,
			Waktu:        waktu,
			Status:       string(a.StatusPublikasi),
			StatusCSS:    statusCSS,
			Slug:         a.Slug,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
	}

	data := ArtikelListData{
		PageTitle:   "Artikel & Kegiatan - Admin Panel",
		User:        &UserInfo{NamaLengkap: "Admin", Peran: "Administrator"},
		TotalKonten: totalKonten,
		TelahTerbit: telahTerbit,
		Draft:       draft,
		ArtikelList: artikelList,
		Success:     success,
		Error:       errMsg,
	}

	return c.Render(http.StatusOK, "admin/artikel_list.html", data)
}

type ArtikelFormData struct {
	PageTitle string
	User      *UserInfo
	Artikel   *models.Artikel
	IsEdit    bool
	Success   string
	Error     string
}

func (h *ArtikelHandler) Form(c echo.Context) error {
	data := ArtikelFormData{
		PageTitle: "Tulis Artikel - Admin Panel",
		User:      &UserInfo{NamaLengkap: "Admin", Peran: "Administrator"},
		IsEdit:    false,
	}

	return c.Render(http.StatusOK, "admin/artikel_form.html", data)
}

func (h *ArtikelHandler) EditForm(c echo.Context) error {
	id := c.Param("id")

	artikel, err := h.service.GetByID(id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/admin/artikel?error=Artikel tidak ditemukan")
	}

	data := ArtikelFormData{
		PageTitle: "Edit Artikel - Admin Panel",
		User:      &UserInfo{NamaLengkap: "Admin", Peran: "Administrator"},
		Artikel:   artikel,
		IsEdit:    true,
	}

	return c.Render(http.StatusOK, "admin/artikel_form.html", data)
}

func (h *ArtikelHandler) Create(c echo.Context) error {
	judul := c.FormValue("judul")
	konten := c.FormValue("konten")
	slug := c.FormValue("slug")
	kategoriID := c.FormValue("kategori_id")
	status := c.FormValue("status")
	metaDeskripsi := c.FormValue("meta_deskripsi")

	if judul == "" {
		return c.Redirect(http.StatusFound, "/admin/artikel?error=Judul artikel wajib diisi")
	}

	if konten == "" {
		konten = "<p>Artikel sedang dalam proses penulisan.</p>"
	}

	if slug == "" {
		slug = generateSlug(judul)
	}

	statusPublikasi := models.StatusPublikasiDraft
	if status == "terbit" || status == "publish" {
		statusPublikasi = models.StatusPublikasiTerbit
	}

	katID := 1
	if kategoriID != "" {
		if id, err := strconv.Atoi(kategoriID); err == nil {
			katID = id
		}
	}

	var thumbnailBase64 *string
	file, err := c.FormFile("gambar_thumbnail")
	if err == nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()
			buf := make([]byte, file.Size)
			src.Read(buf)

			mimeType := file.Header.Get("Content-Type")
			base64Str := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(buf)
			thumbnailBase64 = &base64Str
		}
	}

	artikel := &models.Artikel{
		IDPengurus:      "admin-001",
		Judul:           judul,
		Slug:            slug,
		KontenHTML:      konten,
		IDKategori:      katID,
		StatusPublikasi: statusPublikasi,
		MetaDeskripsi:   metaDeskripsi,
		GambarThumbnail: thumbnailBase64,
	}

	if statusPublikasi == models.StatusPublikasiTerbit {
		now := time.Now()
		artikel.TanggalTerbit = &now
	}

	if err := h.service.Create(artikel); err != nil {
		log.Printf("Error creating artikel: %v", err)
		return c.Redirect(http.StatusFound, "/admin/artikel?error=Gagal menyimpan artikel: "+err.Error())
	}

	return c.Redirect(http.StatusFound, "/admin/artikel?success=Artikel berhasil disimpan")
}

func (h *ArtikelHandler) Update(c echo.Context) error {
	id := c.Param("id")

	judul := c.FormValue("judul")
	konten := c.FormValue("konten")
	slug := c.FormValue("slug")
	kategoriID := c.FormValue("kategori_id")
	status := c.FormValue("status")
	metaDeskripsi := c.FormValue("meta_deskripsi")

	if judul == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Judul artikel wajib diisi",
		})
	}

	artikel, err := h.service.GetByID(id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/admin/artikel?error=Artikel tidak ditemukan")
	}

	if slug == "" {
		slug = generateSlug(judul)
	}

	statusPublikasi := models.StatusPublikasiDraft
	if status == "terbit" || status == "publish" {
		statusPublikasi = models.StatusPublikasiTerbit
	}

	katID := 1
	if kategoriID != "" {
		if kid, err := strconv.Atoi(kategoriID); err == nil {
			katID = kid
		}
	}

	file, err := c.FormFile("gambar_thumbnail")
	if err == nil {
		src, err := file.Open()
		if err == nil {
			defer src.Close()
			buf := make([]byte, file.Size)
			src.Read(buf)

			mimeType := file.Header.Get("Content-Type")
			base64Str := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(buf)
			artikel.GambarThumbnail = &base64Str
		}
	}

	artikel.Judul = judul
	artikel.Slug = slug
	artikel.KontenHTML = konten
	artikel.IDKategori = katID
	artikel.StatusPublikasi = statusPublikasi
	artikel.MetaDeskripsi = metaDeskripsi

	if statusPublikasi == models.StatusPublikasiTerbit && artikel.TanggalTerbit == nil {
		now := time.Now()
		artikel.TanggalTerbit = &now
	}

	if err := h.service.Update(artikel); err != nil {
		return c.Redirect(http.StatusFound, "/admin/artikel?error=Gagal mengupdate artikel")
	}

	return c.Redirect(http.StatusFound, "/admin/artikel?success=Artikel berhasil diupdate")
}

func (h *ArtikelHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Gagal menghapus artikel",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Artikel berhasil dihapus",
	})
}

func (h *ArtikelHandler) Publish(c echo.Context) error {
	id := c.Param("id")

	artikel, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Artikel tidak ditemukan",
		})
	}

	artikel.StatusPublikasi = models.StatusPublikasiTerbit
	now := time.Now()
	artikel.TanggalTerbit = &now

	if err := h.service.Update(artikel); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Gagal mempublikasikan artikel",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Artikel berhasil dipublikasikan",
	})
}

func generateSlug(title string) string {
	slug := make([]byte, 0, len(title))
	for i, c := range title {
		if c >= 'A' && c <= 'Z' {
			slug = append(slug, byte(c+32))
		} else if c >= 'a' && c <= 'z' {
			slug = append(slug, byte(c))
		} else if c >= '0' && c <= '9' {
			slug = append(slug, byte(c))
		} else if c == ' ' || c == '-' || c == '_' {
			if i > 0 && len(slug) > 0 && slug[len(slug)-1] != '-' {
				slug = append(slug, '-')
			}
		} else if strings.Contains("àáâãäåæçèéêëìíîïðñòóôõöøùúûüýþÿ", strings.ToLower(string(c))) {
			replacements := map[rune]byte{
				'ä': 'a', 'æ': 'a', 'á': 'a', 'à': 'a', 'â': 'a', 'ã': 'a', 'å': 'a',
				'ç': 'c',
				'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e',
				'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i',
				'ñ': 'n',
				'ó': 'o', 'ò': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o', 'ø': 'o',
				'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u',
				'ý': 'y', 'ÿ': 'y',
				'Ä': 'a', 'Æ': 'a', 'Á': 'a', 'À': 'a', 'Â': 'a', 'Ã': 'a', 'Å': 'a',
				'Ç': 'c',
				'É': 'e', 'È': 'e', 'Ê': 'e', 'Ë': 'e',
				'Í': 'i', 'Ì': 'i', 'Î': 'i', 'Ï': 'i',
				'Ñ': 'n',
				'Ó': 'o', 'Ò': 'o', 'Ô': 'o', 'Õ': 'o', 'Ö': 'o', 'Ø': 'o',
				'Ú': 'u', 'Ù': 'u', 'Û': 'u', 'Ü': 'u',
				'Ý': 'y',
			}
			if replacement, ok := replacements[c]; ok {
				slug = append(slug, replacement)
			}
		}
	}
	if len(slug) > 0 && slug[len(slug)-1] == '-' {
		slug = slug[:len(slug)-1]
	}
	return string(slug)
}
