package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
)

type KeuanganHandler struct {
	service         *services.KeuanganService
	anakAsuhService *services.AnakAsuhService
}

func NewKeuanganHandler(service *services.KeuanganService, anakAsuhService *services.AnakAsuhService) *KeuanganHandler {
	return &KeuanganHandler{
		service:         service,
		anakAsuhService: anakAsuhService,
	}
}

type KeuanganPageData struct {
	PageTitle             string
	User                  *UserInfo
	TotalSaldo            string
	PemasukanBulan        string
	PengeluaranBulan      string
	PemasukanChange       string
	PengeluaranChange     string
	PemasukanChangeCSS    string
	PengeluaranChangeCSS  string
	PemasukanChangeIcon   string
	PengeluaranChangeIcon string
	CurrentMonth          string
	MonthOptions          []MonthOption
	SelectedMonth         string
	SelectedType          string
	Transactions          []KasTransaction
	TotalPemasukan        string
	TotalPengeluaran      string
	FilteredPemasukan     string
	FilteredPengeluaran   string
	Flash                 *FlashMessage
}

type MonthOption struct {
	Value   string
	Label   string
	Current bool
}

type UserInfo struct {
	NamaLengkap string
	Peran       string
}

type KasTransaction struct {
	ID          string
	Tanggal     string
	Waktu       string
	Deskripsi   string
	Sumber      string
	Kategori    string
	KategoriCSS string
	Jumlah      string
	Type        string
	AnakAsuh    string
	Status      string
	StatusLabel string
	StatusCSS   string
}

func formatPercentChange(change float64) string {
	if change == 0 {
		return "0%"
	}

	sign := ""
	if change > 0 {
		sign = "+"
	}

	return sign + strconv.FormatFloat(change, 'f', 1, 64) + "%"
}

func getChangeCSS(change float64, isPemasukan bool) string {
	if change > 0 {
		if isPemasukan {
			return "text-emerald-600"
		}
		return "text-red-500"
	} else if change < 0 {
		if isPemasukan {
			return "text-red-500"
		}
		return "text-emerald-600"
	}
	return "text-gray-500"
}

func getChangeIcon(change float64) string {
	if change > 0 {
		return "fa-arrow-up"
	} else if change < 0 {
		return "fa-arrow-down"
	}
	return "fa-minus"
}

func formatMonthYear(t time.Time) string {
	months := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return months[int(t.Month())-1] + " " + strconv.Itoa(t.Year())
}

func generateMonthOptions() []MonthOption {
	now := time.Now()
	options := make([]MonthOption, 0, 12)

	months := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

	for i := 0; i < 12; i++ {
		date := now.AddDate(0, -i, 0)
		value := date.Format("01-2006")
		label := months[int(date.Month())-1] + " " + strconv.Itoa(date.Year())
		current := date.Month() == now.Month() && date.Year() == now.Year()

		options = append(options, MonthOption{
			Value:   value,
			Label:   label,
			Current: current,
		})
	}

	return options
}

