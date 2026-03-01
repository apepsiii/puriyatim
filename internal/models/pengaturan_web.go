package models

type PengaturanWeb struct {
	ID                   int     `json:"id_pengaturan" db:"id_pengaturan"`
	NamaLembaga          string  `json:"nama_lembaga" db:"nama_lembaga"`
	DeskripsiTentangKami string  `json:"deskripsi_tentang_kami" db:"deskripsi_tentang_kami"`
	LogoURL              *string `json:"logo_url,omitempty" db:"logo_url"`
	HeroImageURL         *string `json:"hero_image_url,omitempty" db:"hero_image_url"`
	NomorWA              string  `json:"nomor_wa" db:"nomor_wa"`
	EmailLembaga         string  `json:"email_lembaga" db:"email_lembaga"`
	AlamatLengkap        string  `json:"alamat_lengkap" db:"alamat_lengkap"`
	LinkInstagram        *string `json:"link_instagram,omitempty" db:"link_instagram"`
	LinkYouTube          *string `json:"link_youtube,omitempty" db:"link_youtube"`
	OverlayGaleriURL     *string `json:"overlay_galeri_url,omitempty" db:"overlay_galeri_url"`
	RekeningBSI          *string `json:"rekening_bsi,omitempty" db:"rekening_bsi"`
	RekeningMandiri      *string `json:"rekening_mandiri,omitempty" db:"rekening_mandiri"`
	NamaPemilikRekening  *string `json:"nama_pemilik_rekening,omitempty" db:"nama_pemilik_rekening"`
}
