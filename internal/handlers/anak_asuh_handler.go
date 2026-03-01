package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AnakAsuhHandler struct {
	service            *services.AnakAsuhService
	keuanganService    *services.KeuanganService
	jumatBerkahService *services.JumatBerkahService
	exportImportService *services.ExportImportService
}

func NewAnakAsuhHandler(service *services.AnakAsuhService, keuanganService *services.KeuanganService, jumatBerkahService *services.JumatBerkahService, exportImportService *services.ExportImportService) *AnakAsuhHandler {
	return &AnakAsuhHandler{
		service:            service,
		keuanganService:    keuanganService,
		jumatBerkahService: jumatBerkahService,
		exportImportService: exportImportService,
	}
}

type AnakAsuhListData struct {
	Title        string
	AnakAsuh     []AnakAsuhItem
	TotalCount   int
	StartIndex   int
	EndIndex     int
	CurrentPage  int
	Pages        []int
	HasPrev      bool
	HasNext      bool
	PrevPage     int
	NextPage     int
	Flash        *FlashMessage
	StatusList   []string
	RTList       []string
	RWList       []string
}

type AnakAsuhItem struct {
	ID            string
	NamaLengkap   string
	NamaPanggilan string
	Initials      string
	FotoProfilURL *string
	Domisili      string
	Kelurahan     string
	NamaSekolah   string
	Kelas         string
	Status        string
	StatusBg      string
	StatusText    string
	StatusBorder  string
	StatusDot     string
	AvatarBg      string
	AvatarText    string
	Wilayah       string
	Jenjang       string
	RT            string
	RW            string
}

type AnakAsuhDetailData struct {
	Title              string
	AnakAsuh           AnakAsuhDetailItem
	RiwayatBantuan     []RiwayatBantuanItem
	RiwayatJumatBerkah []RiwayatJumatBerkahItem
	Flash              *FlashMessage
}

type RiwayatJumatBerkahItem struct {
	Tanggal    string
	Status     string
	StatusBg   string
	StatusText string
	Keterangan string
}

type AnakAsuhDetailItem struct {
	ID                    string
	NamaLengkap           string
	NamaPanggilan         string
	Initials              string
	NIK                   string
	TempatLahir           string
	TanggalLahir          string
	TanggalLahirFormatted string
	TanggalMasuk          string
	TanggalMasukFormatted string
	JenisKelamin          string
	AlamatJalan           string
	RT                    string
	RW                    string
	DesaKelurahan         string
	Kecamatan             string
	Kota                  string
	Wilayah               string
	Status                string
	StatusBg              string
	StatusText            string
	StatusBorder          string
	AvatarBg              string
	AvatarText            string
	JenjangPendidikan     string
	NamaSekolah           string
	Kelas                 string
	KondisiKesehatan      string
	CatatanKhusus         string
	NamaWali              string
	HubunganWali          string
	KontakWali            string
	FotoProfilURL         *string
	Usia                  int
	StatusAktif           string
}

type RiwayatBantuanItem struct {
	Tanggal    string
	Jenis      string
	Keterangan string
	Nominal    string
}

type AnakAsuhFormData struct {
	Title    string
	IsEdit   bool
	AnakAsuh *AnakAsuhDetailItem
	Flash    *FlashMessage
}

type FlashMessage struct {
	Type    string
	Title   string
	Message string
}

