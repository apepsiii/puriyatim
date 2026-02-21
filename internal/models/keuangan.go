package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type KategoriDana string

const (
	KategoriDanaInfaq    KategoriDana = "Infaq"
	KategoriDanaSedekah  KategoriDana = "Sedekah"
	KategoriDanaWakaf    KategoriDana = "Wakaf"
	KategoriDanaZakat    KategoriDana = "Zakat"
	KategoriDanaLainnya  KategoriDana = "Lainnya"
)

type PemasukanDonasi struct {
	ID            string       `json:"id_pemasukan" db:"id_pemasukan"`
	IDDonatur     *string      `json:"id_donatur,omitempty" db:"id_donatur"`
	TanggalDonasi time.Time    `json:"tanggal_donasi" db:"tanggal_donasi"`
	Nominal       float64      `json:"nominal" db:"nominal"`
	KategoriDana  KategoriDana `json:"kategori_dana" db:"kategori_dana"`
	
	// Join fields
	Donatur       *Donatur     `json:"donatur,omitempty"`
}

type Pengeluaran struct {
ID             string     `json:"id_pengeluaran" db:"id_pengeluaran"`
	TanggalPengeluaran time.Time `json:"tanggal_pengeluaran" db:"tanggal_pengeluaran"`
	Nominal       float64    `json:"nominal" db:"nominal"`
	IDAnak        *string    `json:"id_anak,omitempty" db:"id_anak"`
	
	// Join fields
	Anak          *AnakAsuh  `json:"anak,omitempty"`
}

// Value implements driver.Valuer interface
func (k KategoriDana) Value() (driver.Value, error) {
	return string(k), nil
}

// Scan implements sql.Scanner interface
func (k *KategoriDana) Scan(value interface{}) error {
	if value == nil {
		*k = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*k = KategoriDana(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into KategoriDana", value)
}