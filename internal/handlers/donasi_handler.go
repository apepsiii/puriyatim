package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"
	"puriyatim-app/pkg/pakasir"

	"github.com/labstack/echo/v4"
)

type DonasiHandler struct {
	donasiService *services.DonasiOnlineService
}

func NewDonasiHandler(donasiService *services.DonasiOnlineService) *DonasiHandler {
	return &DonasiHandler{donasiService: donasiService}
}

// ---- POST /api/donasi/create ----
// Body JSON:
//
//	{
//	  "jenis": "donasi" | "zakat" | "jumat_berkah",
//	  "nama_donatur": "...",
//	  "nominal": 50000,
//	  "payment_method": "qris" | "bri_va" | ...,
//	  "catatan": "..."
//	}
func (h *DonasiHandler) CreateTransaction(c echo.Context) error {
	if !h.donasiService.IsPakasirConfigured() {
		return JSONInternalError(c, "Payment gateway belum dikonfigurasi. Hubungi admin.")
	}

	type reqBody struct {
		Jenis         string  `json:"jenis"`
		NamaDonatur   string  `json:"nama_donatur"`
		Nominal       float64 `json:"nominal"`
		PaymentMethod string  `json:"payment_method"`
		Catatan       string  `json:"catatan"`
	}

	var body reqBody
	if err := c.Bind(&body); err != nil {
		return JSONBadRequest(c, "Format request tidak valid")
	}

	body.Jenis = strings.TrimSpace(body.Jenis)
	body.PaymentMethod = strings.TrimSpace(body.PaymentMethod)

	// Validasi jenis
	switch body.Jenis {
	case "donasi", "zakat", "jumat_berkah":
	default:
		return JSONBadRequest(c, "Jenis donasi tidak valid. Gunakan: donasi, zakat, atau jumat_berkah")
	}

	// Validasi payment method
	if _, ok := pakasir.PaymentMethodLabels[body.PaymentMethod]; !ok {
		return JSONBadRequest(c, "Metode pembayaran tidak valid")
	}

	if body.Nominal <= 0 {
		return JSONBadRequest(c, "Nominal harus lebih dari 0")
	}

	req := services.CreateTransactionRequest{
		Jenis:         models.JenisDonasiOnline(body.Jenis),
		NamaDonatur:   body.NamaDonatur,
		Nominal:       body.Nominal,
		PaymentMethod: body.PaymentMethod,
		Catatan:       body.Catatan,
	}

	donasi, err := h.donasiService.CreateTransaction(req)
	if err != nil {
		log.Printf("donasi create error: %v", err)
		return JSONInternalError(c, "Gagal membuat transaksi: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"order_id": donasi.OrderID,
		"data": map[string]interface{}{
			"order_id":       donasi.OrderID,
			"jenis":          string(donasi.Jenis),
			"nama_donatur":   donasi.NamaDonatur,
			"nominal":        donasi.Nominal,
			"total_payment":  donasi.TotalPayment,
			"fee":            donasi.Fee,
			"payment_method": donasi.PaymentMethod,
			"qr_string":      donasi.QRString,
			"va_number":      donasi.VANumber,
			"expired_at":     donasi.ExpiredAt,
			"status":         string(donasi.Status),
		},
	})
}

// ---- GET /api/donasi/status/:order_id ----
// Polling dari frontend untuk cek status pembayaran
func (h *DonasiHandler) GetStatus(c echo.Context) error {
	orderID := strings.TrimSpace(c.Param("order_id"))
	if orderID == "" {
		return JSONBadRequest(c, "order_id tidak valid")
	}

	donasi, err := h.donasiService.CheckAndUpdateStatus(orderID)
	if err != nil {
		return JSONNotFound(c, "Transaksi tidak ditemukan")
	}

	remaining := 0
	if donasi.ExpiredAt != nil && donasi.Status == models.StatusDonasiPending {
		remaining = int(time.Until(*donasi.ExpiredAt).Seconds())
		if remaining < 0 {
			remaining = 0
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"order_id":         donasi.OrderID,
			"status":           string(donasi.Status),
			"payment_method":   donasi.PaymentMethod,
			"nominal":          donasi.Nominal,
			"total_payment":    donasi.TotalPayment,
			"remaining_seconds": remaining,
			"completed_at":     donasi.CompletedAt,
		},
	})
}

