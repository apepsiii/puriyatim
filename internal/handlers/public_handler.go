package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
)

type PublicHandler struct {
	jumatBerkahService *services.JumatBerkahService
	anakAsuhService    *services.AnakAsuhService
	artikelService     *services.ArtikelService
	keuanganService    *services.KeuanganService
	pengaturanService  *services.PengaturanService
	donasiMinNominal   int64
}

func NewPublicHandler(jumatBerkahService *services.JumatBerkahService, anakAsuhService *services.AnakAsuhService, artikelService *services.ArtikelService, keuanganService *services.KeuanganService, pengaturanService *services.PengaturanService, donasiMinNominal int64) *PublicHandler {
	return &PublicHandler{
		jumatBerkahService: jumatBerkahService,
		anakAsuhService:    anakAsuhService,
		artikelService:     artikelService,
		keuanganService:    keuanganService,
		pengaturanService:  pengaturanService,
		donasiMinNominal:   donasiMinNominal,
	}
}

func (h *PublicHandler) LandingPage(c echo.Context) error {
	berita, err := h.artikelService.GetPublished(5)
	if err != nil {
		berita = []*models.Artikel{}
	}

	for _, b := range berita {
		if b.GambarThumbnail != nil && *b.GambarThumbnail != "" {
			normalized := NormalizeArtikelThumbnailURL(*b.GambarThumbnail)
			b.GambarThumbnail = &normalized
		}
	}

	data := map[string]interface{}{
		"Title":      "Beranda",
		"ActivePage": "home",
		"Year":       time.Now().Year(),
		"Berita":     berita,
	}

	return c.Render(http.StatusOK, "public/landing.html", data)
}

func (h *PublicHandler) AboutPage(c echo.Context) error {
	data := map[string]interface{}{
		"Title":      "Tentang Kami",
		"ActivePage": "about",
		"Year":       time.Now().Year(),
	}

	return c.Render(http.StatusOK, "public/about.html", data)
}

type JumatBerkahFormData struct {
	Title         string
	ActivePage    string
	Year          int
	CurrentPeriod string
	FormOpen      bool
	Quota         int
	QuotaFilled   int
	Remaining     int
	PenerimaList  []JumatBerkahPenerimaItem
	RWList        []string
	RTByRW        map[string][]string
	AnakByWilayah map[string][]AnakAsuhForForm
}

type AnakAsuhForForm struct {
	ID            string
	NamaLengkap   string
	NamaPanggilan string
	Jenjang       string
	Status        string
	AlreadyReg    bool
}

type JumatBerkahPenerimaItem struct {
	NamaLengkap string
	Initials    string
	Wilayah     string
	Jenjang     string
	StatusAnak  string
}

type DoaItem struct {
	Judul  string `json:"judul"`
	Arab   string `json:"arab"`
	Indo   string `json:"indo"`
	Source string `json:"source"`
}

type doaAPIResponse struct {
	Status int       `json:"status"`
	Data   []DoaItem `json:"data"`
}

type DzikirItem struct {
	Type  string `json:"type"`
	Arab  string `json:"arab"`
	Indo  string `json:"indo"`
	Ulang string `json:"ulang"`
}

type dzikirAPIResponse struct {
	Status int          `json:"status"`
	Data   []DzikirItem `json:"data"`
}