func (h *KeuanganHandler) BukuKas(c echo.Context) error {
	flash := getFlash(c)

	stats, err := h.service.GetStatistics()
	if err != nil {
		stats = &models.KeuanganStats{}
	}

	selectedMonth := c.QueryParam("month")
	selectedType := c.QueryParam("type")

	var transactions []models.KasTransaction
	var filteredPemasukan, filteredPengeluaran float64

	if selectedMonth != "" {
		transactions, filteredPemasukan, filteredPengeluaran, err = h.service.GetBukuKasFiltered(selectedMonth, selectedType)
		if err != nil {
			transactions = []models.KasTransaction{}
		}
	} else {
		transactions, err = h.service.GetBukuKas()
		filteredPemasukan = stats.TotalPemasukan
		filteredPengeluaran = stats.TotalPengeluaran
		if err != nil {
			transactions = []models.KasTransaction{}
		}
	}

	kasTransactions := make([]KasTransaction, 0, len(transactions))
	for _, t := range transactions {
		var kategoriCSS string
		if t.Type == "masuk" {
			switch t.Kategori {
			case "Zakat":
				kategoriCSS = "bg-purple-50 text-purple-700 border border-purple-100"
			case "Infaq":
				kategoriCSS = "bg-blue-50 text-blue-700 border border-blue-100"
			case "Sedekah":
				kategoriCSS = "bg-emerald-50 text-emerald-700 border border-emerald-100"
			case "Wakaf":
				kategoriCSS = "bg-amber-50 text-amber-700 border border-amber-100"
			default:
				kategoriCSS = "bg-gray-50 text-gray-700 border border-gray-100"
			}
		} else {
			kategoriCSS = "bg-orange-50 text-orange-700 border border-orange-100"
		}

		kasTransactions = append(kasTransactions, KasTransaction{
			ID:          t.ID,
			Tanggal:     t.Tanggal.Format("02 Jan 2006"),
			Waktu:       t.CreatedAt.Format("15:04") + " WIB",
			Deskripsi:   t.Deskripsi,
			Sumber:      t.Donatur,
			Kategori:    t.Kategori,
			KategoriCSS: kategoriCSS,
			Jumlah:      formatRupiah(t.Jumlah),
			Type:        t.Type,
			AnakAsuh:    t.AnakAsuh,
			Status:      string(t.Status),
			StatusLabel: statusLabel(t.Status),
			StatusCSS:   statusCSS(t.Status),
		})
	}

	monthOptions := generateMonthOptions()
	if selectedMonth != "" {
		for i := range monthOptions {
			if monthOptions[i].Value == selectedMonth {
				monthOptions[i].Current = true
			} else {
				monthOptions[i].Current = false
			}
		}
	}

	data := KeuanganPageData{
		PageTitle:             "Buku Kas - Admin Panel",
		User:                  &UserInfo{NamaLengkap: "Admin", Peran: "Administrator"},
		TotalSaldo:            formatRupiah(stats.TotalSaldo),
		PemasukanBulan:        formatRupiah(stats.PemasukanBulanIni),
		PengeluaranBulan:      formatRupiah(stats.PengeluaranBulanIni),
		PemasukanChange:       formatPercentChange(stats.PemasukanChange),
		PengeluaranChange:     formatPercentChange(stats.PengeluaranChange),
		PemasukanChangeCSS:    getChangeCSS(stats.PemasukanChange, true),
		PengeluaranChangeCSS:  getChangeCSS(stats.PengeluaranChange, false),
		PemasukanChangeIcon:   getChangeIcon(stats.PemasukanChange),
		PengeluaranChangeIcon: getChangeIcon(stats.PengeluaranChange),
		CurrentMonth:          formatMonthYear(time.Now()),
		MonthOptions:          monthOptions,
		SelectedMonth:         selectedMonth,
		SelectedType:          selectedType,
		Transactions:          kasTransactions,
		TotalPemasukan:        formatRupiah(stats.TotalPemasukan),
		TotalPengeluaran:      formatRupiah(stats.TotalPengeluaran),
		FilteredPemasukan:     formatRupiah(filteredPemasukan),
		FilteredPengeluaran:   formatRupiah(filteredPengeluaran),
		Flash:                 flash,
	}

	return c.Render(http.StatusOK, "admin/keuangan_buku_kas.html", data)
}

type PemasukanFormData struct {
	PageTitle   string
	User        *UserInfo
	Today       string
	DonaturList []DonaturOption
	Flash       *FlashMessage
}

type DonaturOption struct {
	ID          string
	NamaDonatur string
	Catatan     string
}

func (h *KeuanganHandler) CatatPemasukan(c echo.Context) error {
	flash := getFlash(c)

	donaturList, err := h.service.GetAllDonatur()
	if err != nil {
		donaturList = []models.Donatur{}
	}

	donaturOptions := make([]DonaturOption, 0, len(donaturList))
	for _, d := range donaturList {
		label := d.NamaDonatur
		if d.CatatanKhusus != "" {
			label = d.NamaDonatur + " (" + d.CatatanKhusus + ")"
		}
		donaturOptions = append(donaturOptions, DonaturOption{
			ID:          d.ID,
			NamaDonatur: label,
			Catatan:     d.CatatanKhusus,
		})
	}

	data := PemasukanFormData{
		PageTitle:   "Catat Pemasukan - Admin Panel",
		User:        &UserInfo{NamaLengkap: "Admin", Peran: "Administrator"},
		Today:       time.Now().Format("2006-01-02"),
		DonaturList: donaturOptions,
		Flash:       flash,
	}

	return c.Render(http.StatusOK, "admin/keuangan_pemasukan.html", data)
}