func (h *AnakAsuhHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	flash := getFlash(c)

	allData, err := h.service.GetAll()
	if err != nil {
		setFlash(c, "error", "Kesalahan Database", "Gagal mengambil data anak asuh: "+err.Error())
		return c.Redirect(http.StatusFound, "/admin/dashboard")
	}

	anakAsuhList := make([]AnakAsuhItem, 0, len(allData))

	for _, item := range allData {
		bg, text, border, dot := item.GetStatusStyle()
		avatarBg, avatarText := item.GetAvatarStyle()
		fotoProfilURL := normalizeFotoProfilURLPtr(item.FotoProfilURL)

		anakAsuhList = append(anakAsuhList, AnakAsuhItem{
			ID:            item.ID,
			NamaLengkap:   item.NamaLengkap,
			NamaPanggilan: item.NamaPanggilan,
			Initials:      item.GetInitials(),
			FotoProfilURL: fotoProfilURL,
			Domisili:      item.GetDomisili(),
			Kelurahan:     item.DesaKelurahan,
			NamaSekolah:   item.NamaSekolah,
			Kelas:         item.Kelas,
			Status:        string(item.StatusAnak),
			StatusBg:      bg,
			StatusText:    text,
			StatusBorder:  border,
			StatusDot:     dot,
			AvatarBg:      avatarBg,
			AvatarText:    avatarText,
			Wilayah:       item.GetWilayah(),
			Jenjang:       item.JenjangPendidikan,
			RT:            item.RT,
			RW:            item.RW,
		})
	}

	// Get filter options from database
	statusList, _ := h.service.GetUniqueStatusAnak()
	rtList, _ := h.service.GetUniqueRT()
	rwList, _ := h.service.GetUniqueRW()

	data := AnakAsuhListData{
		Title:       "Data Anak Asuh",
		AnakAsuh:    anakAsuhList,
		TotalCount:  len(anakAsuhList),
		StartIndex:  1,
		EndIndex:    len(anakAsuhList),
		CurrentPage: page,
		Pages:       []int{1},
		HasPrev:     false,
		HasNext:     false,
		Flash:       flash,
		StatusList:  statusList,
		RTList:      rtList,
		RWList:      rwList,
	}

	return c.Render(http.StatusOK, "admin/anak_asuh_list.html", data)
}

func (h *AnakAsuhHandler) Detail(c echo.Context) error {
	id := c.Param("id")
	flash := getFlash(c)

	item, err := h.service.GetByID(id)
	if err != nil {
		setFlash(c, "error", "Data Tidak Ditemukan", "Data anak asuh dengan ID tersebut tidak ditemukan.")
		return c.Redirect(http.StatusFound, "/admin/anak-asuh")
	}

	bg, text, border, _ := item.GetStatusStyle()
	avatarBg, avatarText := item.GetAvatarStyle()

	usia := 0
	if !item.TanggalLahir.IsZero() {
		usia = int(time.Since(item.TanggalLahir).Hours() / 24 / 365)
	}

	jenisKelamin := "Laki-laki"
	if item.JenisKelamin == "P" {
		jenisKelamin = "Perempuan"
	}

	nik := ""
	if item.NIK != nil {
		nik = *item.NIK
	}

	anakAsuh := AnakAsuhDetailItem{
		ID:                    item.ID,
		NamaLengkap:           item.NamaLengkap,
		NamaPanggilan:         item.NamaPanggilan,
		Initials:              item.GetInitials(),
		NIK:                   nik,
		TempatLahir:           item.TempatLahir,
		TanggalLahir:          item.TanggalLahir.Format("2006-01-02"),
		TanggalLahirFormatted: formatDate(item.TanggalLahir),
		TanggalMasuk:          item.TanggalMasuk.Format("2006-01-02"),
		TanggalMasukFormatted: formatDate(item.TanggalMasuk),
		JenisKelamin:          jenisKelamin,
		AlamatJalan:           item.AlamatJalan,
		RT:                    item.RT,
		RW:                    item.RW,
		DesaKelurahan:         item.DesaKelurahan,
		Kecamatan:             item.Kecamatan,
		Kota:                  item.Kota,
		Wilayah:               item.GetWilayah(),
		Status:                string(item.StatusAnak),
		StatusBg:              bg,
		StatusText:            text,
		StatusBorder:          border,
		AvatarBg:              avatarBg,
		AvatarText:            avatarText,
		JenjangPendidikan:     item.JenjangPendidikan,
		NamaSekolah:           item.NamaSekolah,
		Kelas:                 item.Kelas,
		KondisiKesehatan:      item.KondisiKesehatan,
		CatatanKhusus:         item.CatatanKhusus,
		NamaWali:              item.NamaWali,
		HubunganWali:          item.HubunganWali,
		KontakWali:            item.KontakWali,
		FotoProfilURL:         normalizeFotoProfilURLPtr(item.FotoProfilURL),
		Usia:                  usia,
		StatusAktif:           string(item.StatusAktif),
	}

	// Get Riwayat Bantuan (Pengeluaran for this child)
	riwayatBantuan := []RiwayatBantuanItem{}
	if h.keuanganService != nil {
		pengeluaranList, err := h.keuanganService.GetPengeluaranByAnakID(id)
		if err == nil {
			for _, p := range pengeluaranList {
				riwayatBantuan = append(riwayatBantuan, RiwayatBantuanItem{
					Tanggal:    p.TanggalPengeluaran.Format("02 Jan 2006"),
					Jenis:      "Santunan/Bantuan",
					Keterangan: p.Keterangan,
					Nominal:    fmt.Sprintf("%.0f", p.Nominal),
				})
			}
		}
	}

	// Get Riwayat Jumat Berkah
	riwayatJumatBerkah := []RiwayatJumatBerkahItem{}
	if h.jumatBerkahService != nil {
		pendaftarList, err := h.jumatBerkahService.GetPendaftarByAnakID(id)
		if err == nil {
			for _, p := range pendaftarList {
				statusBg, statusText := "", ""
				switch p.StatusApproval {
				case models.StatusApprovalDisetujui:
					statusBg = "bg-green-100"
					statusText = "text-green-700"
				case models.StatusApprovalDitolak:
					statusBg = "bg-red-100"
					statusText = "text-red-700"
				default:
					statusBg = "bg-yellow-100"
					statusText = "text-yellow-700"
				}

				tanggal := ""
				if p.Kegiatan != nil {
					tanggal = p.Kegiatan.TanggalKegiatan.Format("02 Jan 2006")
				}

				riwayatJumatBerkah = append(riwayatJumatBerkah, RiwayatJumatBerkahItem{
					Tanggal:    tanggal,
					Status:     string(p.StatusApproval),
					StatusBg:   statusBg,
					StatusText: statusText,
					Keterangan: p.Catatan,
				})
			}
		}
	}

	data := AnakAsuhDetailData{
		Title:              anakAsuh.NamaLengkap,
		AnakAsuh:           anakAsuh,
		RiwayatBantuan:     riwayatBantuan,
		RiwayatJumatBerkah: riwayatJumatBerkah,
		Flash:              flash,
	}

	return c.Render(http.StatusOK, "admin/anak_asuh_detail.html", data)
}