func (h *PublicHandler) JumatBerkahForm(c echo.Context) error {
	kegiatan, _ := h.jumatBerkahService.GetCurrentKegiatan()
	anakList, err := h.anakAsuhService.GetAll()
	if err != nil {
		anakList = []*models.AnakAsuh{}
	}

	currentPeriod := ""
	quota := 20
	formOpen := true
	if kegiatan != nil {
		currentPeriod = kegiatan.TanggalKegiatan.Format("Monday, 02 January 2006")
		quota = kegiatan.KuotaMaksimal
		formOpen = kegiatan.StatusKegiatan == models.StatusKegiatanDibuka
	}

	quotaFilled := h.jumatBerkahService.GetApprovedCount()
	remaining := h.jumatBerkahService.GetRemainingQuota()
	penerimaList := []JumatBerkahPenerimaItem{}

	if kegiatan != nil {
		approvedRegs, err := h.jumatBerkahService.GetPendaftarByStatus(kegiatan.ID, models.StatusApprovalDisetujui)
		if err == nil {
			for _, reg := range approvedRegs {
				if reg.Anak == nil {
					continue
				}

				jenjang := reg.Anak.JenjangPendidikan
				if jenjang == "" {
					jenjang = "Belum Sekolah"
				}

				wilayah := "-"
				if reg.Anak.RT != "" || reg.Anak.RW != "" {
					wilayah = fmt.Sprintf("RT %s / RW %s", reg.Anak.RT, reg.Anak.RW)
				}

				penerimaList = append(penerimaList, JumatBerkahPenerimaItem{
					NamaLengkap: reg.Anak.NamaLengkap,
					Initials:    GetInitials(reg.Anak.NamaLengkap),
					Wilayah:     wilayah,
					Jenjang:     jenjang,
					StatusAnak:  string(reg.Anak.StatusAnak),
				})
			}
		}
	}

	rwSet := make(map[string]bool)
	rtByRWSet := make(map[string]map[string]bool)
	for _, anak := range anakList {
		if anak.RW != "" {
			rwSet[anak.RW] = true
		}
		if anak.RT != "" && anak.RW != "" {
			if rtByRWSet[anak.RW] == nil {
				rtByRWSet[anak.RW] = make(map[string]bool)
			}
			rtByRWSet[anak.RW][anak.RT] = true
		}
	}

	rwList := make([]string, 0, len(rwSet))
	for rw := range rwSet {
		rwList = append(rwList, rw)
	}
	for i := 0; i < len(rwList); i++ {
		for j := i + 1; j < len(rwList); j++ {
			if rwList[i] > rwList[j] {
				rwList[i], rwList[j] = rwList[j], rwList[i]
			}
		}
	}

	rtByRW := make(map[string][]string)
	for rw, rtSet := range rtByRWSet {
		rtList := make([]string, 0, len(rtSet))
		for rt := range rtSet {
			if !strings.HasPrefix(rt, "RT ") {
				rtList = append(rtList, "RT "+rt)
			} else {
				rtList = append(rtList, rt)
			}
		}
		for i := 0; i < len(rtList); i++ {
			for j := i + 1; j < len(rtList); j++ {
				if rtList[i] > rtList[j] {
					rtList[i], rtList[j] = rtList[j], rtList[i]
				}
			}
		}
		rtByRW[rw] = rtList
	}

	anakByWilayah := make(map[string][]AnakAsuhForForm)
	for _, anak := range anakList {
		key := "RW " + anak.RW + "-RT " + strings.TrimPrefix(anak.RT, "RT ")
		jenjang := anak.JenjangPendidikan
		if jenjang == "" {
			jenjang = "Belum Sekolah"
		}
		anakByWilayah[key] = append(anakByWilayah[key], AnakAsuhForForm{
			ID:            anak.ID,
			NamaLengkap:   anak.NamaLengkap,
			NamaPanggilan: anak.NamaPanggilan,
			Jenjang:       jenjang,
			Status:        string(anak.StatusAnak),
			AlreadyReg:    h.jumatBerkahService.IsAnakAsuhRegistered(anak.ID),
		})
	}

	data := JumatBerkahFormData{
		Title:         "Jumat Berkah",
		ActivePage:    "jumat-berkah",
		Year:          time.Now().Year(),
		CurrentPeriod: currentPeriod,
		FormOpen:      formOpen,
		Quota:         quota,
		QuotaFilled:   quotaFilled,
		Remaining:     remaining,
		PenerimaList:  penerimaList,
		RWList:        rwList,
		RTByRW:        rtByRW,
		AnakByWilayah: anakByWilayah,
	}

	return c.Render(http.StatusOK, "public/jumat_berkah_form.html", data)
}

func (h *PublicHandler) SubmitJumatBerkahRegistration(c echo.Context) error {
	anakIDs := c.FormValue("anak_ids")
	if anakIDs == "" {
		return JSONBadRequest(c, "Tidak ada anak asuh yang dipilih")
	}

	ids := strings.Split(anakIDs, ",")
	count := 0
	var lastErr error
	for _, id := range ids {
		if id == "" {
			continue
		}
		anak, err := h.anakAsuhService.GetByID(id)
		if err != nil {
			lastErr = err
			continue
		}
		if h.jumatBerkahService.IsAnakAsuhRegistered(id) {
			continue
		}
		_, err = h.jumatBerkahService.CreateRegistration(
			anak.ID,
			anak.NamaLengkap,
			anak.JenjangPendidikan,
			string(anak.StatusAnak),
			anak.RW,
			anak.RT,
		)
		if err != nil {
			lastErr = err
			continue
		}
		count++
	}

	if count == 0 {
		message := "Tidak ada pendaftaran yang berhasil diproses"
		if lastErr != nil {
			message = lastErr.Error()
		}
		return JSONBadRequest(c, message)
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": strconv.Itoa(count) + " anak asuh berhasil didaftarkan",
		"count":   count,
	})
}

