package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"

	"github.com/labstack/echo/v4"
)

type JumatBerkahHandler struct {
	service *services.JumatBerkahService
}

func NewJumatBerkahHandler(service *services.JumatBerkahService) *JumatBerkahHandler {
	return &JumatBerkahHandler{
		service: service,
	}
}

type JumatBerkahListData struct {
	Title          string
	PendingCount   int
	ApprovedCount  int
	Quota          int
	QuotaFilled    int
	RemainingQuota int
	CurrentPeriod  string
	FormOpen       bool
	Registrations  []JumatBerkahRegItem
	Flash          *FlashMessage
	StatusFilter   string
}

type JumatBerkahRegItem struct {
	ID          string
	AnakAsuhID  string
	NamaAnak    string
	Jenjang     string
	StatusAnak  string
	Wilayah     string
	SubmittedAt string
	Status      string
	StatusCSS   string
	Initials    string
	AvatarBg    string
	AvatarText  string
}

func (h *JumatBerkahHandler) List(c echo.Context) error {
	flash := getFlash(c)
	statusFilter := c.QueryParam("status")

	kegiatan, _ := h.service.GetCurrentKegiatan()

	var regs []*models.PendaftarJumatBerkah
	var err error

	if kegiatan != nil {
		if statusFilter != "" {
			switch statusFilter {
			case "pending":
				regs, err = h.service.GetPendaftarByStatus(kegiatan.ID, models.StatusApprovalMenunggu)
			case "approved":
				regs, err = h.service.GetPendaftarByStatus(kegiatan.ID, models.StatusApprovalDisetujui)
			case "rejected":
				regs, err = h.service.GetPendaftarByStatus(kegiatan.ID, models.StatusApprovalDitolak)
			default:
				regs, err = h.service.GetPendaftarByKegiatan(kegiatan.ID)
			}
		} else {
			regs, err = h.service.GetPendaftarByKegiatan(kegiatan.ID)
		}
	}

	if err != nil {
		regs = []*models.PendaftarJumatBerkah{}
	}

	regItems := make([]JumatBerkahRegItem, 0, len(regs))

	for _, reg := range regs {
		var statusCSS string
		switch reg.StatusApproval {
		case models.StatusApprovalMenunggu:
			statusCSS = "bg-orange-50 text-orange-700 border border-orange-200"
		case models.StatusApprovalDisetujui:
			statusCSS = "bg-emerald-50 text-emerald-700 border border-emerald-200"
		case models.StatusApprovalDitolak:
			statusCSS = "bg-red-50 text-red-700 border border-red-200"
		}

		namaAnak := ""
		jenjang := ""
		statusAnak := ""
		wilayah := ""
		initials := "??"

		if reg.Anak != nil {
			namaAnak = reg.Anak.NamaLengkap
			jenjang = reg.Anak.JenjangPendidikan
			statusAnak = string(reg.Anak.StatusAnak)
			wilayah = fmt.Sprintf("RT %s / RW %s", reg.Anak.RT, reg.Anak.RW)
			initials = getInitials(namaAnak)
		}

		avatarBg, avatarText := getAvatarStyle(namaAnak)

		regItems = append(regItems, JumatBerkahRegItem{
			ID:          reg.ID,
			AnakAsuhID:  reg.IDAnak,
			NamaAnak:    namaAnak,
			Jenjang:     jenjang,
			StatusAnak:  statusAnak,
			Wilayah:     wilayah,
			SubmittedAt: formatTimeAgo(reg.WaktuSubmit),
			Status:      reg.StatusApproval.ToServiceStatus(),
			StatusCSS:   statusCSS,
			Initials:    initials,
			AvatarBg:    avatarBg,
			AvatarText:  avatarText,
		})
	}

	quota := 20
	quotaFilled := 0
	remaining := 0
	currentPeriod := ""
	formOpen := false

	if kegiatan != nil {
		quota = kegiatan.KuotaMaksimal
		quotaFilled, _ = h.service.CountApprovedByKegiatan(kegiatan.ID)
		remaining = quota - quotaFilled
		if remaining < 0 {
			remaining = 0
		}
		currentPeriod = kegiatan.TanggalKegiatan.Format("Monday, 02 January 2006")
		formOpen = kegiatan.StatusKegiatan == models.StatusKegiatanDibuka
	}

	data := JumatBerkahListData{
		Title:          "Jumat Berkah",
		PendingCount:   h.service.GetPendingCount(),
		ApprovedCount:  h.service.GetApprovedCount(),
		Quota:          quota,
		QuotaFilled:    quotaFilled,
		RemainingQuota: remaining,
		CurrentPeriod:  currentPeriod,
		FormOpen:       formOpen,
		Registrations:  regItems,
		Flash:          flash,
		StatusFilter:   statusFilter,
	}

	return c.Render(http.StatusOK, "admin/jumat_berkah_list.html", data)
}

func (h *JumatBerkahHandler) Approve(c echo.Context) error {
	id := c.Param("id")

	err := h.service.ApproveRegistration(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Pendaftar berhasil disetujui",
	})
}

func (h *JumatBerkahHandler) Reject(c echo.Context) error {
	id := c.Param("id")

	err := h.service.RejectRegistration(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Pendaftar berhasil ditolak",
	})
}

func (h *JumatBerkahHandler) BulkApprove(c echo.Context) error {
	ids := c.FormValue("ids")
	var idList []string
	for _, id := range splitIDs(ids) {
		if id != "" {
			idList = append(idList, id)
		}
	}

	count := h.service.ApproveMultiple(idList)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": strconv.Itoa(count) + " pendaftar berhasil disetujui",
		"count":   count,
	})
}

func (h *JumatBerkahHandler) BulkReject(c echo.Context) error {
	ids := c.FormValue("ids")
	var idList []string
	for _, id := range splitIDs(ids) {
		if id != "" {
			idList = append(idList, id)
		}
	}

	count := h.service.RejectMultiple(idList)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": strconv.Itoa(count) + " pendaftar berhasil ditolak",
		"count":   count,
	})
}

func (h *JumatBerkahHandler) ApproveAll(c echo.Context) error {
	count := h.service.ApproveAllPending()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": strconv.Itoa(count) + " pendaftar berhasil disetujui",
		"count":   count,
	})
}

func (h *JumatBerkahHandler) UpdateQuota(c echo.Context) error {
	quota, err := strconv.Atoi(c.FormValue("quota"))
	if err != nil || quota < 1 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Kuota tidak valid",
		})
	}

	err = h.service.UpdateQuota(quota)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Kuota berhasil diperbarui",
	})
}

func (h *JumatBerkahHandler) ToggleForm(c echo.Context) error {
	open := c.FormValue("open") == "true"

	err := h.service.ToggleForm(open)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	status := "ditutup"
	if open {
		status = "dibuka"
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Form pendaftaran berhasil " + status,
	})
}

func getAvatarStyle(name string) (bg, text string) {
	colors := []struct{ bg, text string }{
		{"bg-emerald-100", "text-emerald-600"},
		{"bg-purple-100", "text-purple-600"},
		{"bg-indigo-100", "text-indigo-600"},
		{"bg-pink-100", "text-pink-600"},
		{"bg-teal-100", "text-teal-600"},
		{"bg-blue-100", "text-blue-600"},
	}
	sum := 0
	for _, c := range name {
		sum += int(c)
	}
	style := colors[sum%len(colors)]
	return style.bg, style.text
}

func splitIDs(s string) []string {
	var ids []string
	current := ""
	for _, r := range s {
		if r == ',' {
			if current != "" {
				ids = append(ids, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		ids = append(ids, current)
	}
	return ids
}
