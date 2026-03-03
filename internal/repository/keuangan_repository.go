package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"puriyatim-app/internal/models"

	"github.com/google/uuid"
)

type KeuanganRepository struct {
	db *sql.DB
}

func NewKeuanganRepository(db *sql.DB) *KeuanganRepository {
	repo := &KeuanganRepository{db: db}
	repo.ensurePemasukanStatusColumn()
	return repo
}

func (r *KeuanganRepository) generateID() string {
	return uuid.New().String()[:8]
}

func (r *KeuanganRepository) ensurePemasukanStatusColumn() {
	query := `
		ALTER TABLE PEMASUKAN_DONASI
		ADD COLUMN status_verifikasi TEXT NOT NULL DEFAULT 'verified'
	`
	_, _ = r.db.Exec(query)
}

func (r *KeuanganRepository) CreatePemasukan(p *models.PemasukanDonasi) error {
	if p.ID == "" {
		p.ID = r.generateID()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.StatusVerifikasi == "" {
		p.StatusVerifikasi = models.StatusVerifikasiPending
	}

	query := `
		INSERT OR IGNORE INTO PEMASUKAN_DONASI (id_pemasukan, nama_donatur, tanggal_donasi, nominal, kategori_dana, catatan, bukti_transaksi, status_verifikasi, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, p.ID, p.NamaDonatur, p.TanggalDonasi, p.Nominal, p.KategoriDana, p.Catatan, p.BuktiTransaksi, p.StatusVerifikasi, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create pemasukan: %w", err)
	}
	return nil
}

func (r *KeuanganRepository) GetPemasukanByID(id string) (*models.PemasukanDonasi, error) {
	query := `
		SELECT id_pemasukan, nama_donatur, tanggal_donasi, nominal, kategori_dana, catatan, bukti_transaksi, status_verifikasi, created_at, updated_at
		FROM PEMASUKAN_DONASI WHERE id_pemasukan = ?
	`

	var p models.PemasukanDonasi
	var namaDonatur, catatan, bukti, status sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &namaDonatur, &p.TanggalDonasi, &p.Nominal, &p.KategoriDana, &catatan, &bukti, &status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pemasukan not found")
		}
		return nil, fmt.Errorf("failed to get pemasukan: %w", err)
	}

	if namaDonatur.Valid {
		p.NamaDonatur = namaDonatur.String
	} else {
		p.NamaDonatur = "Hamba Allah"
	}
	if catatan.Valid {
		p.Catatan = catatan.String
	}
	if bukti.Valid {
		p.BuktiTransaksi = bukti.String
	}
	p.StatusVerifikasi = normalizeStatus(status.String)

	return &p, nil
}

func (r *KeuanganRepository) GetAllPemasukan() ([]*models.PemasukanDonasi, error) {
	query := `
		SELECT id_pemasukan, nama_donatur, tanggal_donasi, nominal, kategori_dana, catatan, bukti_transaksi, status_verifikasi, created_at, updated_at
		FROM PEMASUKAN_DONASI ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pemasukan: %w", err)
	}
	defer rows.Close()

	var list []*models.PemasukanDonasi
	for rows.Next() {
		var p models.PemasukanDonasi
		var namaDonatur, catatan, bukti, status sql.NullString

		if err := rows.Scan(&p.ID, &namaDonatur, &p.TanggalDonasi, &p.Nominal, &p.KategoriDana, &catatan, &bukti, &status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pemasukan: %w", err)
		}

		if namaDonatur.Valid {
			p.NamaDonatur = namaDonatur.String
		} else {
			p.NamaDonatur = "Hamba Allah"
		}
		if catatan.Valid {
			p.Catatan = catatan.String
		}
		if bukti.Valid {
			p.BuktiTransaksi = bukti.String
		}
		p.StatusVerifikasi = normalizeStatus(status.String)

		list = append(list, &p)
	}

	return list, nil
}