func (h *PublicHandler) ZakatCalculator(c echo.Context) error {
	data := map[string]interface{}{
		"Title":      "Kalkulator Zakat",
		"ActivePage": "zakat",
		"Year":       time.Now().Year(),
	}

	return c.Render(http.StatusOK, "public/zakat_calculator.html", data)
}

func (h *PublicHandler) ProgramDonasiPage(c echo.Context) error {
	data := map[string]interface{}{
		"Title":      "Program Donasi",
		"ActivePage": "program",
		"Year":       time.Now().Year(),
		"MinNominal": h.donasiMinNominal,
		"Programs": []map[string]string{
			{"key": "zakat", "label": "Zakat", "desc": "Penyaluran zakat sesuai asnaf", "icon": "fa-mosque", "color": "emerald"},
			{"key": "infaq", "label": "Infaq", "desc": "Dukungan kebutuhan operasional harian", "icon": "fa-hand-holding-dollar", "color": "blue"},
			{"key": "sedekah", "label": "Sedekah", "desc": "Bantuan cepat untuk kebutuhan anak asuh", "icon": "fa-heart", "color": "rose"},
			{"key": "wakaf", "label": "Wakaf", "desc": "Investasi amal jariyah jangka panjang", "icon": "fa-building-columns", "color": "amber"},
			{"key": "lainnya", "label": "Program Lainnya", "desc": "Fidyah, kafarat, atau donasi khusus", "icon": "fa-box-open", "color": "slate"},
		},
	}
	return c.Render(http.StatusOK, "public/program_donasi.html", data)
}

func (h *PublicHandler) ProgramDonasiConfirmationPage(c echo.Context) error {
	rekeningBSI := "7123456789"
	rekeningMandiri := "1400012345678"
	namaRekening := "Puri Yatim"
	nomorWA := "6281234567890"
	if h.pengaturanService != nil {
		setting, err := h.pengaturanService.Get()
		if err == nil && setting != nil {
			if setting.RekeningBSI != nil && strings.TrimSpace(*setting.RekeningBSI) != "" {
				rekeningBSI = strings.TrimSpace(*setting.RekeningBSI)
			}
			if setting.RekeningMandiri != nil && strings.TrimSpace(*setting.RekeningMandiri) != "" {
				rekeningMandiri = strings.TrimSpace(*setting.RekeningMandiri)
			}
			if setting.NamaPemilikRekening != nil && strings.TrimSpace(*setting.NamaPemilikRekening) != "" {
				namaRekening = strings.TrimSpace(*setting.NamaPemilikRekening)
			}
			if strings.TrimSpace(setting.NomorWA) != "" {
				nomorWA = strings.TrimSpace(setting.NomorWA)
			}
		}
	}

	data := map[string]interface{}{
		"Title":           "Konfirmasi Donasi",
		"ActivePage":      "program",
		"Year":            time.Now().Year(),
		"RekeningBSI":     rekeningBSI,
		"RekeningMandiri": rekeningMandiri,
		"NamaRekening":    namaRekening,
		"NomorWA":         nomorWA,
	}
	return c.Render(http.StatusOK, "public/program_donasi_confirmation.html", data)
}

func (h *PublicHandler) DoaHarianPage(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("q"))
	source := strings.ToLower(strings.TrimSpace(c.QueryParam("source")))
	if source == "" {
		source = "harian"
	}
	page, _ := strconv.Atoi(strings.TrimSpace(c.QueryParam("page")))
	if page < 1 {
		page = 1
	}
	const perPage = 10

	doaList, err := fetchDoaList(source, query)
	if err != nil {
		doaList = []DoaItem{}
	}

	totalItems := len(doaList)
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + perPage - 1) / perPage
	}
	if totalPages == 0 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * perPage
	end := start + perPage
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}
	pageItems := doaList[start:end]

	hasPrev := page > 1
	hasNext := page < totalPages
	prevURL := buildDoaPageURL(source, query, page-1)
	nextURL := buildDoaPageURL(source, query, page+1)

	data := map[string]interface{}{
		"Title":       "Doa Harian",
		"ActivePage":  "doa-harian",
		"Year":        time.Now().Year(),
		"DoaList":     pageItems,
		"Source":      source,
		"Query":       query,
		"Sources":     []string{"harian", "quran", "hadits", "pilihan", "ibadah", "haji", "lainnya"},
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  totalItems,
		"HasPrev":     hasPrev,
		"HasNext":     hasNext,
		"PrevURL":     prevURL,
		"NextURL":     nextURL,
		"StartItem":   start + 1,
		"EndItem":     end,
	}

	return c.Render(http.StatusOK, "public/doa_harian.html", data)
}

