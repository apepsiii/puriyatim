package repository

import (
	"database/sql"
	"fmt"

	"puriyatim-app/internal/models"
)

type PengaturanRepository struct {
	db *sql.DB
}

func NewPengaturanRepository(db *sql.DB) *PengaturanRepository {
	repo := &PengaturanRepository{db: db}
	repo.ensureOverlayColumn()
	repo.ensureTransferColumns()
	return repo
}

func (r *PengaturanRepository) ensureOverlayColumn() {
	query := `ALTER TABLE PENGATURAN_WEB ADD COLUMN overlay_galeri_url TEXT`
	_, _ = r.db.Exec(query)
}

func (r *PengaturanRepository) ensureTransferColumns() {
	_, _ = r.db.Exec(`ALTER TABLE PENGATURAN_WEB ADD COLUMN rekening_bsi TEXT`)
	_, _ = r.db.Exec(`ALTER TABLE PENGATURAN_WEB ADD COLUMN rekening_mandiri TEXT`)
	_, _ = r.db.Exec(`ALTER TABLE PENGATURAN_WEB ADD COLUMN nama_pemilik_rekening TEXT`)
}

func (r *PengaturanRepository) Get() (*models.PengaturanWeb, error) {
	query := `
		SELECT id_pengaturan, nama_lembaga, deskripsi_tentang_kami, logo_url, hero_image_url, nomor_wa, email_lembaga, alamat_lengkap, link_instagram, link_youtube, overlay_galeri_url, rekening_bsi, rekening_mandiri, nama_pemilik_rekening
		FROM PENGATURAN_WEB
		WHERE id_pengaturan = 1
	`

	var p models.PengaturanWeb
	var logoURL, heroURL, ig, yt, overlay, rekeningBSI, rekeningMandiri, namaPemilik sql.NullString
	err := r.db.QueryRow(query).Scan(
		&p.ID, &p.NamaLembaga, &p.DeskripsiTentangKami, &logoURL, &heroURL, &p.NomorWA, &p.EmailLembaga, &p.AlamatLengkap, &ig, &yt, &overlay, &rekeningBSI, &rekeningMandiri, &namaPemilik,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get pengaturan: %w", err)
	}

	if logoURL.Valid {
		p.LogoURL = &logoURL.String
	}
	if heroURL.Valid {
		p.HeroImageURL = &heroURL.String
	}
	if ig.Valid {
		p.LinkInstagram = &ig.String
	}
	if yt.Valid {
		p.LinkYouTube = &yt.String
	}
	if overlay.Valid {
		p.OverlayGaleriURL = &overlay.String
	}
	if rekeningBSI.Valid {
		p.RekeningBSI = &rekeningBSI.String
	}
	if rekeningMandiri.Valid {
		p.RekeningMandiri = &rekeningMandiri.String
	}
	if namaPemilik.Valid {
		p.NamaPemilikRekening = &namaPemilik.String
	}

	return &p, nil
}

func (r *PengaturanRepository) Save(p *models.PengaturanWeb) error {
	query := `
		UPDATE PENGATURAN_WEB
		SET nama_lembaga = ?, deskripsi_tentang_kami = ?, nomor_wa = ?, email_lembaga = ?, alamat_lengkap = ?, link_instagram = ?, link_youtube = ?, overlay_galeri_url = ?, rekening_bsi = ?, rekening_mandiri = ?, nama_pemilik_rekening = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id_pengaturan = 1
	`

	var ig interface{}
	var yt interface{}
	var overlay interface{}
	var rekeningBSI interface{}
	var rekeningMandiri interface{}
	var namaPemilik interface{}
	if p.LinkInstagram != nil {
		ig = *p.LinkInstagram
	}
	if p.LinkYouTube != nil {
		yt = *p.LinkYouTube
	}
	if p.OverlayGaleriURL != nil {
		overlay = *p.OverlayGaleriURL
	}
	if p.RekeningBSI != nil {
		rekeningBSI = *p.RekeningBSI
	}
	if p.RekeningMandiri != nil {
		rekeningMandiri = *p.RekeningMandiri
	}
	if p.NamaPemilikRekening != nil {
		namaPemilik = *p.NamaPemilikRekening
	}

	_, err := r.db.Exec(query, p.NamaLembaga, p.DeskripsiTentangKami, p.NomorWA, p.EmailLembaga, p.AlamatLengkap, ig, yt, overlay, rekeningBSI, rekeningMandiri, namaPemilik)
	if err != nil {
		return fmt.Errorf("failed to save pengaturan: %w", err)
	}

	return nil
}
