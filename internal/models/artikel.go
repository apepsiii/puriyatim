package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type StatusPublikasi string

const (
	StatusPublikasiDraft  StatusPublikasi = "Draft"
	StatusPublikasiTerbit StatusPublikasi = "Terbit"
	StatusPublikasiArsip  StatusPublikasi = "Arsip"
)

type KategoriKonten struct {
	ID           int    `json:"id_kategori" db:"id_kategori"`
	NamaKategori string `json:"nama_kategori" db:"nama_kategori"`
	Slug         string `json:"slug" db:"slug"`
}

type Artikel struct {
	ID               string          `json:"id_artikel" db:"id_artikel"`
	IDPengurus       string          `json:"id_pengurus" db:"id_pengurus"`
	IDKategori       int             `json:"id_kategori" db:"id_kategori"`
	Judul            string          `json:"judul" db:"judul"`
	Slug             string          `json:"slug" db:"slug"`
	KontenHTML       string          `json:"konten_html_markdown" db:"konten_html_markdown"`
	GambarThumbnail  *string         `json:"gambar_thumbnail_url,omitempty" db:"gambar_thumbnail_url"`
	MetaDeskripsi    string          `json:"meta_deskripsi" db:"meta_deskripsi"`
	StatusPublikasi  StatusPublikasi `json:"status_publikasi" db:"status_publikasi"`
	TanggalTerbit    *time.Time      `json:"tanggal_terbit,omitempty" db:"tanggal_terbit"`
	
	// Join fields
	Pengurus         *Pengurus       `json:"pengurus,omitempty"`
	Kategori         *KategoriKonten `json:"kategori,omitempty"`
}

// Value implements driver.Valuer interface
func (s StatusPublikasi) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements sql.Scanner interface
func (s *StatusPublikasi) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*s = StatusPublikasi(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into StatusPublikasi", value)
}