func (h *PublicHandler) DzikirPage(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("q"))
	dzikirType := strings.ToLower(strings.TrimSpace(c.QueryParam("type")))
	if dzikirType == "" {
		dzikirType = "pagi"
	}

	dzikirList, err := fetchDzikirList(dzikirType, query)
	if err != nil {
		dzikirList = []DzikirItem{}
	}

	data := map[string]interface{}{
		"Title":      "Dzikir Harian",
		"ActivePage": "dzikir",
		"Year":       time.Now().Year(),
		"DzikirList": dzikirList,
		"Type":       dzikirType,
		"Query":      query,
		"Types":      []string{"all", "pagi", "sore", "solat"},
	}

	return c.Render(http.StatusOK, "public/dzikir.html", data)
}

func (h *PublicHandler) ZakatPayment(c echo.Context) error {
	data := map[string]interface{}{
		"Title":      "Pembayaran Zakat",
		"ActivePage": "zakat",
		"Year":       time.Now().Year(),
		"MinNominal": h.donasiMinNominal,
	}

	return c.Render(http.StatusOK, "public/zakat_payment.html", data)
}

func (h *PublicHandler) NewsList(c echo.Context) error {
	berita, err := h.artikelService.GetPublished(20)
	if err != nil {
		berita = []*models.Artikel{}
	}
	for _, b := range berita {
		if b.GambarThumbnail != nil && *b.GambarThumbnail != "" {
			normalized := NormalizeArtikelThumbnailURL(*b.GambarThumbnail)
			b.GambarThumbnail = &normalized
		}
	}

	data := map[string]interface{}{
		"Title":      "Berita",
		"ActivePage": "news",
		"Year":       time.Now().Year(),
		"Berita":     berita,
	}

	return c.Render(http.StatusOK, "public/news_list.html", data)
}

func (h *PublicHandler) NewsDetail(c echo.Context) error {
	slug := c.Param("id")

	artikel, err := h.artikelService.GetBySlug(slug)
	if err != nil {
		data := map[string]interface{}{
			"Title":      "Artikel Tidak Ditemukan",
			"ActivePage": "news",
			"Year":       time.Now().Year(),
			"Error":      "Artikel yang Anda cari tidak ditemukan",
		}
		return c.Render(http.StatusNotFound, "public/news_list.html", data)
	}
	if artikel.GambarThumbnail != nil && *artikel.GambarThumbnail != "" {
		normalized := NormalizeArtikelThumbnailURL(*artikel.GambarThumbnail)
		artikel.GambarThumbnail = &normalized
	}

	data := map[string]interface{}{
		"Title":      artikel.Judul,
		"ActivePage": "news",
		"Year":       time.Now().Year(),
		"Artikel":    artikel,
	}

	return c.Render(http.StatusOK, "public/news_detail.html", data)
}

func (h *PublicHandler) SubmitJumatBerkah(c echo.Context) error {
	nama := c.FormValue("nama")
	telepon := c.FormValue("telepon")
	email := c.FormValue("email")
	rt := c.FormValue("rt")
	rw := c.FormValue("rw")
	jumlah := c.FormValue("jumlah")
	pengiriman := c.FormValue("pengiriman")
	tanggal := c.FormValue("tanggal")
	catatan := c.FormValue("catatan")

	if nama == "" || telepon == "" || rt == "" || rw == "" || jumlah == "" || pengiriman == "" || tanggal == "" {
		return JSONBadRequest(c, "Semua field wajib diisi")
	}

	jumlahInt, err := strconv.Atoi(jumlah)
	if err != nil || jumlahInt <= 0 {
		return JSONBadRequest(c, "Jumlah tidak valid")
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": "Pendaftaran Jumat Berkah berhasil",
		"data": map[string]interface{}{
			"nama":       nama,
			"telepon":    telepon,
			"email":      email,
			"rt":         rt,
			"rw":         rw,
			"jumlah":     jumlahInt,
			"pengiriman": pengiriman,
			"tanggal":    tanggal,
			"catatan":    catatan,
		},
	})
}

