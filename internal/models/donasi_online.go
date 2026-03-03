package models

import "time"

// Jenis donasi online
type JenisDonasiOnline string

const (
	JenisDonasiUmum    JenisDonasiOnline = "donasi"
	JenisDonasiZakat   JenisDonasiOnline = "zakat"
	JenisDonasiJumat   JenisDonasiOnline = "jumat_berkah"
)

// Status transaksi donasi online
type StatusDonasiOnline string

const (
	StatusDonasiPending   StatusDonasiOnline = "pending"
	StatusDonasiCompleted StatusDonasiOnline = "completed"
	StatusDonasiExpired   StatusDonasiOnline = "expired"
	StatusDonasiCancelled StatusDonasiOnline = "cancelled"
)

// DonasiOnline adalah record satu transaksi pembayaran digital via Pakasir
type DonasiOnline struct {
	ID            int64              `json:"id"             db:"id"`
	OrderID       string             `json:"order_id"       db:"order_id"`
	Jenis         JenisDonasiOnline  `json:"jenis"          db:"jenis"`
	NamaDonatur   string             `json:"nama_donatur"   db:"nama_donatur"`
	Nominal       float64            `json:"nominal"        db:"nominal"`
	PaymentMethod string             `json:"payment_method" db:"payment_method"`
	Status        StatusDonasiOnline `json:"status"         db:"status"`
	QRString      string             `json:"qr_string"      db:"qr_string"`
	VANumber      string             `json:"va_number"      db:"va_number"`
	TotalPayment  float64            `json:"total_payment"  db:"total_payment"`
	Fee           float64            `json:"fee"            db:"fee"`
	ExpiredAt     *time.Time         `json:"expired_at"     db:"expired_at"`
	CompletedAt   *time.Time         `json:"completed_at"   db:"completed_at"`
	Catatan       string             `json:"catatan"        db:"catatan"`
	CreatedAt     time.Time          `json:"created_at"     db:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"     db:"updated_at"`
}

// IsQRIS memeriksa apakah metode pembayaran adalah QRIS
func (d *DonasiOnline) IsQRIS() bool {
	return d.PaymentMethod == "qris"
}

// IsVA memeriksa apakah metode pembayaran adalah Virtual Account
func (d *DonasiOnline) IsVA() bool {
	return !d.IsQRIS() && d.PaymentMethod != "paypal"
}

// IsExpired memeriksa apakah transaksi sudah kadaluarsa
func (d *DonasiOnline) IsExpired() bool {
	if d.ExpiredAt == nil {
		return false
	}
	return time.Now().After(*d.ExpiredAt)
}
