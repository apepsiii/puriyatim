package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type TipeDonatur string

const (
	TipeDonaturIndividu TipeDonatur = "Individu"
	TipeDonaturKelompok TipeDonatur = "Kelompok"
	TipeDonaturInstansi  TipeDonatur = "Instansi"
)

type Donatur struct {
	ID             string       `json:"id_donatur" db:"id_donatur"`
	NamaDonatur    string       `json:"nama_donatur" db:"nama_donatur"`
	TipeDonatur    TipeDonatur  `json:"tipe_donatur" db:"tipe_donatur"`
	NoTelepon      string       `json:"no_telepon" db:"no_telepon"`
	Email          *string      `json:"email,omitempty" db:"email"`
	Alamat         *string      `json:"alamat,omitempty" db:"alamat"`
	CatatanKhusus  string       `json:"catatan_khusus" db:"catatan_khusus"`
	TanggalBergabung time.Time   `json:"tanggal_bergabung" db:"tanggal_bergabung"`
}

// Value implements driver.Valuer interface
func (t TipeDonatur) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements sql.Scanner interface
func (t *TipeDonatur) Scan(value interface{}) error {
	if value == nil {
		*t = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*t = TipeDonatur(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into TipeDonatur", value)
}