func (h *PublicHandler) SubmitZakatPayment(c echo.Context) error {
	nama := c.FormValue("nama")
	telepon := c.FormValue("telepon")
	email := c.FormValue("email")
	anonymous := c.FormValue("anonymous")
	zakatType := c.FormValue("zakat_type")
	jumlah := c.FormValue("jumlah")
	paymentMethod := c.FormValue("payment_method")
	doa := c.FormValue("doa")
	receipt := c.FormValue("receipt")
	newsletter := c.FormValue("newsletter")

	if nama == "" || telepon == "" || zakatType == "" || jumlah == "" || paymentMethod == "" {
		return JSONBadRequest(c, "Semua field wajib diisi")
	}

	jumlahInt, err := strconv.Atoi(jumlah)
	if err != nil || jumlahInt <= 0 {
		return JSONBadRequest(c, "Jumlah tidak valid")
	}

	catatan := "Zakat via halaman pembayaran"
	if doa != "" {
		catatan = catatan + " | Doa: " + doa
	}
	if paymentMethod != "" {
		catatan = catatan + " | Metode: " + paymentMethod
	}

	if err := h.createPendingDonation(nama, float64(jumlahInt), models.KategoriDanaZakat, catatan); err != nil {
		return JSONBadRequest(c, "Gagal mencatat donasi: "+err.Error())
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": "Pembayaran zakat berhasil diproses",
		"data": map[string]interface{}{
			"nama":           nama,
			"telepon":        telepon,
			"email":          email,
			"anonymous":      anonymous,
			"zakat_type":     zakatType,
			"jumlah":         jumlahInt,
			"payment_method": paymentMethod,
			"doa":            doa,
			"receipt":        receipt,
			"newsletter":     newsletter,
		},
	})
}

func (h *PublicHandler) SubmitProgramDonasi(c echo.Context) error {
	nama := strings.TrimSpace(c.FormValue("nama"))
	if nama == "" {
		nama = "Hamba Allah"
	}

	program := strings.ToLower(strings.TrimSpace(c.FormValue("program")))
	if program == "" {
		return JSONBadRequest(c, "Program donasi wajib dipilih")
	}

	nominalRaw := strings.TrimSpace(c.FormValue("nominal"))
	nominal, err := parseNominalFromForm(nominalRaw)
	if err != nil || nominal <= 0 {
		return JSONBadRequest(c, "Nominal tidak valid")
	}

	kontak := strings.TrimSpace(c.FormValue("kontak"))
	metode := strings.TrimSpace(c.FormValue("metode"))
	catatanUser := strings.TrimSpace(c.FormValue("catatan"))

	kategori := mapProgramToKategori(program)
	catatan := "Program: " + program
	if metode != "" {
		catatan += " | Metode: " + metode
	}
	if kontak != "" {
		catatan += " | Kontak: " + kontak
	}
	if catatanUser != "" {
		catatan += " | Catatan: " + catatanUser
	}

	if err := h.createPendingDonation(nama, nominal, kategori, catatan); err != nil {
		return JSONBadRequest(c, "Gagal mencatat donasi: "+err.Error())
	}

	return JSONOk(c, "Donasi berhasil dicatat dan menunggu verifikasi admin")
}

