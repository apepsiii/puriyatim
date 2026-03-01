package models

// RekeningDonasi menyimpan satu entri rekening bank untuk tujuan transfer donasi.
type RekeningDonasi struct {
	ID            int    `json:"id"             db:"id"`
	NamaBank      string `json:"nama_bank"      db:"nama_bank"`
	LogoURL       string `json:"logo_url"       db:"logo_bank"` // URL path ke file logo yang di-upload
	NomorRekening string `json:"nomor_rekening" db:"nomor_rekening"`
	AtasNama      string `json:"atas_nama"      db:"atas_nama"`
	Urutan        int    `json:"urutan"         db:"urutan"`
	Aktif         bool   `json:"aktif"          db:"aktif"`
}
