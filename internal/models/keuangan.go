package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type KategoriDana string

const (
	KategoriDanaInfaq   KategoriDana = "Infaq"
	KategoriDanaSedekah KategoriDana = "Sedekah"
	KategoriDanaWakaf   KategoriDana = "Wakaf"
	KategoriDanaZakat   KategoriDana = "Zakat"
	KategoriDanaLainnya KategoriDana = "Lainnya"
)

type PemasukanDonasi struct {
	ID             string       `json:"id_pemasukan" db:"id_pemasukan"`
	NamaDonatur    string       `json:"nama_donatur" db:"nama_donatur"`
	TanggalDonasi  time.Time    `json:"tanggal_donasi" db:"tanggal_donasi"`
	Nominal        float64      `json:"nominal" db:"nominal"`
	KategoriDana   KategoriDana `json:"kategori_dana" db:"kategori_dana"`
	Catatan        string       `json:"catatan" db:"catatan"`
	BuktiTransaksi string       `json:"bukti_transaksi" db:"bukti_transaksi"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
}

type Pengeluaran struct {
	ID                 string    `json:"id_pengeluaran" db:"id_pengeluaran"`
	TanggalPengeluaran time.Time `json:"tanggal_pengeluaran" db:"tanggal_pengeluaran"`
	Nominal            float64   `json:"nominal" db:"nominal"`
	IDAnak             *string   `json:"id_anak,omitempty" db:"id_anak"`
	Keterangan         string    `json:"keterangan" db:"keterangan"`
	BuktiTransaksi     string    `json:"bukti_transaksi" db:"bukti_transaksi"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`

	Anak *AnakAsuh `json:"anak,omitempty"`
}

type KasTransaction struct {
	ID        string
	Tanggal   time.Time
	CreatedAt time.Time
	Deskripsi string
	Kategori  string
	Jumlah    float64
	Type      string
	Donatur   string
	AnakAsuh  string
}

type KeuanganStats struct {
	TotalSaldo           float64
	TotalPemasukan       float64
	TotalPengeluaran     float64
	PemasukanBulanIni    float64
	PengeluaranBulanIni  float64
	PemasukanBulanLalu   float64
	PengeluaranBulanLalu float64
	PemasukanChange      float64
	PengeluaranChange    float64
}

func (k KategoriDana) Value() (driver.Value, error) {
	return string(k), nil
}

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

func (k KategoriDana) GetStyle() (bg, text, border string) {
	switch k {
	case KategoriDanaZakat:
		return "bg-purple-50", "text-purple-700", "border-purple-100"
	case KategoriDanaInfaq:
		return "bg-blue-50", "text-blue-700", "border-blue-100"
	case KategoriDanaSedekah:
		return "bg-emerald-50", "text-emerald-700", "border-emerald-100"
	case KategoriDanaWakaf:
		return "bg-amber-50", "text-amber-700", "border-amber-100"
	default:
		return "bg-gray-50", "text-gray-700", "border-gray-100"
	}
}