func (h *PublicHandler) SubmitProgramDonasiConfirmation(c echo.Context) error {
	nama := strings.TrimSpace(c.FormValue("nama"))
	if nama == "" {
		nama = "Hamba Allah"
	}
	nominal, err := parseNominalFromForm(strings.TrimSpace(c.FormValue("nominal")))
	if err != nil || nominal <= 0 {
		return JSONBadRequest(c, "Nominal tidak valid")
	}
	program := strings.ToLower(strings.TrimSpace(c.FormValue("program")))
	if program == "" {
		return JSONBadRequest(c, "Program donasi wajib dipilih")
	}
	nomorHP := strings.TrimSpace(c.FormValue("nomor_hp"))
	if nomorHP == "" {
		return JSONBadRequest(c, "Nomor HP wajib diisi")
	}

	buktiFile, err := c.FormFile("bukti_transfer")
	if err != nil {
		return JSONBadRequest(c, "Bukti transfer wajib diupload")
	}
	buktiURL, err := saveDonasiProofFile(buktiFile)
	if err != nil {
		return JSONBadRequest(c, err.Error())
	}

	metode := strings.TrimSpace(c.FormValue("metode"))
	catatan := "Program: " + program + " | Nomor HP: " + nomorHP
	if metode != "" {
		catatan += " | Metode: " + metode
	}

	kategori := mapProgramToKategori(program)
	pemasukan := &models.PemasukanDonasi{
		NamaDonatur:      nama,
		TanggalDonasi:    time.Now(),
		Nominal:          nominal,
		KategoriDana:     kategori,
		Catatan:          catatan,
		BuktiTransaksi:   buktiURL,
		StatusVerifikasi: models.StatusVerifikasiPending,
	}
	if err := h.keuanganService.CreatePemasukan(pemasukan); err != nil {
		return JSONBadRequest(c, "Gagal mencatat donasi: "+err.Error())
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": "Bukti transfer berhasil dikirim. Status: menunggu verifikasi.",
		"id":      pemasukan.ID,
		"status":  "pending",
		"bukti":   buktiURL,
	})
}

func (h *PublicHandler) GetProgramDonasiStatus(c echo.Context) error {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		return JSONBadRequest(c, "ID donasi tidak valid")
	}
	p, err := h.keuanganService.GetPemasukanByID(id)
	if err != nil {
		return JSONNotFound(c, "Data donasi tidak ditemukan")
	}

	label := "Menunggu Verifikasi"
	if p.StatusVerifikasi == models.StatusVerifikasiVerified {
		label = "Terverifikasi"
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":         p.ID,
			"status":     string(p.StatusVerifikasi),
			"statusText": label,
		},
	})
}

