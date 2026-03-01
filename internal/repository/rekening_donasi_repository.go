package repository

import (
	"database/sql"
	"puriyatim-app/internal/models"
	"time"
)

type RekeningDonasiRepository struct {
	db *sql.DB
}

func NewRekeningDonasiRepository(db *sql.DB) *RekeningDonasiRepository {
	return &RekeningDonasiRepository{db: db}
}

func (r *RekeningDonasiRepository) GetAll() ([]*models.RekeningDonasi, error) {
	rows, err := r.db.Query(`
		SELECT id, nama_bank, logo_bank, nomor_rekening, atas_nama, urutan, aktif
		FROM REKENING_DONASI
		ORDER BY urutan ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.RekeningDonasi
	for rows.Next() {
		var item models.RekeningDonasi
		var aktifInt int
		if err := rows.Scan(&item.ID, &item.NamaBank, &item.LogoBank, &item.NomorRekening, &item.AtasNama, &item.Urutan, &aktifInt); err != nil {
			return nil, err
		}
		item.Aktif = aktifInt == 1
		list = append(list, &item)
	}
	return list, rows.Err()
}

func (r *RekeningDonasiRepository) Create(item *models.RekeningDonasi) error {
	_, err := r.db.Exec(`
		INSERT INTO REKENING_DONASI (nama_bank, logo_bank, nomor_rekening, atas_nama, urutan, aktif, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, item.NamaBank, item.LogoBank, item.NomorRekening, item.AtasNama, item.Urutan, boolToInt(item.Aktif),
		time.Now(), time.Now())
	return err
}

func (r *RekeningDonasiRepository) Update(item *models.RekeningDonasi) error {
	_, err := r.db.Exec(`
		UPDATE REKENING_DONASI
		SET nama_bank=?, logo_bank=?, nomor_rekening=?, atas_nama=?, urutan=?, aktif=?, updated_at=?
		WHERE id=?
	`, item.NamaBank, item.LogoBank, item.NomorRekening, item.AtasNama, item.Urutan, boolToInt(item.Aktif),
		time.Now(), item.ID)
	return err
}

func (r *RekeningDonasiRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM REKENING_DONASI WHERE id = ?`, id)
	return err
}

func (r *RekeningDonasiRepository) ReorderAll(ids []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE REKENING_DONASI SET urutan=?, updated_at=? WHERE id=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for i, id := range ids {
		if _, err := stmt.Exec(i+1, now, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
