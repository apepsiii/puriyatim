package repository

import (
	"database/sql"
	"fmt"
	"time"

	"puriyatim-app/internal/models"

	"github.com/google/uuid"
)

type ArtikelRepository struct {
	db *sql.DB
}

func NewArtikelRepository(db *sql.DB) *ArtikelRepository {
	return &ArtikelRepository{db: db}
}

func (r *ArtikelRepository) generateID() string {
	return uuid.New().String()[:8]
}

func (r *ArtikelRepository) Create(a *models.Artikel) error {
	if a.ID == "" {
		a.ID = r.generateID()
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	query := `
		INSERT INTO ARTIKEL (id_artikel, id_pengurus, id_kategori, judul, slug, konten_html_markdown, gambar_thumbnail_url, meta_deskripsi, status_publikasi, tanggal_terbit)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var tanggalTerbit interface{}
	if a.TanggalTerbit != nil {
		tanggalTerbit = *a.TanggalTerbit
	}

	var thumbnail interface{}
	if a.GambarThumbnail != nil {
		thumbnail = *a.GambarThumbnail
	}

	_, err := r.db.Exec(query, a.ID, a.IDPengurus, a.IDKategori, a.Judul, a.Slug, a.KontenHTML, thumbnail, a.MetaDeskripsi, a.StatusPublikasi, tanggalTerbit)
	if err != nil {
		return fmt.Errorf("failed to create artikel: %w", err)
	}
	return nil
}

func (r *ArtikelRepository) GetByID(id string) (*models.Artikel, error) {
	query := `
		SELECT id_artikel, id_pengurus, id_kategori, judul, slug, konten_html_markdown, gambar_thumbnail_url, meta_deskripsi, status_publikasi, tanggal_terbit, created_at, updated_at
		FROM ARTIKEL WHERE id_artikel = ?
	`

	var a models.Artikel
	var thumbnail, tanggalTerbit sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.IDPengurus, &a.IDKategori, &a.Judul, &a.Slug, &a.KontenHTML, &thumbnail, &a.MetaDeskripsi, &a.StatusPublikasi, &tanggalTerbit, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artikel not found")
		}
		return nil, fmt.Errorf("failed to get artikel: %w", err)
	}

	if thumbnail.Valid {
		a.GambarThumbnail = &thumbnail.String
	}
	if tanggalTerbit.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", tanggalTerbit.String)
		a.TanggalTerbit = &t
	}
	if createdAt.Valid {
		a.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		a.UpdatedAt = updatedAt.Time
	}

	return &a, nil
}

func (r *ArtikelRepository) GetBySlug(slug string) (*models.Artikel, error) {
	query := `
		SELECT a.id_artikel, a.id_pengurus, a.id_kategori, a.judul, a.slug, a.konten_html_markdown, a.gambar_thumbnail_url, a.meta_deskripsi, a.status_publikasi, a.tanggal_terbit, a.created_at, a.updated_at,
			   k.id_kategori, k.nama_kategori, k.slug
		FROM ARTIKEL a
		LEFT JOIN KATEGORI_KONTEN k ON a.id_kategori = k.id_kategori
		WHERE a.slug = ?
	`

	var a models.Artikel
	var thumbnail, tanggalTerbit sql.NullString
	var createdAt, updatedAt sql.NullTime
	var k models.KategoriKonten
	var kID sql.NullInt64
	var kNama, kSlug sql.NullString

	err := r.db.QueryRow(query, slug).Scan(
		&a.ID, &a.IDPengurus, &a.IDKategori, &a.Judul, &a.Slug, &a.KontenHTML, &thumbnail, &a.MetaDeskripsi, &a.StatusPublikasi, &tanggalTerbit, &createdAt, &updatedAt,
		&kID, &kNama, &kSlug,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artikel not found")
		}
		return nil, fmt.Errorf("failed to get artikel: %w", err)
	}

	if thumbnail.Valid && thumbnail.String != "" {
		a.GambarThumbnail = &thumbnail.String
	}
	if tanggalTerbit.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", tanggalTerbit.String)
		a.TanggalTerbit = &t
	}
	if createdAt.Valid {
		a.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		a.UpdatedAt = updatedAt.Time
	}
	if kID.Valid {
		k.ID = int(kID.Int64)
	}
	if kNama.Valid {
		k.NamaKategori = kNama.String
	}
	if kSlug.Valid {
		k.Slug = kSlug.String
	}
	a.Kategori = &k

	return &a, nil
}

func (r *ArtikelRepository) GetAll() ([]*models.Artikel, error) {
	query := `
		SELECT a.id_artikel, a.id_pengurus, a.id_kategori, a.judul, a.slug, a.konten_html_markdown, a.gambar_thumbnail_url, a.meta_deskripsi, a.status_publikasi, a.tanggal_terbit, a.created_at, a.updated_at,
			   k.id_kategori, k.nama_kategori, k.slug
		FROM ARTIKEL a
		LEFT JOIN KATEGORI_KONTEN k ON a.id_kategori = k.id_kategori
		ORDER BY a.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query artikel: %w", err)
	}
	defer rows.Close()

	var list []*models.Artikel
	for rows.Next() {
		var a models.Artikel
		var thumbnail, tanggalTerbit sql.NullString
		var createdAt, updatedAt sql.NullTime
		var k models.KategoriKonten
		var kID sql.NullInt64
		var kNama, kSlug sql.NullString

		err := rows.Scan(
			&a.ID, &a.IDPengurus, &a.IDKategori, &a.Judul, &a.Slug, &a.KontenHTML, &thumbnail, &a.MetaDeskripsi, &a.StatusPublikasi, &tanggalTerbit, &createdAt, &updatedAt,
			&kID, &kNama, &kSlug,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artikel: %w", err)
		}

		if thumbnail.Valid {
			a.GambarThumbnail = &thumbnail.String
		}
		if tanggalTerbit.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", tanggalTerbit.String)
			a.TanggalTerbit = &t
		}
		if createdAt.Valid {
			a.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			a.UpdatedAt = updatedAt.Time
		}
		if kID.Valid {
			k.ID = int(kID.Int64)
		}
		if kNama.Valid {
			k.NamaKategori = kNama.String
		}
		if kSlug.Valid {
			k.Slug = kSlug.String
		}
		a.Kategori = &k

		list = append(list, &a)
	}

	return list, nil
}