func (h *KeuanganHandler) SavePemasukan(c echo.Context) error {
	tanggalStr := c.FormValue("tanggal")
	nominalStr := c.FormValue("nominal")
	kategori := c.FormValue("kategori")
	namaDonatur := c.FormValue("nama_donatur")
	catatan := c.FormValue("keterangan")

	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Format tanggal tidak valid",
		})
	}

	nominalStr = strings.ReplaceAll(nominalStr, ".", "")
	nominalStr = strings.ReplaceAll(nominalStr, ",", "")
	nominal, err := strconv.ParseFloat(nominalStr, 64)
	if err != nil || nominal <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Nominal tidak valid",
		})
	}

	if namaDonatur == "" {
		namaDonatur = "Hamba Allah"
	}

	pemasukan := &models.PemasukanDonasi{
		NamaDonatur:      namaDonatur,
		TanggalDonasi:    tanggal,
		Nominal:          nominal,
		KategoriDana:     models.KategoriDana(kategori),
		Catatan:          catatan,
		StatusVerifikasi: models.StatusVerifikasiPending,
	}

	err = h.service.CreatePemasukan(pemasukan)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "failed to create pemasukan: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Pemasukan berhasil dicatat (status pending, menunggu verifikasi admin)",
		"id":      pemasukan.ID,
	})
}

type PengeluaranFormData struct {
	PageTitle    string
	User         *UserInfo
	Today        string
	AnakAsuhList []AnakAsuhOption
	Flash        *FlashMessage
}

type AnakAsuhOption struct {
	ID          string
	NamaLengkap string
}

func (h *KeuanganHandler) CatatPengeluaran(c echo.Context) error {
	flash := getFlash(c)

	anakList, _ := h.anakAsuhService.GetAll()
	anakOptions := make([]AnakAsuhOption, 0, len(anakList))
	for _, a := range anakList {
		anakOptions = append(anakOptions, AnakAsuhOption{
			ID:          a.ID,
			NamaLengkap: a.NamaLengkap,
		})
	}

	data := PengeluaranFormData{
		PageTitle:    "Catat Pengeluaran - Admin Panel",
		User:         &UserInfo{NamaLengkap: "Admin", Peran: "Administrator"},
		Today:        time.Now().Format("2006-01-02"),
		AnakAsuhList: anakOptions,
		Flash:        flash,
	}

	return c.Render(http.StatusOK, "admin/keuangan_pengeluaran.html", data)
}

func (h *KeuanganHandler) SavePengeluaran(c echo.Context) error {
	tanggalStr := c.FormValue("tanggal")
	nominalStr := c.FormValue("nominal")
	keterangan := c.FormValue("keterangan")
	idAnak := c.FormValue("id_anak")

	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Format tanggal tidak valid",
		})
	}

	nominalStr = strings.ReplaceAll(nominalStr, ".", "")
	nominalStr = strings.ReplaceAll(nominalStr, ",", "")
	nominal, err := strconv.ParseFloat(nominalStr, 64)
	if err != nil || nominal <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Nominal tidak valid",
		})
	}

	if keterangan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Keterangan wajib diisi",
		})
	}

	pengeluaran := &models.Pengeluaran{
		TanggalPengeluaran: tanggal,
		Nominal:            nominal,
		Keterangan:         keterangan,
	}

	if idAnak != "" {
		pengeluaran.IDAnak = &idAnak
	}

	err = h.service.CreatePengeluaran(pengeluaran)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Pengeluaran berhasil dicatat",
		"id":      pengeluaran.ID,
	})
}

func (h *KeuanganHandler) CreateDonatur(c echo.Context) error {
	nama := c.FormValue("nama")
	telepon := c.FormValue("telepon")
	alamat := c.FormValue("alamat")

	if nama == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Nama wajib diisi",
		})
	}

	if telepon == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "No. telepon wajib diisi",
		})
	}

	donatur := &models.Donatur{
		NamaDonatur: nama,
		NoTelepon:   telepon,
		TipeDonatur: models.TipeDonaturIndividu,
	}

	if alamat != "" {
		donatur.Alamat = &alamat
	}

	err := h.service.CreateDonatur(donatur)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Donatur berhasil ditambahkan",
		"id":      donatur.ID,
		"nama":    donatur.NamaDonatur,
	})
}