func (r *KeuanganRepository) GetPemasukanByDateRange(start, end time.Time) ([]*models.PemasukanDonasi, error) {
	query := `
		SELECT id_pemasukan, nama_donatur, tanggal_donasi, nominal, kategori_dana, catatan, bukti_transaksi, status_verifikasi, created_at, updated_at
		FROM PEMASUKAN_DONASI WHERE tanggal_donasi BETWEEN ? AND ? ORDER BY tanggal_donasi DESC
	`

	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query pemasukan by date: %w", err)
	}
	defer rows.Close()

	var list []*models.PemasukanDonasi
	for rows.Next() {
		var p models.PemasukanDonasi
		var namaDonatur, catatan, bukti, status sql.NullString

		if err := rows.Scan(&p.ID, &namaDonatur, &p.TanggalDonasi, &p.Nominal, &p.KategoriDana, &catatan, &bukti, &status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pemasukan: %w", err)
		}

		if namaDonatur.Valid {
			p.NamaDonatur = namaDonatur.String
		} else {
			p.NamaDonatur = "Hamba Allah"
		}
		if catatan.Valid {
			p.Catatan = catatan.String
		}
		if bukti.Valid {
			p.BuktiTransaksi = bukti.String
		}
		p.StatusVerifikasi = normalizeStatus(status.String)

		list = append(list, &p)
	}

	return list, nil
}

func (r *KeuanganRepository) CreatePengeluaran(p *models.Pengeluaran) error {
	if p.ID == "" {
		p.ID = r.generateID()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	query := `
		INSERT INTO PENGELUARAN (id_pengeluaran, tanggal_pengeluaran, nominal, id_anak, keterangan, bukti_transaksi, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	var idAnak interface{}
	if p.IDAnak != nil {
		idAnak = *p.IDAnak
	}

	_, err := r.db.Exec(query, p.ID, p.TanggalPengeluaran, p.Nominal, idAnak, p.Keterangan, p.BuktiTransaksi, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create pengeluaran: %w", err)
	}
	return nil
}

func (r *KeuanganRepository) GetPengeluaranByID(id string) (*models.Pengeluaran, error) {
	query := `
		SELECT id_pengeluaran, tanggal_pengeluaran, nominal, id_anak, keterangan, bukti_transaksi, created_at, updated_at
		FROM PENGELUARAN WHERE id_pengeluaran = ?
	`

	var p models.Pengeluaran
	var idAnak, bukti sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.TanggalPengeluaran, &p.Nominal, &idAnak, &p.Keterangan, &bukti, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pengeluaran not found")
		}
		return nil, fmt.Errorf("failed to get pengeluaran: %w", err)
	}

	if idAnak.Valid {
		p.IDAnak = &idAnak.String
	}
	if bukti.Valid {
		p.BuktiTransaksi = bukti.String
	}

	return &p, nil
}

func (r *KeuanganRepository) GetPengeluaranByAnakID(anakID string) ([]*models.Pengeluaran, error) {
	query := `
		SELECT id_pengeluaran, tanggal_pengeluaran, nominal, id_anak, keterangan, bukti_transaksi, created_at, updated_at
		FROM PENGELUARAN WHERE id_anak = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, anakID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pengeluaran by anak: %w", err)
	}
	defer rows.Close()

	var list []*models.Pengeluaran
	for rows.Next() {
		var p models.Pengeluaran
		var idAnak, bukti sql.NullString

		if err := rows.Scan(&p.ID, &p.TanggalPengeluaran, &p.Nominal, &idAnak, &p.Keterangan, &bukti, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pengeluaran: %w", err)
		}

		if idAnak.Valid {
			p.IDAnak = &idAnak.String
		}
		if bukti.Valid {
			p.BuktiTransaksi = bukti.String
		}

		list = append(list, &p)
	}

	return list, nil
}