func (r *ArtikelRepository) GetPublished(limit int) ([]*models.Artikel, error) {
	query := `
		SELECT a.id_artikel, a.id_pengurus, a.id_kategori, a.judul, a.slug, a.konten_html_markdown, a.gambar_thumbnail_url, a.meta_deskripsi, a.status_publikasi, a.tanggal_terbit, a.created_at, a.updated_at,
			   k.id_kategori, k.nama_kategori, k.slug
		FROM ARTIKEL a
		LEFT JOIN KATEGORI_KONTEN k ON a.id_kategori = k.id_kategori
		WHERE a.status_publikasi = 'Terbit'
		ORDER BY a.tanggal_terbit DESC, a.created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query published artikel: %w", err)
	}
	defer rows.Close()

	var list []*models.Artikel
	for rows.Next() {
		var a models.Artikel
		var thumbnail, tanggalTerbit sql.NullString
		var createdAt, updatedAt sql.NullTime
		var k models.KategoriKonten
		var kID sql.NullInt64
		var kNama, kSlug sql.NullString

		err := rows.Scan(
			&a.ID, &a.IDPengurus, &a.IDKategori, &a.Judul, &a.Slug, &a.KontenHTML, &thumbnail, &a.MetaDeskripsi, &a.StatusPublikasi, &tanggalTerbit, &createdAt, &updatedAt,
			&kID, &kNama, &kSlug,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artikel: %w", err)
		}

		if thumbnail.Valid && thumbnail.String != "" {
			a.GambarThumbnail = &thumbnail.String
		}
		if tanggalTerbit.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", tanggalTerbit.String)
			a.TanggalTerbit = &t
		}
		if createdAt.Valid {
			a.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			a.UpdatedAt = updatedAt.Time
		}
		if kID.Valid {
			k.ID = int(kID.Int64)
		}
		if kNama.Valid {
			k.NamaKategori = kNama.String
		}
		if kSlug.Valid {
			k.Slug = kSlug.String
		}
		a.Kategori = &k

		list = append(list, &a)
	}

	return list, nil
}

func (r *ArtikelRepository) Update(a *models.Artikel) error {
	a.UpdatedAt = time.Now()

	query := `
		UPDATE ARTIKEL SET
			id_pengurus = ?, id_kategori = ?, judul = ?, slug = ?, konten_html_markdown = ?, gambar_thumbnail_url = ?, meta_deskripsi = ?, status_publikasi = ?, tanggal_terbit = ?, updated_at = ?
		WHERE id_artikel = ?
	`

	var tanggalTerbit interface{}
	if a.TanggalTerbit != nil {
		tanggalTerbit = *a.TanggalTerbit
	}

	var thumbnail interface{}
	if a.GambarThumbnail != nil {
		thumbnail = *a.GambarThumbnail
	}

	_, err := r.db.Exec(query, a.IDPengurus, a.IDKategori, a.Judul, a.Slug, a.KontenHTML, thumbnail, a.MetaDeskripsi, a.StatusPublikasi, tanggalTerbit, a.UpdatedAt, a.ID)
	if err != nil {
		return fmt.Errorf("failed to update artikel: %w", err)
	}
	return nil
}

func (r *ArtikelRepository) Delete(id string) error {
	query := `DELETE FROM ARTIKEL WHERE id_artikel = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete artikel: %w", err)
	}
	return nil
}

func (r *ArtikelRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM ARTIKEL`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count artikel: %w", err)
	}
	return count, nil
}

func (r *ArtikelRepository) CountByStatus(status models.StatusPublikasi) (int, error) {
	query := `SELECT COUNT(*) FROM ARTIKEL WHERE status_publikasi = ?`
	var count int
	err := r.db.QueryRow(query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count artikel by status: %w", err)
	}
	return count, nil
}
