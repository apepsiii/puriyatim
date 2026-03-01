package repository

import (
	"database/sql"
	"fmt"
	"time"

	"puriyatim-app/internal/models"

	"github.com/google/uuid"
)

type GaleriRepository struct {
	db *sql.DB
}

func NewGaleriRepository(db *sql.DB) *GaleriRepository {
	repo := &GaleriRepository{db: db}
	repo.ensureTable()
	return repo
}

func (r *GaleriRepository) ensureTable() {
	query := `
		CREATE TABLE IF NOT EXISTS GALERI_FOTO (
			id_foto TEXT PRIMARY KEY,
			judul TEXT NOT NULL,
			deskripsi TEXT,
			gambar_asli_url TEXT NOT NULL,
			gambar_overlay_url TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, _ = r.db.Exec(query)
}

func (r *GaleriRepository) generateID() string {
	return uuid.New().String()[:8]
}

func (r *GaleriRepository) Create(item *models.GaleriFoto) error {
	if item.ID == "" {
		item.ID = r.generateID()
	}
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	query := `
		INSERT INTO GALERI_FOTO (id_foto, judul, deskripsi, gambar_asli_url, gambar_overlay_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, item.ID, item.Judul, item.Deskripsi, item.GambarAsliURL, item.GambarOverlayURL, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create galeri foto: %w", err)
	}
	return nil
}

func (r *GaleriRepository) ListAll() ([]*models.GaleriFoto, error) {
	query := `
		SELECT id_foto, judul, deskripsi, gambar_asli_url, gambar_overlay_url, created_at, updated_at
		FROM GALERI_FOTO
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query galeri: %w", err)
	}
	defer rows.Close()

	var list []*models.GaleriFoto
	for rows.Next() {
		var item models.GaleriFoto
		var deskripsi sql.NullString
		if err := rows.Scan(&item.ID, &item.Judul, &deskripsi, &item.GambarAsliURL, &item.GambarOverlayURL, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan galeri: %w", err)
		}
		if deskripsi.Valid {
			item.Deskripsi = deskripsi.String
		}
		list = append(list, &item)
	}
	return list, nil
}

func (r *GaleriRepository) GetByID(id string) (*models.GaleriFoto, error) {
	query := `
		SELECT id_foto, judul, deskripsi, gambar_asli_url, gambar_overlay_url, created_at, updated_at
		FROM GALERI_FOTO
		WHERE id_foto = ?
	`
	var item models.GaleriFoto
	var deskripsi sql.NullString
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.Judul, &deskripsi, &item.GambarAsliURL, &item.GambarOverlayURL, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("foto galeri tidak ditemukan")
		}
		return nil, fmt.Errorf("failed to get galeri: %w", err)
	}
	if deskripsi.Valid {
		item.Deskripsi = deskripsi.String
	}
	return &item, nil
}

func (r *GaleriRepository) Update(item *models.GaleriFoto) error {
	item.UpdatedAt = time.Now()
	query := `
		UPDATE GALERI_FOTO
		SET judul = ?, deskripsi = ?, gambar_asli_url = ?, gambar_overlay_url = ?, updated_at = ?
		WHERE id_foto = ?
	`
	result, err := r.db.Exec(query, item.Judul, item.Deskripsi, item.GambarAsliURL, item.GambarOverlayURL, item.UpdatedAt, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update galeri: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update galeri: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("foto galeri tidak ditemukan")
	}
	return nil
}

func (r *GaleriRepository) Delete(id string) error {
	query := `DELETE FROM GALERI_FOTO WHERE id_foto = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete galeri: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete galeri: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("foto galeri tidak ditemukan")
	}
	return nil
}