func (r *KeuanganRepository) GetAllPengeluaran() ([]*models.Pengeluaran, error) {
	query := `
		SELECT p.id_pengeluaran, p.tanggal_pengeluaran, p.nominal, p.id_anak, p.keterangan, p.bukti_transaksi, p.created_at, p.updated_at,
			   a.nama_lengkap
		FROM PENGELUARAN p
		LEFT JOIN ANAK_ASUH a ON p.id_anak = a.id_anak
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pengeluaran: %w", err)
	}
	defer rows.Close()

	var list []*models.Pengeluaran
	for rows.Next() {
		var p models.Pengeluaran
		var idAnak, bukti, anakNama sql.NullString

		if err := rows.Scan(&p.ID, &p.TanggalPengeluaran, &p.Nominal, &idAnak, &p.Keterangan, &bukti, &p.CreatedAt, &p.UpdatedAt, &anakNama); err != nil {
			return nil, fmt.Errorf("failed to scan pengeluaran: %w", err)
		}

		if idAnak.Valid {
			p.IDAnak = &idAnak.String
			if anakNama.Valid {
				p.Anak = &models.AnakAsuh{NamaLengkap: anakNama.String}
			}
		}
		if bukti.Valid {
			p.BuktiTransaksi = bukti.String
		}

		list = append(list, &p)
	}

	return list, nil
}

func (r *KeuanganRepository) GetPengeluaranByDateRange(start, end time.Time) ([]*models.Pengeluaran, error) {
	query := `
		SELECT p.id_pengeluaran, p.tanggal_pengeluaran, p.nominal, p.id_anak, p.keterangan, p.bukti_transaksi, p.created_at, p.updated_at,
			   a.nama_lengkap
		FROM PENGELUARAN p
		LEFT JOIN ANAK_ASUH a ON p.id_anak = a.id_anak
		WHERE p.tanggal_pengeluaran BETWEEN ? AND ? ORDER BY p.tanggal_pengeluaran DESC
	`

	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query pengeluaran by date: %w", err)
	}
	defer rows.Close()

	var list []*models.Pengeluaran
	for rows.Next() {
		var p models.Pengeluaran
		var idAnak, bukti, anakNama sql.NullString

		if err := rows.Scan(&p.ID, &p.TanggalPengeluaran, &p.Nominal, &idAnak, &p.Keterangan, &bukti, &p.CreatedAt, &p.UpdatedAt, &anakNama); err != nil {
			return nil, fmt.Errorf("failed to scan pengeluaran: %w", err)
		}

		if idAnak.Valid {
			p.IDAnak = &idAnak.String
			if anakNama.Valid {
				p.Anak = &models.AnakAsuh{NamaLengkap: anakNama.String}
			}
		}
		if bukti.Valid {
			p.BuktiTransaksi = bukti.String
		}

		list = append(list, &p)
	}

	return list, nil
}

func (r *KeuanganRepository) GetTotalPemasukan() (float64, error) {
	query := `SELECT COALESCE(SUM(nominal), 0) FROM PEMASUKAN_DONASI WHERE status_verifikasi = 'verified'`
	var total float64
	err := r.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total pemasukan: %w", err)
	}
	return total, nil
}

func (r *KeuanganRepository) GetTotalPengeluaran() (float64, error) {
	query := `SELECT COALESCE(SUM(nominal), 0) FROM PENGELUARAN`
	var total float64
	err := r.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total pengeluaran: %w", err)
	}
	return total, nil
}

func (r *KeuanganRepository) GetTotalPemasukanByMonth(year, month int) (float64, error) {
	query := `SELECT COALESCE(SUM(nominal), 0) FROM PEMASUKAN_DONASI WHERE status_verifikasi = 'verified' AND strftime('%Y', tanggal_donasi) = ? AND strftime('%m', tanggal_donasi) = ?`
	var total float64
	monthStr := fmt.Sprintf("%02d", month)
	err := r.db.QueryRow(query, fmt.Sprintf("%d", year), monthStr).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total pemasukan by month: %w", err)
	}
	return total, nil
}

func (r *KeuanganRepository) GetTotalPengeluaranByMonth(year, month int) (float64, error) {
	query := `SELECT COALESCE(SUM(nominal), 0) FROM PENGELUARAN WHERE strftime('%Y', tanggal_pengeluaran) = ? AND strftime('%m', tanggal_pengeluaran) = ?`
	var total float64
	monthStr := fmt.Sprintf("%02d", month)
	err := r.db.QueryRow(query, fmt.Sprintf("%d", year), monthStr).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total pengeluaran by month: %w", err)
	}
	return total, nil
}

func (r *KeuanganRepository) CountPemasukan() (int, error) {
	query := `SELECT COUNT(*) FROM PEMASUKAN_DONASI`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pemasukan: %w", err)
	}
	return count, nil
}

func (r *KeuanganRepository) CountPengeluaran() (int, error) {
	query := `SELECT COUNT(*) FROM PENGELUARAN`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pengeluaran: %w", err)
	}
	return count, nil
}

func (r *KeuanganRepository) GetAllDonatur() ([]models.Donatur, error) {
	query := `SELECT id_donatur, nama_donatur, tipe_donatur, no_telepon, email, alamat, catatan_khusus, tanggal_bergabung, created_at, updated_at FROM DONATUR ORDER BY nama_donatur ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query donatur: %w", err)
	}
	defer rows.Close()

	var list []models.Donatur
	for rows.Next() {
		var d models.Donatur
		var email, alamat sql.NullString
		var catatan sql.NullString

		err := rows.Scan(
			&d.ID, &d.NamaDonatur, &d.TipeDonatur, &d.NoTelepon, &email, &alamat, &catatan, &d.TanggalBergabung, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan donatur: %w", err)
		}

		if email.Valid {
			d.Email = &email.String
		}
		if alamat.Valid {
			d.Alamat = &alamat.String
		}
		if catatan.Valid {
			d.CatatanKhusus = catatan.String
		}

		list = append(list, d)
	}

	return list, nil
}