func (h *AnakAsuhHandler) Form(c echo.Context) error {
	flash := getFlash(c)
	data := AnakAsuhFormData{
		Title:  "Tambah Anak Asuh",
		IsEdit: false,
		Flash:  flash,
	}
	return c.Render(http.StatusOK, "admin/anak_asuh_form.html", data)
}

func (h *AnakAsuhHandler) EditForm(c echo.Context) error {
	id := c.Param("id")
	flash := getFlash(c)

	item, err := h.service.GetByID(id)
	if err != nil {
		setFlash(c, "error", "Data Tidak Ditemukan", "Data anak asuh dengan ID tersebut tidak ditemukan.")
		return c.Redirect(http.StatusFound, "/admin/anak-asuh")
	}

	nik := ""
	if item.NIK != nil {
		nik = *item.NIK
	}

	anakAsuh := &AnakAsuhDetailItem{
		ID:                item.ID,
		NamaLengkap:       item.NamaLengkap,
		NamaPanggilan:     item.NamaPanggilan,
		NIK:               nik,
		TempatLahir:       item.TempatLahir,
		TanggalLahir:      item.TanggalLahir.Format("2006-01-02"),
		JenisKelamin:      string(item.JenisKelamin),
		AlamatJalan:       item.AlamatJalan,
		RT:                item.RT,
		RW:                item.RW,
		DesaKelurahan:     item.DesaKelurahan,
		Kecamatan:         item.Kecamatan,
		Kota:              item.Kota,
		Status:            string(item.StatusAnak),
		TanggalMasuk:      item.TanggalMasuk.Format("2006-01-02"),
		JenjangPendidikan: item.JenjangPendidikan,
		NamaSekolah:       item.NamaSekolah,
		Kelas:             item.Kelas,
		KondisiKesehatan:  item.KondisiKesehatan,
		CatatanKhusus:     item.CatatanKhusus,
		NamaWali:          item.NamaWali,
		HubunganWali:      item.HubunganWali,
		KontakWali:        item.KontakWali,
		FotoProfilURL:     normalizeFotoProfilURLPtr(item.FotoProfilURL),
		StatusAktif:       string(item.StatusAktif),
	}

	data := AnakAsuhFormData{
		Title:    "Edit Anak Asuh",
		IsEdit:   true,
		AnakAsuh: anakAsuh,
		Flash:    flash,
	}
	return c.Render(http.StatusOK, "admin/anak_asuh_form.html", data)
}

