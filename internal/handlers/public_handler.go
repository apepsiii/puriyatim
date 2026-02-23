package handlers

import (
	"log"
	"net/http"
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
}

func NewPublicHandler(jumatBerkahService *services.JumatBerkahService, anakAsuhService *services.AnakAsuhService, artikelService *services.ArtikelService) *PublicHandler {
	return &PublicHandler{
		jumatBerkahService: jumatBerkahService,
		anakAsuhService:    anakAsuhService,
		artikelService:     artikelService,
	}
}

func (h *PublicHandler) LandingPage(c echo.Context) error {
	berita, err := h.artikelService.GetPublished(5)
	if err != nil {
		log.Printf("Error getting published articles: %v", err)
		berita = []*models.Artikel{}
	}

	log.Printf("Found %d articles", len(berita))
	for i, b := range berita {
		hasThumb := b.GambarThumbnail != nil && *b.GambarThumbnail != ""
		log.Printf("Article %d: %s, hasThumbnail: %v", i, b.Judul, hasThumb)
		if hasThumb {
			log.Printf("  Thumbnail length: %d, starts with: %s", len(*b.GambarThumbnail), (*b.GambarThumbnail)[:50])
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
		RWList:        rwList,
		RTByRW:        rtByRW,
		AnakByWilayah: anakByWilayah,
	}

	return c.Render(http.StatusOK, "public/jumat_berkah_form.html", data)
}

func (h *PublicHandler) SubmitJumatBerkahRegistration(c echo.Context) error {
	anakIDs := c.FormValue("anak_ids")
	if anakIDs == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Tidak ada anak asuh yang dipilih",
		})
	}

	ids := strings.Split(anakIDs, ",")
	count := 0
	for _, id := range ids {
		if id == "" {
			continue
		}
		anak, err := h.anakAsuhService.GetByID(id)
		if err != nil {
			continue
		}
		if h.jumatBerkahService.IsAnakAsuhRegistered(id) {
			continue
		}
		h.jumatBerkahService.CreateRegistration(
			anak.ID,
			anak.NamaLengkap,
			anak.JenjangPendidikan,
			string(anak.StatusAnak),
			anak.RW,
			anak.RT,
		)
		count++
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
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

func (h *PublicHandler) ZakatPayment(c echo.Context) error {
	data := map[string]interface{}{
		"Title":      "Pembayaran Zakat",
		"ActivePage": "zakat",
		"Year":       time.Now().Year(),
	}

	return c.Render(http.StatusOK, "public/zakat_payment.html", data)
}

func (h *PublicHandler) NewsList(c echo.Context) error {
	berita, err := h.artikelService.GetPublished(20)
	if err != nil {
		berita = []*models.Artikel{}
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Semua field wajib diisi",
		})
	}

	jumlahInt, err := strconv.Atoi(jumlah)
	if err != nil || jumlahInt <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Jumlah tidak valid",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Semua field wajib diisi",
		})
	}

	jumlahInt, err := strconv.Atoi(jumlah)
	if err != nil || jumlahInt <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Jumlah tidak valid",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
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

func (h *PublicHandler) SubscribeNewsletter(c echo.Context) error {
	email := c.FormValue("email")

	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email wajib diisi",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}