func (r *KeuanganRepository) CreateDonatur(d *models.Donatur) error {
	if d.ID == "" {
		d.ID = r.generateID()
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.TanggalBergabung.IsZero() {
		d.TanggalBergabung = now
	}

	query := `
		INSERT INTO DONATUR (id_donatur, nama_donatur, tipe_donatur, no_telepon, email, alamat, catatan_khusus, tanggal_bergabung, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var email, alamat interface{}
	if d.Email != nil {
		email = *d.Email
	}
	if d.Alamat != nil {
		alamat = *d.Alamat
	}

	_, err := r.db.Exec(query, d.ID, d.NamaDonatur, d.TipeDonatur, d.NoTelepon, email, alamat, d.CatatanKhusus, d.TanggalBergabung, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create donatur: %w", err)
	}
	return nil
}

func (r *KeuanganRepository) DeletePemasukan(id string) error {
	query := `DELETE FROM PEMASUKAN_DONASI WHERE id_pemasukan = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pemasukan: %w", err)
	}
	return nil
}

func (r *KeuanganRepository) DeletePengeluaran(id string) error {
	query := `DELETE FROM PENGELUARAN WHERE id_pengeluaran = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pengeluaran: %w", err)
	}
	return nil
}

func (r *KeuanganRepository) UpdatePemasukan(p *models.PemasukanDonasi) error {
	p.UpdatedAt = time.Now()

	query := `
		UPDATE PEMASUKAN_DONASI SET
			nama_donatur = ?, tanggal_donasi = ?, nominal = ?, kategori_dana = ?, catatan = ?, bukti_transaksi = ?, status_verifikasi = ?, updated_at = ?
		WHERE id_pemasukan = ?
	`

	_, err := r.db.Exec(query, p.NamaDonatur, p.TanggalDonasi, p.Nominal, p.KategoriDana, p.Catatan, p.BuktiTransaksi, p.StatusVerifikasi, p.UpdatedAt, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update pemasukan: %w", err)
	}
	return nil
}

func (r *KeuanganRepository) UpdatePengeluaran(p *models.Pengeluaran) error {
	p.UpdatedAt = time.Now()

	query := `
		UPDATE PENGELUARAN SET
			tanggal_pengeluaran = ?, nominal = ?, keterangan = ?, id_anak = ?, bukti_transaksi = ?, updated_at = ?
		WHERE id_pengeluaran = ?
	`

	var idAnak interface{}
	if p.IDAnak != nil {
		idAnak = *p.IDAnak
	}

	_, err := r.db.Exec(query, p.TanggalPengeluaran, p.Nominal, p.Keterangan, idAnak, p.BuktiTransaksi, p.UpdatedAt, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update pengeluaran: %w", err)
	}
	return nil
}

func (r *KeuanganRepository) GetPemasukanByMonthStr(monthStr string) ([]*models.PemasukanDonasi, error) {
	if monthStr == "" {
		return r.GetAllPemasukan()
	}

	parts := strings.Split(monthStr, "-")
	if len(parts) != 2 {
		return r.GetAllPemasukan()
	}
	month := parts[0]
	year := parts[1]

	query := `
		SELECT id_pemasukan, nama_donatur, tanggal_donasi, nominal, kategori_dana, catatan, bukti_transaksi, status_verifikasi, created_at, updated_at
		FROM PEMASUKAN_DONASI 
		WHERE strftime('%m', tanggal_donasi) = ? AND strftime('%Y', tanggal_donasi) = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to query pemasukan by month: %w", err)
	}
	defer rows.Close()

	var list []*models.PemasukanDonasi
	for rows.Next() {
		var p models.PemasukanDonasi
		var namaDonatur, catatan, bukti, status sql.NullString

		if err := rows.Scan(&p.ID, &namaDonatur, &p.TanggalDonasi, &p.Nominal, &p.KategoriDana, &catatan, &bukti, &status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pemasukan: %w", err)
		}

		if namaDonatur.Valid {
			p.NamaDonatur = namaDonatur.String
		} else {
			p.NamaDonatur = "Hamba Allah"
		}
		if catatan.Valid {
			p.Catatan = catatan.String
		}
		if bukti.Valid {
			p.BuktiTransaksi = bukti.String
		}
		p.StatusVerifikasi = normalizeStatus(status.String)

		list = append(list, &p)
	}

	return list, nil
}

