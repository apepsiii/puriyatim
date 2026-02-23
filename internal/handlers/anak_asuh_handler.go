package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
)

type AnakAsuhHandler struct {
	service            *services.AnakAsuhService
	keuanganService    *services.KeuanganService
	jumatBerkahService *services.JumatBerkahService
}

func NewAnakAsuhHandler(service *services.AnakAsuhService, keuanganService *services.KeuanganService, jumatBerkahService *services.JumatBerkahService) *AnakAsuhHandler {
	return &AnakAsuhHandler{
		service:            service,
		keuanganService:    keuanganService,
		jumatBerkahService: jumatBerkahService,
	}
}

type AnakAsuhListData struct {
	Title       string
	AnakAsuh    []AnakAsuhItem
	TotalCount  int
	StartIndex  int
	EndIndex    int
	CurrentPage int
	Pages       []int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
	Flash       *FlashMessage
}

type AnakAsuhItem struct {
	ID            string
	NamaLengkap   string
	NamaPanggilan string
	Initials      string
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

		anakAsuhList = append(anakAsuhList, AnakAsuhItem{
			ID:            item.ID,
			NamaLengkap:   item.NamaLengkap,
			NamaPanggilan: item.NamaPanggilan,
			Initials:      item.GetInitials(),
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
		})
	}

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

func (h *AnakAsuhHandler) Create(c echo.Context) error {
	tanggalLahir, _ := time.Parse("2006-01-02", c.FormValue("tanggal_lahir"))
	tanggalMasuk, _ := time.Parse("2006-01-02", c.FormValue("tanggal_masuk"))

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

	nik := c.FormValue("nik")
	if nik != "" {
		anakAsuh.NIK = &nik
	}

	if anakAsuh.NamaLengkap == "" || anakAsuh.NamaPanggilan == "" || anakAsuh.StatusAnak == "" {
		setFlash(c, "error", "Data Tidak Lengkap", "Mohon lengkapi semua field yang wajib diisi.")
		return c.Redirect(http.StatusFound, "/admin/anak-asuh/tambah")
	}

	err := h.service.Create(anakAsuh)
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