func (h *PublicHandler) GetProgramDonasiHistory(c echo.Context) error {
	nomorHP := strings.TrimSpace(c.QueryParam("nomor_hp"))
	if nomorHP == "" {
		return JSONBadRequest(c, "Nomor HP wajib diisi")
	}

	list, err := h.keuanganService.GetPemasukanByNomorHP(nomorHP)
	if err != nil {
		return JSONInternalError(c, "Gagal mengambil riwayat donasi")
	}

	result := make([]map[string]interface{}, 0)
	for _, p := range list {
		statusText := "Menunggu Verifikasi"
		if p.StatusVerifikasi == models.StatusVerifikasiVerified {
			statusText = "Terverifikasi"
		}
		result = append(result, map[string]interface{}{
			"id":          p.ID,
			"program":     extractCatatanValue(p.Catatan, "Program"),
			"nominal":     p.Nominal,
			"tanggal":     p.TanggalDonasi.Format("02 Jan 2006"),
			"status":      string(p.StatusVerifikasi),
			"status_text": statusText,
		})
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

func (h *PublicHandler) SubscribeNewsletter(c echo.Context) error {
	email := c.FormValue("email")

	if email == "" {
		return JSONBadRequest(c, "Email wajib diisi")
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": "Berhasil berlangganan newsletter",
		"email":   email,
	})
}

func (h *PublicHandler) ZakatSuccess(c echo.Context) error {
	data := map[string]interface{}{
		"Title":      "Pembayaran Berhasil",
		"ActivePage": "zakat",
		"Year":       time.Now().Year(),
	}

	return c.Render(http.StatusOK, "public/zakat_success.html", data)
}

func (h *PublicHandler) GetJumatBerkahData(c echo.Context) error {
	rw := c.QueryParam("rw")
	rt := c.QueryParam("rt")

	anakList, err := h.anakAsuhService.GetAll()
	if err != nil {
		anakList = []*models.AnakAsuh{}
	}
	result := make([]map[string]interface{}, 0)

	for _, anak := range anakList {
		if anak.RW == rw && anak.RT == rt {
			result = append(result, map[string]interface{}{
				"id":             anak.ID,
				"nama_lengkap":   anak.NamaLengkap,
				"nama_panggilan": anak.NamaPanggilan,
				"jenjang":        anak.JenjangPendidikan,
				"status":         string(anak.StatusAnak),
				"already_reg":    h.jumatBerkahService.IsAnakAsuhRegistered(anak.ID),
			})
		}
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

func fetchDoaList(source, query string) ([]DoaItem, error) {
	baseURL := "https://muslim-api-three.vercel.app/v1/doa"
	requestURL := baseURL

	if query != "" {
		requestURL = baseURL + "/find?query=" + url.QueryEscape(query)
	} else {
		requestURL = baseURL + "?source=" + url.QueryEscape(source)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api doa returned status %d", resp.StatusCode)
	}

	var apiResp doaAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	normalizedSource := strings.ToLower(source)
	result := make([]DoaItem, 0, len(apiResp.Data))
	for _, item := range apiResp.Data {
		item.Source = strings.ToLower(strings.TrimSpace(item.Source))
		if normalizedSource != "" && query != "" && item.Source != normalizedSource {
			continue
		}
		result = append(result, item)
	}

	return result, nil
}

func fetchDzikirList(dzikirType, query string) ([]DzikirItem, error) {
	baseURL := "https://muslim-api-three.vercel.app/v1/dzikir"
	requestURL := baseURL
	if dzikirType != "" && dzikirType != "all" {
		requestURL = baseURL + "?type=" + url.QueryEscape(dzikirType)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api dzikir returned status %d", resp.StatusCode)
	}

	var apiResp dzikirAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	q := strings.ToLower(strings.TrimSpace(query))
	result := make([]DzikirItem, 0, len(apiResp.Data))
	for _, item := range apiResp.Data {
		item.Type = strings.ToLower(strings.TrimSpace(item.Type))
		if dzikirType != "" && dzikirType != "all" && item.Type != dzikirType {
			continue
		}
		if q != "" {
			haystack := strings.ToLower(item.Arab + " " + item.Indo + " " + item.Ulang + " " + item.Type)
			if !strings.Contains(haystack, q) {
				continue
			}
		}
		result = append(result, item)
	}

	return result, nil
}

func buildDoaPageURL(source, query string, page int) string {
	if page < 1 {
		page = 1
	}

	values := url.Values{}
	values.Set("source", source)
	if query != "" {
		values.Set("q", query)
	}
	values.Set("page", strconv.Itoa(page))
	return "/doa-harian?" + values.Encode()
}

func (h *PublicHandler) createPendingDonation(nama string, nominal float64, kategori models.KategoriDana, catatan string) error {
	pemasukan := &models.PemasukanDonasi{
		NamaDonatur:      nama,
		TanggalDonasi:    time.Now(),
		Nominal:          nominal,
		KategoriDana:     kategori,
		Catatan:          catatan,
		StatusVerifikasi: models.StatusVerifikasiPending,
	}
	return h.keuanganService.CreatePemasukan(pemasukan)
}

func parseNominalFromForm(raw string) (float64, error) {
	cleaned := strings.ReplaceAll(raw, ".", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	return strconv.ParseFloat(cleaned, 64)
}

func mapProgramToKategori(program string) models.KategoriDana {
	switch program {
	case "zakat":
		return models.KategoriDanaZakat
	case "infaq":
		return models.KategoriDanaInfaq
	case "sedekah":
		return models.KategoriDanaSedekah
	case "wakaf":
		return models.KategoriDanaWakaf
	default:
		return models.KategoriDanaLainnya
	}
}

func saveDonasiProofFile(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file bukti transfer tidak valid")
	}
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("ukuran bukti transfer maksimal 5MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowed[ext] {
		return "", fmt.Errorf("format file harus JPG, PNG, atau WEBP")
	}

	dir := filepath.Join("static", "uploads", "donasi")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file bukti transfer")
	}
	defer src.Close()

	filename := fmt.Sprintf("proof-%d%s", time.Now().UnixNano(), ext)
	targetPath := filepath.Join(dir, filename)
	dst, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan bukti transfer")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("gagal menulis bukti transfer")
	}

	return "/" + filepath.ToSlash(targetPath), nil
}

func extractCatatanValue(catatan, key string) string {
	parts := strings.Split(catatan, "|")
	prefix := strings.ToLower(strings.TrimSpace(key)) + ":"
	for _, part := range parts {
		p := strings.TrimSpace(part)
		lower := strings.ToLower(p)
		if strings.HasPrefix(lower, prefix) {
			return strings.TrimSpace(p[len(prefix):])
		}
	}
	return "-"
}

func (h *PublicHandler) OfflinePage(c echo.Context) error {
	return c.Render(http.StatusOK, "public/offline.html", map[string]interface{}{
		"Title": "Tidak Ada Koneksi",
		"Year":  time.Now().Year(),
	})
}
