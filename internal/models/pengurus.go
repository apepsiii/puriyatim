package models

import (
	"database/sql/driver"
	"fmt"
)

type PeranPengurus string

const (
	PeranSuperadmin   PeranPengurus = "Superadmin"
	PeranKeuangan     PeranPengurus = "Keuangan"
	PeranPenulisBerita PeranPengurus = "Penulis Berita"
)

type StatusPengurus string

const (
	StatusPengurusAktif   StatusPengurus = "Aktif"
	StatusPengurusNonaktif StatusPengurus = "Nonaktif"
)

type Pengurus struct {
	ID           string         `json:"id_pengurus" db:"id_pengurus"`
	NamaLengkap  string         `json:"nama_lengkap" db:"nama_lengkap"`
	Email        string         `json:"email" db:"email"`
	PasswordHash string         `json:"-" db:"password_hash"` // Hidden from JSON
	Peran        PeranPengurus  `json:"peran" db:"peran"`
	Status       StatusPengurus `json:"status" db:"status"`
}

// Value implements driver.Valuer interface
func (p PeranPengurus) Value() (driver.Value, error) {
	return string(p), nil
}

// Scan implements sql.Scanner interface
func (p *PeranPengurus) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*p = PeranPengurus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into PeranPengurus", value)
}

// Value implements driver.Valuer interface
func (s StatusPengurus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements sql.Scanner interface
func (s *StatusPengurus) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*s = StatusPengurus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into StatusPengurus", value)
}