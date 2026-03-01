package models

// RekeningDonasi menyimpan satu entri rekening bank untuk tujuan transfer donasi.
type RekeningDonasi struct {
	ID             int    `json:"id"              db:"id"`
	NamaBank       string `json:"nama_bank"       db:"nama_bank"`
	LogoBank       string `json:"logo_bank"       db:"logo_bank"` // slug: bsi, mandiri, bri, bni, dst
	NomorRekening  string `json:"nomor_rekening"  db:"nomor_rekening"`
	AtasNama       string `json:"atas_nama"       db:"atas_nama"`
	Urutan         int    `json:"urutan"          db:"urutan"`
	Aktif          bool   `json:"aktif"           db:"aktif"`
}
