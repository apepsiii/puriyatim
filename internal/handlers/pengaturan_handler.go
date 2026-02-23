package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PengaturanHandler struct{}

func NewPengaturanHandler() *PengaturanHandler {
	return &PengaturanHandler{}
}

type PengaturanData struct {
	PageTitle   string
	User        *UserInfo
	NamaLembaga string
	Deskripsi   string
	Whatsapp    string
	Email       string
	Alamat      string
	Instagram   string
	YouTube     string
}

func (h *PengaturanHandler) Page(c echo.Context) error {
	data := PengaturanData{
		PageTitle:   "Pengaturan Web - Admin Panel",
		User:        &UserInfo{NamaLengkap: "Admin", Peran: "Administrator"},
		NamaLembaga: "Panti Asuhan Puri Yatim",
		Deskripsi:   "Kami adalah lembaga sosial yang berfokus pada pendidikan dan kesejahteraan anak-anak yatim dan dhuafa di wilayah Bogor dan sekitarnya. Berdiri sejak tahun 2010.",
		Whatsapp:    "6281234567890",
		Email:       "info@puriyatim.org",
		Alamat:      "Jl. Pahlawan No. 45, RT 02/RW 05, Kel. Mulyaharja, Kec. Bogor Selatan, Kota Bogor, Jawa Barat 16135",
		Instagram:   "https://instagram.com/puri.yatim",
		YouTube:     "",
	}

	return c.Render(http.StatusOK, "admin/pengaturan.html", data)
}

func (h *PengaturanHandler) Save(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Pengaturan berhasil disimpan",
	})
}
