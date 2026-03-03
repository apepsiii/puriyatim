package repository

import (
	"database/sql"
	"fmt"
	"time"

	"puriyatim-app/internal/models"
)

type DonasiOnlineRepository struct {
	db *sql.DB
}

func NewDonasiOnlineRepository(db *sql.DB) *DonasiOnlineRepository {
	return &DonasiOnlineRepository{db: db}
}

// Create menyimpan record donasi online baru
func (r *DonasiOnlineRepository) Create(d *models.DonasiOnline) error {
	query := `
		INSERT INTO DONASI_ONLINE
			(order_id, jenis, nama_donatur, nominal, payment_method, status,
			 qr_string, va_number, total_payment, fee, expired_at, catatan)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		d.OrderID, string(d.Jenis), d.NamaDonatur, d.Nominal, d.PaymentMethod, string(d.Status),
		d.QRString, d.VANumber, d.TotalPayment, d.Fee, d.ExpiredAt, d.Catatan,
	)
	if err != nil {
		return fmt.Errorf("donasi_online create: %w", err)
	}

	id, _ := result.LastInsertId()
	d.ID = id
	return nil
}

// GetByOrderID mengambil satu record berdasarkan order_id
func (r *DonasiOnlineRepository) GetByOrderID(orderID string) (*models.DonasiOnline, error) {
	query := `
		SELECT id, order_id, jenis, nama_donatur, nominal, payment_method, status,
		       COALESCE(qr_string,''), COALESCE(va_number,''), COALESCE(total_payment,0),
		       COALESCE(fee,0), expired_at, completed_at, COALESCE(catatan,''),
		       created_at, updated_at
		FROM DONASI_ONLINE
		WHERE order_id = ?`

	row := r.db.QueryRow(query, orderID)
	return scanDonasiOnline(row)
}

// GetByID mengambil record berdasarkan primary key
func (r *DonasiOnlineRepository) GetByID(id int64) (*models.DonasiOnline, error) {
	query := `
		SELECT id, order_id, jenis, nama_donatur, nominal, payment_method, status,
		       COALESCE(qr_string,''), COALESCE(va_number,''), COALESCE(total_payment,0),
		       COALESCE(fee,0), expired_at, completed_at, COALESCE(catatan,''),
		       created_at, updated_at
		FROM DONASI_ONLINE
		WHERE id = ?`

	row := r.db.QueryRow(query, id)
	return scanDonasiOnline(row)
}

// UpdateStatus mengubah status transaksi, dan jika completed juga set completed_at
func (r *DonasiOnlineRepository) UpdateStatus(orderID string, status models.StatusDonasiOnline, completedAt *time.Time) error {
	query := `UPDATE DONASI_ONLINE SET status = ?, completed_at = ? WHERE order_id = ?`
	_, err := r.db.Exec(query, string(status), completedAt, orderID)
	if err != nil {
		return fmt.Errorf("donasi_online update_status: %w", err)
	}
	return nil
}

// GetAll mengambil semua record, diurutkan terbaru dulu
func (r *DonasiOnlineRepository) GetAll(limit int) ([]*models.DonasiOnline, error) {
	query := `
		SELECT id, order_id, jenis, nama_donatur, nominal, payment_method, status,
		       COALESCE(qr_string,''), COALESCE(va_number,''), COALESCE(total_payment,0),
		       COALESCE(fee,0), expired_at, completed_at, COALESCE(catatan,''),
		       created_at, updated_at
		FROM DONASI_ONLINE
		ORDER BY created_at DESC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("donasi_online get_all: %w", err)
	}
	defer rows.Close()

	var result []*models.DonasiOnline
	for rows.Next() {
		d, err := scanDonasiOnlineRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

// GetByStatus mengambil record berdasarkan status
func (r *DonasiOnlineRepository) GetByStatus(status models.StatusDonasiOnline) ([]*models.DonasiOnline, error) {
	query := `
		SELECT id, order_id, jenis, nama_donatur, nominal, payment_method, status,
		       COALESCE(qr_string,''), COALESCE(va_number,''), COALESCE(total_payment,0),
		       COALESCE(fee,0), expired_at, completed_at, COALESCE(catatan,''),
		       created_at, updated_at
		FROM DONASI_ONLINE
		WHERE status = ?
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, string(status))
	if err != nil {
		return nil, fmt.Errorf("donasi_online get_by_status: %w", err)
	}
	defer rows.Close()

	var result []*models.DonasiOnline
	for rows.Next() {
		d, err := scanDonasiOnlineRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

// MarkExpiredOld menandai semua transaksi pending yang sudah lewat expired_at sebagai expired
func (r *DonasiOnlineRepository) MarkExpiredOld() (int64, error) {
	result, err := r.db.Exec(`
		UPDATE DONASI_ONLINE
		SET status = 'expired'
		WHERE status = 'pending'
		  AND expired_at IS NOT NULL
		  AND expired_at < ?`, time.Now())
	if err != nil {
		return 0, fmt.Errorf("donasi_online mark_expired: %w", err)
	}
	n, _ := result.RowsAffected()
	return n, nil
}

// ---- helpers ----

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanDonasiOnline(row *sql.Row) (*models.DonasiOnline, error) {
	d := &models.DonasiOnline{}
	err := row.Scan(
		&d.ID, &d.OrderID, &d.Jenis, &d.NamaDonatur, &d.Nominal,
		&d.PaymentMethod, &d.Status,
		&d.QRString, &d.VANumber, &d.TotalPayment, &d.Fee,
		&d.ExpiredAt, &d.CompletedAt, &d.Catatan,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("donasi tidak ditemukan")
	}
	if err != nil {
		return nil, fmt.Errorf("donasi_online scan: %w", err)
	}
	return d, nil
}

func scanDonasiOnlineRow(rows *sql.Rows) (*models.DonasiOnline, error) {
	d := &models.DonasiOnline{}
	err := rows.Scan(
		&d.ID, &d.OrderID, &d.Jenis, &d.NamaDonatur, &d.Nominal,
		&d.PaymentMethod, &d.Status,
		&d.QRString, &d.VANumber, &d.TotalPayment, &d.Fee,
		&d.ExpiredAt, &d.CompletedAt, &d.Catatan,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("donasi_online scan rows: %w", err)
	}
	return d, nil
}