// handleFileUpload processes the uploaded photo file and returns the file URL
func (h *AnakAsuhHandler) handleFileUpload(c echo.Context) (string, error) {
	file, err := c.FormFile("foto_profil")
	if err != nil {
		// No file uploaded, return empty string (not an error)
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}

	// Validate file size (2MB max)
	if file.Size > 2*1024*1024 {
		return "", fmt.Errorf("ukuran file terlalu besar, maksimal 2MB")
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", fmt.Errorf("format file tidak didukung, gunakan JPG atau PNG")
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadPath := filepath.Join("static", "uploads", filename)

	// Create uploads directory if not exists
	uploadsDir := filepath.Join("static", "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", err
	}

	// Create destination file
	dst, err := os.Create(uploadPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return URL path served by Echo static handler.
	return "/static/uploads/" + filename, nil
}

func (h *AnakAsuhHandler) Create(c echo.Context) error {
	tanggalLahir, _ := time.Parse("2006-01-02", c.FormValue("tanggal_lahir"))
	tanggalMasuk, _ := time.Parse("2006-01-02", c.FormValue("tanggal_masuk"))

	// Handle file upload
	fotoProfilURL, err := h.handleFileUpload(c)
	if err != nil {
		setFlash(c, "error", "Gagal Upload Foto", err.Error())
		return c.Redirect(http.StatusFound, "/admin/anak-asuh/tambah")
	}

	anakAsuh := &models.AnakAsuh{
		NamaLengkap:       c.FormValue("nama_lengkap"),
		NamaPanggilan:     c.FormValue("nama_panggilan"),
		TempatLahir:       c.FormValue("tempat_lahir"),
		TanggalLahir:      tanggalLahir,
		JenisKelamin:      models.JenisKelamin(c.FormValue("jenis_kelamin")),
		AlamatJalan:       c.FormValue("alamat_jalan"),
		RT:                c.FormValue("rt"),
		RW:                c.FormValue("rw"),
		DesaKelurahan:     c.FormValue("desa_kelurahan"),
		Kecamatan:         c.FormValue("kecamatan"),
		Kota:              c.FormValue("kota"),
		StatusAnak:        models.StatusAnak(c.FormValue("status_anak")),
		StatusAktif:       models.StatusAktifAktif,
		TanggalMasuk:      tanggalMasuk,
		JenjangPendidikan: c.FormValue("jenjang_pendidikan"),
		NamaSekolah:       c.FormValue("nama_sekolah"),
		Kelas:             c.FormValue("kelas"),
		KondisiKesehatan:  c.FormValue("kondisi_kesehatan"),
		CatatanKhusus:     c.FormValue("catatan_khusus"),
		NamaWali:          c.FormValue("nama_wali"),
		HubunganWali:      c.FormValue("hubungan_wali"),
		KontakWali:        c.FormValue("kontak_wali"),
	}

	// Set foto profil URL if uploaded
	if fotoProfilURL != "" {
		normalized := normalizeFotoProfilURL(fotoProfilURL)
		anakAsuh.FotoProfilURL = &normalized
	}

	nik := c.FormValue("nik")
	if nik != "" {
		anakAsuh.NIK = &nik
	}

	if anakAsuh.NamaLengkap == "" || anakAsuh.NamaPanggilan == "" || anakAsuh.StatusAnak == "" {
		setFlash(c, "error", "Data Tidak Lengkap", "Mohon lengkapi semua field yang wajib diisi.")
		return c.Redirect(http.StatusFound, "/admin/anak-asuh/tambah")
	}

	err = h.service.Create(anakAsuh)
	if err != nil {
		setFlash(c, "error", "Gagal Menyimpan", "Terjadi kesalahan saat menyimpan data: "+err.Error())
		return c.Redirect(http.StatusFound, "/admin/anak-asuh/tambah")
	}

	setFlash(c, "success", "Berhasil Menyimpan", "Data anak asuh "+anakAsuh.NamaLengkap+" berhasil ditambahkan.")
	return c.Redirect(http.StatusFound, "/admin/anak-asuh")
}

func (h *AnakAsuhHandler) Update(c echo.Context) error {
	id := c.Param("id")

	existing, err := h.service.GetByID(id)
	if err != nil {
		setFlash(c, "error", "Data Tidak Ditemukan", "Data anak asuh dengan ID tersebut tidak ditemukan.")
		return c.Redirect(http.StatusFound, "/admin/anak-asuh")
	}

	tanggalLahir, _ := time.Parse("2006-01-02", c.FormValue("tanggal_lahir"))
	tanggalMasuk, _ := time.Parse("2006-01-02", c.FormValue("tanggal_masuk"))

	// Handle file upload
	fotoProfilURL, err := h.handleFileUpload(c)
	if err != nil {
		setFlash(c, "error", "Gagal Upload Foto", err.Error())
		return c.Redirect(http.StatusFound, "/admin/anak-asuh/"+id+"/edit")
	}

	anakAsuh := &models.AnakAsuh{
		ID:                id,
		NamaLengkap:       c.FormValue("nama_lengkap"),
		NamaPanggilan:     c.FormValue("nama_panggilan"),
		TempatLahir:       c.FormValue("tempat_lahir"),
		TanggalLahir:      tanggalLahir,
		JenisKelamin:      models.JenisKelamin(c.FormValue("jenis_kelamin")),
		AlamatJalan:       c.FormValue("alamat_jalan"),
		RT:                c.FormValue("rt"),
		RW:                c.FormValue("rw"),
		DesaKelurahan:     c.FormValue("desa_kelurahan"),
		Kecamatan:         c.FormValue("kecamatan"),
		Kota:              c.FormValue("kota"),
		StatusAnak:        models.StatusAnak(c.FormValue("status_anak")),
		StatusAktif:       existing.StatusAktif,
		TanggalMasuk:      tanggalMasuk,
		JenjangPendidikan: c.FormValue("jenjang_pendidikan"),
		NamaSekolah:       c.FormValue("nama_sekolah"),
		Kelas:             c.FormValue("kelas"),
		KondisiKesehatan:  c.FormValue("kondisi_kesehatan"),
		CatatanKhusus:     c.FormValue("catatan_khusus"),
		NamaWali:          c.FormValue("nama_wali"),
		HubunganWali:      c.FormValue("hubungan_wali"),
		KontakWali:        c.FormValue("kontak_wali"),
		CreatedAt:         existing.CreatedAt,
	}

	// Update foto profil URL if new file uploaded
	if fotoProfilURL != "" {
		// Delete old file if exists
		if existing.FotoProfilURL != nil && *existing.FotoProfilURL != "" {
			oldFilePath := fotoProfilFilePath(*existing.FotoProfilURL)
			if oldFilePath != "" {
				os.Remove(oldFilePath) // Ignore error if file doesn't exist
			}
		}
		normalized := normalizeFotoProfilURL(fotoProfilURL)
		anakAsuh.FotoProfilURL = &normalized
	} else {
		// Keep existing photo if no new upload
		anakAsuh.FotoProfilURL = normalizeFotoProfilURLPtr(existing.FotoProfilURL)
	}

	nik := c.FormValue("nik")
	if nik != "" {
		anakAsuh.NIK = &nik
	}

	err = h.service.Update(anakAsuh)
	if err != nil {
		setFlash(c, "error", "Gagal Mengupdate", "Terjadi kesalahan saat mengupdate data: "+err.Error())
		return c.Redirect(http.StatusFound, "/admin/anak-asuh/"+id+"/edit")
	}

	setFlash(c, "success", "Berhasil Mengupdate", "Data anak asuh "+anakAsuh.NamaLengkap+" berhasil diperbarui.")
	return c.Redirect(http.StatusFound, "/admin/anak-asuh")
}

func (h *AnakAsuhHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	item, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "Data tidak ditemukan",
		})
	}

	nama := item.NamaLengkap
	err = h.service.Delete(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Gagal menghapus data: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Data " + nama + " berhasil dihapus",
	})
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()-1], t.Year())
}