func (h *KeuanganHandler) GetTransactionDetail(c echo.Context) error {
	id := c.Param("id")
	tipe := c.QueryParam("type")

	if tipe == "masuk" {
		pemasukan, err := h.service.GetPemasukanByID(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "Transaksi tidak ditemukan",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":           pemasukan.ID,
				"type":         "masuk",
				"tanggal":      pemasukan.TanggalDonasi.Format("02 January 2006"),
				"tanggal_raw":  pemasukan.TanggalDonasi.Format("2006-01-02"),
				"waktu":        pemasukan.CreatedAt.Format("15:04") + " WIB",
				"nominal":      formatRupiah(pemasukan.Nominal),
				"kategori":     string(pemasukan.KategoriDana),
				"sumber":       pemasukan.NamaDonatur,
				"deskripsi":    "Donasi " + string(pemasukan.KategoriDana),
				"catatan":      pemasukan.Catatan,
				"bukti":        pemasukan.BuktiTransaksi,
				"status":       string(pemasukan.StatusVerifikasi),
				"status_label": statusLabel(pemasukan.StatusVerifikasi),
				"created_at":   pemasukan.CreatedAt.Format("02 January 2006 15:04"),
			},
		})
	} else {
		pengeluaran, err := h.service.GetPengeluaranByID(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "Transaksi tidak ditemukan",
			})
		}

		anakNama := ""
		if pengeluaran.Anak != nil {
			anakNama = pengeluaran.Anak.NamaLengkap
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":          pengeluaran.ID,
				"type":        "keluar",
				"tanggal":     pengeluaran.TanggalPengeluaran.Format("02 January 2006"),
				"tanggal_raw": pengeluaran.TanggalPengeluaran.Format("2006-01-02"),
				"waktu":       pengeluaran.CreatedAt.Format("15:04") + " WIB",
				"nominal":     formatRupiah(pengeluaran.Nominal),
				"kategori":    "Pengeluaran",
				"sumber":      anakNama,
				"deskripsi":   pengeluaran.Keterangan,
				"catatan":     "",
				"bukti":       pengeluaran.BuktiTransaksi,
				"created_at":  pengeluaran.CreatedAt.Format("02 January 2006 15:04"),
			},
		})
	}
}

func (h *KeuanganHandler) DeleteTransaction(c echo.Context) error {
	id := c.Param("id")
	tipe := c.QueryParam("type")

	var err error
	if tipe == "masuk" {
		err = h.service.DeletePemasukan(id)
	} else {
		err = h.service.DeletePengeluaran(id)
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Gagal menghapus transaksi: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Transaksi berhasil dihapus",
	})
}

func (h *KeuanganHandler) UpdateTransaction(c echo.Context) error {
	id := c.Param("id")
	tipe := c.QueryParam("type")

	tanggalStr := c.FormValue("tanggal")
	nominalStr := c.FormValue("nominal")

	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Format tanggal tidak valid",
		})
	}

	nominalStr = strings.ReplaceAll(nominalStr, ".", "")
	nominalStr = strings.ReplaceAll(nominalStr, ",", "")
	nominal, err := strconv.ParseFloat(nominalStr, 64)
	if err != nil || nominal <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Nominal tidak valid",
		})
	}

	if tipe == "masuk" {
		kategori := c.FormValue("kategori")
		namaDonatur := c.FormValue("nama_donatur")
		catatan := c.FormValue("catatan")

		if namaDonatur == "" {
			namaDonatur = "Hamba Allah"
		}

		pemasukan := &models.PemasukanDonasi{
			ID:               id,
			NamaDonatur:      namaDonatur,
			TanggalDonasi:    tanggal,
			Nominal:          nominal,
			KategoriDana:     models.KategoriDana(kategori),
			Catatan:          catatan,
			StatusVerifikasi: "",
		}

		err = h.service.UpdatePemasukan(pemasukan)
	} else {
		keterangan := c.FormValue("keterangan")
		idAnak := c.FormValue("id_anak")

		if keterangan == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Keterangan wajib diisi",
			})
		}

		pengeluaran := &models.Pengeluaran{
			ID:                 id,
			TanggalPengeluaran: tanggal,
			Nominal:            nominal,
			Keterangan:         keterangan,
		}

		if idAnak != "" {
			pengeluaran.IDAnak = &idAnak
		}

		err = h.service.UpdatePengeluaran(pengeluaran)
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Transaksi berhasil diperbarui",
	})
}