// ---- POST /api/pakasir/webhook ----
// Endpoint yang didaftarkan di dashboard Pakasir sebagai Webhook URL
func (h *DonasiHandler) Webhook(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("pakasir webhook: gagal baca body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	var payload pakasir.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("pakasir webhook: gagal decode JSON: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	log.Printf("pakasir webhook received: order=%s status=%s amount=%d method=%s",
		payload.OrderID, payload.Status, payload.Amount, payload.PaymentMethod)

	if err := h.donasiService.HandleWebhook(&payload); err != nil {
		log.Printf("pakasir webhook error: %v", err)
		// Return 200 agar Pakasir tidak retry terus (idempotent handling di dalam service)
		return c.JSON(http.StatusOK, map[string]string{"status": "error_logged"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// ---- GET /donasi/bayar/:order_id ----
// Halaman tunggu pembayaran — tampilkan QR / VA + timer countdown
func (h *DonasiHandler) PaymentPage(c echo.Context) error {
	orderID := strings.TrimSpace(c.Param("order_id"))
	if orderID == "" {
		return c.Redirect(http.StatusFound, "/program-donasi")
	}

	donasi, err := h.donasiService.GetByOrderID(orderID)
	if err != nil {
		return c.Render(http.StatusNotFound, "public/404.html", map[string]interface{}{
			"Year": time.Now().Year(),
		})
	}

	// Redirect kalau sudah selesai
	if donasi.Status == models.StatusDonasiCompleted {
		return c.Redirect(http.StatusFound, "/donasi/sukses/"+orderID)
	}

	methodLabel := pakasir.PaymentMethodLabels[donasi.PaymentMethod]
	if methodLabel == "" {
		methodLabel = donasi.PaymentMethod
	}

	data := map[string]interface{}{
		"Title":         "Pembayaran",
		"ActivePage":    "program",
		"Year":          time.Now().Year(),
		"Donasi":        donasi,
		"MethodLabel":   methodLabel,
		"IsQRIS":        donasi.IsQRIS(),
		"IsVA":          donasi.IsVA(),
	}

	return c.Render(http.StatusOK, "public/donasi_payment.html", data)
}

// ---- GET /donasi/sukses/:order_id ----
// Halaman sukses setelah pembayaran berhasil
func (h *DonasiHandler) SuccessPage(c echo.Context) error {
	orderID := strings.TrimSpace(c.Param("order_id"))

	var donasi *models.DonasiOnline
	if orderID != "" {
		donasi, _ = h.donasiService.GetByOrderID(orderID)
	}

	data := map[string]interface{}{
		"Title":      "Pembayaran Berhasil",
		"ActivePage": "program",
		"Year":       time.Now().Year(),
		"Donasi":     donasi,
	}

	return c.Render(http.StatusOK, "public/donasi_sukses.html", data)
}

// ---- POST /api/donasi/cancel ----
// Body JSON: { "order_id": "..." }
func (h *DonasiHandler) CancelTransaction(c echo.Context) error {
	type reqBody struct {
		OrderID string `json:"order_id"`
	}

	var body reqBody
	if err := c.Bind(&body); err != nil || body.OrderID == "" {
		return JSONBadRequest(c, "order_id tidak valid")
	}

	if err := h.donasiService.CancelTransaction(body.OrderID); err != nil {
		return JSONBadRequest(c, err.Error())
	}

	return JSONOk(c, "Transaksi berhasil dibatalkan")
}
