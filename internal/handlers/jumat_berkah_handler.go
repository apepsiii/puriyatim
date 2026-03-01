package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	ManualCandidates []ManualCandidateItem
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

type ManualCandidateItem struct {
	ID            string
	NamaLengkap   string
	NamaPanggilan string
	Jenjang       string
	StatusAnak    string
	Wilayah       string
	AlreadyReg    bool
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
			initials = GetInitials(namaAnak)
		}

		avatarBg, avatarText := getAvatarStyle(namaAnak)

		regItems = append(regItems, JumatBerkahRegItem{
			ID:          reg.ID,
			AnakAsuhID:  reg.IDAnak,
			NamaAnak:    namaAnak,
			Jenjang:     jenjang,
			StatusAnak:  statusAnak,
			Wilayah:     wilayah,
			SubmittedAt: FormatTimeAgo(reg.WaktuSubmit),
			Status:      reg.StatusApproval.ToServiceStatus(),
			StatusCSS:   ApprovalStatusCSS(reg.StatusApproval),
			Initials:    initials,
			AvatarBg:    avatarBg,
			AvatarText:  avatarText,
		})
	}

	candidates, _ := h.service.GetManualRegistrationCandidates()
	manualCandidates := make([]ManualCandidateItem, 0, len(candidates))
	for _, c := range candidates {
		wilayah := "-"
		if c.RT != "" || c.RW != "" {
			wilayah = fmt.Sprintf("RT %s / RW %s", c.RT, c.RW)
		}
		manualCandidates = append(manualCandidates, ManualCandidateItem{
			ID:            c.ID,
			NamaLengkap:   c.NamaLengkap,
			NamaPanggilan: c.NamaPanggilan,
			Jenjang:       c.Jenjang,
			StatusAnak:    c.StatusAnak,
			Wilayah:       wilayah,
			AlreadyReg:    c.AlreadyReg,
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
		ManualCandidates: manualCandidates,
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

	return JSONOk(c, "Pendaftar berhasil disetujui")
}

func (h *JumatBerkahHandler) Reject(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.RejectRegistration(id); err != nil {
		return JSONBadRequest(c, err.Error())
	}

	return JSONOk(c, "Pendaftar berhasil ditolak")
}

func (h *JumatBerkahHandler) BulkApprove(c echo.Context) error {
	idList := filterEmpty(splitIDs(c.FormValue("ids")))
	count := h.service.ApproveMultiple(idList)
	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": strconv.Itoa(count) + " pendaftar berhasil disetujui",
		"count":   count,
	})
}

func (h *JumatBerkahHandler) BulkReject(c echo.Context) error {
	idList := filterEmpty(splitIDs(c.FormValue("ids")))
	count := h.service.RejectMultiple(idList)
	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": strconv.Itoa(count) + " pendaftar berhasil ditolak",
		"count":   count,
	})
}

func (h *JumatBerkahHandler) ApproveAll(c echo.Context) error {
	count := h.service.ApproveAllPending()
	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": strconv.Itoa(count) + " pendaftar berhasil disetujui",
		"count":   count,
	})
}

func (h *JumatBerkahHandler) UpdateQuota(c echo.Context) error {
	quota, err := strconv.Atoi(c.FormValue("quota"))
	if err != nil || quota < 1 {
		return JSONBadRequest(c, "Kuota tidak valid")
	}

	if err = h.service.UpdateQuota(quota); err != nil {
		return JSONBadRequest(c, err.Error())
	}

	return JSONOk(c, "Kuota berhasil diperbarui")
}

func (h *JumatBerkahHandler) ToggleForm(c echo.Context) error {
	open := c.FormValue("open") == "true"

	if err := h.service.ToggleForm(open); err != nil {
		return JSONBadRequest(c, err.Error())
	}

	status := "ditutup"
	if open {
		status = "dibuka"
	}
	return JSONOk(c, "Form pendaftaran berhasil "+status)
}

func (h *JumatBerkahHandler) ManualRegister(c echo.Context) error {
	ids := c.FormValue("anak_ids")
	autoApprove := c.FormValue("auto_approve") == "true"
	catatan := strings.TrimSpace(c.FormValue("catatan"))

	idList := splitIDs(ids)
	if len(idList) == 0 {
		return JSONBadRequest(c, "Pilih minimal satu anak asuh")
	}

	count := 0
	var regErrors []string
	for _, id := range idList {
		if _, err := h.service.CreateManualRegistration(id, autoApprove, catatan); err != nil {
			regErrors = append(regErrors, fmt.Sprintf("%s: %s", id, err.Error()))
			continue
		}
		count++
	}

	if count == 0 {
		msg := "Tidak ada pendaftaran manual yang berhasil diproses"
		if len(regErrors) > 0 {
			msg = regErrors[0]
		}
		return JSONWithFields(c, map[string]interface{}{
			"success": false,
			"message": msg,
			"errors":  regErrors,
		})
	}

	actionText := "berhasil didaftarkan manual"
	if autoApprove {
		actionText = "berhasil didaftarkan manual dan langsung disetujui"
	}

	return JSONWithFields(c, map[string]interface{}{
		"success": true,
		"message": strconv.Itoa(count) + " anak asuh " + actionText,
		"count":   count,
		"errors":  regErrors,
	})
}

// splitIDs memecah string berisi ID yang dipisah koma menjadi slice.
func splitIDs(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool { return r == ',' })
}

// filterEmpty membuang string kosong dari slice.
func filterEmpty(ss []string) []string {
	result := ss[:0]
	for _, s := range ss {
		if t := strings.TrimSpace(s); t != "" {
			result = append(result, t)
		}
	}
	return result
}
