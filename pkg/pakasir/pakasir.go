// Package pakasir menyediakan client untuk integrasi API Pakasir payment gateway.
// Dokumentasi: https://pakasir.com/p/docs
package pakasir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://app.pakasir.com"

// PaymentMethod yang tersedia di Pakasir
const (
	MethodQRIS         = "qris"
	MethodBRIVA        = "bri_va"
	MethodBNIVA        = "bni_va"
	MethodCIMBVA       = "cimb_niaga_va"
	MethodSampoernaVA  = "sampoerna_va"
	MethodBNCVA        = "bnc_va"
	MethodMaybankVA    = "maybank_va"
	MethodPermataVA    = "permata_va"
	MethodATMBersamaVA = "atm_bersama_va"
	MethodArthaGrahaVA = "artha_graha_va"
	MethodPaypal       = "paypal"
)

// PaymentMethodLabels label tampilan per metode
var PaymentMethodLabels = map[string]string{
	MethodQRIS:         "QRIS",
	MethodBRIVA:        "Virtual Account BRI",
	MethodBNIVA:        "Virtual Account BNI",
	MethodCIMBVA:       "Virtual Account CIMB Niaga",
	MethodSampoernaVA:  "Virtual Account Bank Sampoerna",
	MethodBNCVA:        "Virtual Account BNC",
	MethodMaybankVA:    "Virtual Account Maybank",
	MethodPermataVA:    "Virtual Account Permata",
	MethodATMBersamaVA: "ATM Bersama",
	MethodArthaGrahaVA: "Virtual Account Artha Graha",
	MethodPaypal:       "PayPal",
}

// Client adalah Pakasir API client
type Client struct {
	projectSlug string
	apiKey      string
	httpClient  *http.Client
}

// NewClient membuat instance client baru
func NewClient(projectSlug, apiKey string) *Client {
	return &Client{
		projectSlug: projectSlug,
		apiKey:      apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured memeriksa apakah config sudah terisi
func (c *Client) IsConfigured() bool {
	return c.projectSlug != "" && c.apiKey != ""
}

// ---- Request / Response structs ----

type transactionRequest struct {
	Project string `json:"project"`
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
	APIKey  string `json:"api_key"`
}

// Payment adalah detail pembayaran dari response Pakasir
type Payment struct {
	Project       string    `json:"project"`
	OrderID       string    `json:"order_id"`
	Amount        int64     `json:"amount"`
	Fee           int64     `json:"fee"`
	TotalPayment  int64     `json:"total_payment"`
	PaymentMethod string    `json:"payment_method"`
	PaymentNumber string    `json:"payment_number"` // QR string atau VA number
	ExpiredAt     time.Time `json:"expired_at"`
}

// CreateTransactionResponse adalah response dari API create transaction
type CreateTransactionResponse struct {
	Payment *Payment `json:"payment"`
}

// TransactionDetailResponse adalah response dari API transaction detail
type TransactionDetail struct {
	Amount        int64      `json:"amount"`
	OrderID       string     `json:"order_id"`
	Project       string     `json:"project"`
	Status        string     `json:"status"`
	PaymentMethod string     `json:"payment_method"`
	CompletedAt   *time.Time `json:"completed_at"`
}

type TransactionDetailResponse struct {
	Transaction *TransactionDetail `json:"transaction"`
}

// ---- API Methods ----

// CreateTransaction membuat transaksi baru di Pakasir.
// method: salah satu konstanta MethodXxx di atas.
// orderID: ID unik dari sistem kita.
// amount: nominal dalam rupiah (integer).
func (c *Client) CreateTransaction(method, orderID string, amount int64) (*Payment, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("pakasir: project slug atau api key belum dikonfigurasi")
	}

	reqBody := transactionRequest{
		Project: c.projectSlug,
		OrderID: orderID,
		Amount:  amount,
		APIKey:  c.apiKey,
	}

	endpoint := fmt.Sprintf("%s/api/transactioncreate/%s", baseURL, method)
	resp, err := c.doPost(endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("pakasir: gagal membuat transaksi: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("pakasir: server error %d: %s", resp.StatusCode, string(body))
	}

	var result CreateTransactionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("pakasir: gagal decode response: %w", err)
	}

	if result.Payment == nil {
		return nil, fmt.Errorf("pakasir: response tidak valid (payment null). body: %s", string(body))
	}

	return result.Payment, nil
}

// CancelTransaction membatalkan transaksi yang masih pending.
func (c *Client) CancelTransaction(orderID string, amount int64) error {
	if !c.IsConfigured() {
		return fmt.Errorf("pakasir: project slug atau api key belum dikonfigurasi")
	}

	reqBody := transactionRequest{
		Project: c.projectSlug,
		OrderID: orderID,
		Amount:  amount,
		APIKey:  c.apiKey,
	}

	resp, err := c.doPost(baseURL+"/api/transactioncancel", reqBody)
	if err != nil {
		return fmt.Errorf("pakasir: gagal cancel transaksi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pakasir: server error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetTransactionDetail mengambil detail dan status transaksi.
func (c *Client) GetTransactionDetail(orderID string, amount int64) (*TransactionDetail, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("pakasir: project slug atau api key belum dikonfigurasi")
	}

	url := fmt.Sprintf("%s/api/transactiondetail?project=%s&amount=%d&order_id=%s&api_key=%s",
		baseURL, c.projectSlug, amount, orderID, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("pakasir: gagal ambil detail transaksi: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pakasir: server error %d: %s", resp.StatusCode, string(body))
	}

	var result TransactionDetailResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("pakasir: gagal decode detail: %w", err)
	}
	if result.Transaction == nil {
		return nil, fmt.Errorf("pakasir: transaksi tidak ditemukan. body: %s", string(body))
	}

	return result.Transaction, nil
}

// SimulatePayment untuk keperluan testing di mode sandbox.
func (c *Client) SimulatePayment(orderID string, amount int64) error {
	if !c.IsConfigured() {
		return fmt.Errorf("pakasir: project slug atau api key belum dikonfigurasi")
	}

	reqBody := transactionRequest{
		Project: c.projectSlug,
		OrderID: orderID,
		Amount:  amount,
		APIKey:  c.apiKey,
	}

	resp, err := c.doPost(baseURL+"/api/paymentsimulation", reqBody)
	if err != nil {
		return fmt.Errorf("pakasir: gagal simulasi pembayaran: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pakasir: server error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ---- WebhookPayload adalah struktur POST yang dikirim Pakasir saat pembayaran selesai ----
type WebhookPayload struct {
	Amount        int64      `json:"amount"`
	OrderID       string     `json:"order_id"`
	Project       string     `json:"project"`
	Status        string     `json:"status"` // "completed"
	PaymentMethod string     `json:"payment_method"`
	CompletedAt   *time.Time `json:"completed_at"`
}

// ---- helper ----

func (c *Client) doPost(url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}