func normalizeFotoProfilURLPtr(raw *string) *string {
	if raw == nil || *raw == "" {
		return nil
	}
	normalized := normalizeFotoProfilURL(*raw)
	return &normalized
}

func normalizeFotoProfilURL(raw string) string {
	url := strings.TrimSpace(raw)
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	if strings.HasPrefix(url, "/static/") {
		return url
	}
	if strings.HasPrefix(url, "/uploads/") {
		return "/static" + url
	}
	if strings.HasPrefix(url, "uploads/") {
		return "/static/" + url
	}
	if strings.HasPrefix(url, "static/") {
		return "/" + url
	}
	return url
}

func fotoProfilFilePath(url string) string {
	normalized := normalizeFotoProfilURL(url)
	if !strings.HasPrefix(normalized, "/static/") {
		return ""
	}
	relativePath := strings.TrimPrefix(normalized, "/static/")
	return filepath.Join("static", filepath.FromSlash(relativePath))
}

func setFlash(c echo.Context, flashType, title, message string) {
	cookie := &http.Cookie{
		Name:     "flash_type",
		Value:    flashType,
		Path:     "/",
		MaxAge:   60,
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	cookie = &http.Cookie{
		Name:     "flash_title",
		Value:    title,
		Path:     "/",
		MaxAge:   60,
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	cookie = &http.Cookie{
		Name:     "flash_message",
		Value:    message,
		Path:     "/",
		MaxAge:   60,
		HttpOnly: true,
	}
	c.SetCookie(cookie)
}

func getFlash(c echo.Context) *FlashMessage {
	flashType, err := c.Cookie("flash_type")
	if err != nil {
		return nil
	}

	flashTitle, err := c.Cookie("flash_title")
	if err != nil {
		return nil
	}

	flashMessage, err := c.Cookie("flash_message")
	if err != nil {
		return nil
	}

	c.SetCookie(&http.Cookie{
		Name:     "flash_type",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	c.SetCookie(&http.Cookie{
		Name:     "flash_title",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	c.SetCookie(&http.Cookie{
		Name:     "flash_message",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	return &FlashMessage{
		Type:    flashType.Value,
		Title:   flashTitle.Value,
		Message: flashMessage.Value,
	}
}

// ExportExcel exports anak asuh data to Excel format
func (h *AnakAsuhHandler) ExportExcel(c echo.Context) error {
	file, err := h.exportImportService.ExportToExcel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Gagal membuat file Excel: " + err.Error(),
		})
	}

	filename := fmt.Sprintf("data_anak_asuh_%s.xlsx", time.Now().Format("20060102_150405"))
	
	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	return file.Write(c.Response().Writer)
}

// ExportCSV exports anak asuh data to CSV format
func (h *AnakAsuhHandler) ExportCSV(c echo.Context) error {
	filename := fmt.Sprintf("data_anak_asuh_%s.csv", time.Now().Format("20060102_150405"))
	
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	err := h.exportImportService.ExportToCSV(c.Response().Writer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Gagal membuat file CSV: " + err.Error(),
		})
	}
	
	return nil
}

// DownloadTemplate downloads import template
func (h *AnakAsuhHandler) DownloadTemplate(c echo.Context) error {
	file, err := h.exportImportService.GetImportTemplate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Gagal membuat template: " + err.Error(),
		})
	}

	filename := "template_import_anak_asuh.xlsx"
	
	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	return file.Write(c.Response().Writer)
}

// ImportData imports anak asuh data from uploaded file
func (h *AnakAsuhHandler) ImportData(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "File tidak ditemukan",
		})
	}

	// Validate file
	if err := h.exportImportService.ValidateImportFile(file); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// Import based on file type
	var successCount int
	var errors []string

	if strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		successCount, errors, err = h.exportImportService.ImportFromCSV(file)
	} else {
		successCount, errors, err = h.exportImportService.ImportFromExcel(file)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Gagal mengimport data: " + err.Error(),
		})
	}

	stats := h.exportImportService.GetImportStats(successCount, errors)
	
	if successCount > 0 {
		message := fmt.Sprintf("Berhasil mengimport %d data", successCount)
		if len(errors) > 0 {
			message += fmt.Sprintf(", %d data gagal", len(errors))
		}
		
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": message,
			"stats":   stats,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"message": "Tidak ada data yang berhasil diimport",
		"stats":   stats,
	})
}
