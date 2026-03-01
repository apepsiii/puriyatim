package models

import "time"

type GaleriFoto struct {
	ID               string    `json:"id_foto" db:"id_foto"`
	Judul            string    `json:"judul" db:"judul"`
	Deskripsi        string    `json:"deskripsi" db:"deskripsi"`
	GambarAsliURL    string    `json:"gambar_asli_url" db:"gambar_asli_url"`
	GambarOverlayURL string    `json:"gambar_overlay_url" db:"gambar_overlay_url"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