func (h *KeuanganHandler) VerifyPemasukan(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "ID pemasukan tidak valid",
		})
	}

	if err := h.service.VerifyPemasukan(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Pemasukan berhasil diverifikasi",
	})
}

func statusLabel(status models.StatusVerifikasiPemasukan) string {
	switch status {
	case models.StatusVerifikasiPending:
		return "Pending"
	case models.StatusVerifikasiVerified:
		return "Verified"
	default:
		return "Verified"
	}
}

func statusCSS(status models.StatusVerifikasiPemasukan) string {
	switch status {
	case models.StatusVerifikasiPending:
		return "bg-amber-50 text-amber-700 border border-amber-100"
	case models.StatusVerifikasiVerified:
		return "bg-emerald-50 text-emerald-700 border border-emerald-100"
	default:
		return "bg-emerald-50 text-emerald-700 border border-emerald-100"
	}
}

func (h *KeuanganHandler) GetEditFormData(c echo.Context) error {
	anakList, _ := h.anakAsuhService.GetAll()
	anakOptions := make([]map[string]string, 0, len(anakList))
	for _, a := range anakList {
		anakOptions = append(anakOptions, map[string]string{
			"id":   a.ID,
			"nama": a.NamaLengkap,
		})
	}

	donaturList, _ := h.service.GetAllDonatur()
	donaturOptions := make([]map[string]string, 0, len(donaturList))
	for _, d := range donaturList {
		donaturOptions = append(donaturOptions, map[string]string{
			"id":   d.ID,
			"nama": d.NamaDonatur,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":      true,
		"anak_list":    anakOptions,
		"donatur_list": donaturOptions,
	})
}

func (h *KeuanganHandler) ExportCSV(c echo.Context) error {
	month := c.QueryParam("month")
	transType := c.QueryParam("type")

	var transactions []models.KasTransaction
	if month != "" {
		transactions, _, _, _ = h.service.GetBukuKasFiltered(month, transType)
	} else {
		transactions, _ = h.service.GetBukuKas()
	}

	c.Response().Header().Set("Content-Type", "text/csv; charset=utf-8")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=buku_kas_"+time.Now().Format("20060102150405")+".csv")

	writer := c.Response().Writer
	writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer.Write([]byte("Tanggal,Waktu,Kategori,Keterangan,Sumber,Nominal,Tipe\n"))

	for _, t := range transactions {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
			t.Tanggal.Format("02/01/2006"),
			t.Tanggal.Format("15:04"),
			t.Kategori,
			strings.ReplaceAll(t.Deskripsi, ",", ";"),
			t.Donatur+t.AnakAsuh,
			fmt.Sprintf("%.0f", t.Jumlah),
			t.Type,
		)
		writer.Write([]byte(line))
	}

	return nil
}

func (h *KeuanganHandler) ExportPDF(c echo.Context) error {
	month := c.QueryParam("month")
	transType := c.QueryParam("type")

	var transactions []models.KasTransaction
	var totalPemasukan, totalPengeluaran float64

	if month != "" {
		transactions, totalPemasukan, totalPengeluaran, _ = h.service.GetBukuKasFiltered(month, transType)
	} else {
		transactions, _ = h.service.GetBukuKas()
		stats, _ := h.service.GetStatistics()
		totalPemasukan = stats.TotalPemasukan
		totalPengeluaran = stats.TotalPengeluaran
	}

	data := map[string]interface{}{
		"Title":            "Laporan Buku Kas",
		"Transactions":     transactions,
		"TotalPemasukan":   totalPemasukan,
		"TotalPengeluaran": totalPengeluaran,
		"Saldo":            totalPemasukan - totalPengeluaran,
		"PrintDate":        time.Now().Format("02 January 2006 15:04"),
		"Month":            month,
	}

	return c.Render(http.StatusOK, "admin/keuangan_print.html", data)
}