func (r *KeuanganRepository) VerifyPemasukan(id string) error {
	query := `
		UPDATE PEMASUKAN_DONASI
		SET status_verifikasi = 'verified', updated_at = ?
		WHERE id_pemasukan = ?
	`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to verify pemasukan: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to verify pemasukan: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("pemasukan not found")
	}
	return nil
}

func normalizeStatus(status string) models.StatusVerifikasiPemasukan {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case string(models.StatusVerifikasiPending):
		return models.StatusVerifikasiPending
	case string(models.StatusVerifikasiVerified):
		return models.StatusVerifikasiVerified
	default:
		return models.StatusVerifikasiVerified
	}
}

func (r *KeuanganRepository) GetPengeluaranByMonthStr(monthStr string) ([]*models.Pengeluaran, error) {
	if monthStr == "" {
		return r.GetAllPengeluaran()
	}

	parts := strings.Split(monthStr, "-")
	if len(parts) != 2 {
		return r.GetAllPengeluaran()
	}
	month := parts[0]
	year := parts[1]

	query := `
		SELECT p.id_pengeluaran, p.tanggal_pengeluaran, p.nominal, p.id_anak, p.keterangan, p.bukti_transaksi, p.created_at, p.updated_at,
			   a.nama_lengkap
		FROM PENGELUARAN p
		LEFT JOIN ANAK_ASUH a ON p.id_anak = a.id_anak
		WHERE strftime('%m', p.tanggal_pengeluaran) = ? AND strftime('%Y', p.tanggal_pengeluaran) = ?
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to query pengeluaran by month: %w", err)
	}
	defer rows.Close()

	var list []*models.Pengeluaran
	for rows.Next() {
		var p models.Pengeluaran
		var idAnak, bukti, anakNama sql.NullString

		if err := rows.Scan(&p.ID, &p.TanggalPengeluaran, &p.Nominal, &idAnak, &p.Keterangan, &bukti, &p.CreatedAt, &p.UpdatedAt, &anakNama); err != nil {
			return nil, fmt.Errorf("failed to scan pengeluaran: %w", err)
		}

		if idAnak.Valid {
			p.IDAnak = &idAnak.String
			if anakNama.Valid {
				p.Anak = &models.AnakAsuh{NamaLengkap: anakNama.String}
			}
		}
		if bukti.Valid {
			p.BuktiTransaksi = bukti.String
		}

		list = append(list, &p)
	}

	return list, nil
}

func (r *KeuanganRepository) GetPemasukanByNomorHP(nomorHP string) ([]*models.PemasukanDonasi, error) {
	query := `
		SELECT id_pemasukan, nama_donatur, tanggal_donasi, nominal, kategori_dana, catatan, bukti_transaksi, status_verifikasi, created_at, updated_at
		FROM PEMASUKAN_DONASI
		WHERE LOWER(catatan) LIKE ?
		ORDER BY created_at DESC
		LIMIT 100
	`
	pattern := "%nomor hp: " + strings.ToLower(strings.TrimSpace(nomorHP)) + "%"

	rows, err := r.db.Query(query, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to query pemasukan by nomor hp: %w", err)
	}
	defer rows.Close()

	var list []*models.PemasukanDonasi
	for rows.Next() {
		var p models.PemasukanDonasi
		var namaDonatur, catatan, bukti, status sql.NullString

		if err := rows.Scan(&p.ID, &namaDonatur, &p.TanggalDonasi, &p.Nominal, &p.KategoriDana, &catatan, &bukti, &status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pemasukan by nomor hp: %w", err)
		}
		if namaDonatur.Valid {
			p.NamaDonatur = namaDonatur.String
		} else {
			p.NamaDonatur = "Hamba Allah"
		}
		if catatan.Valid {
			p.Catatan = catatan.String
		}
		if bukti.Valid {
			p.BuktiTransaksi = bukti.String
		}
		p.StatusVerifikasi = normalizeStatus(status.String)

		list = append(list, &p)
	}

	return list, nil
}
