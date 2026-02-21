package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type JenisKelamin string

const (
	JenisKelaminLakiLaki JenisKelamin = "L"
	JenisKelaminPerempuan JenisKelamin = "P"
)

type StatusAnak string

const (
	StatusAnakYatim      StatusAnak = "Yatim"
	StatusAnakPiatu      StatusAnak = "Piatu"
	StatusAnakYatimPiatu StatusAnak = "Yatim Piatu"
	StatusAnakDhuafa     StatusAnak = "Dhuafa"
)

type StatusAktif string

const (
	StatusAktifAktif  StatusAktif = "Aktif"
	StatusAktifLulus  StatusAktif = "Lulus"
	StatusAktifKeluar StatusAktif = "Keluar"
)

type AnakAsuh struct {
	ID               string        `json:"id_anak" db:"id_anak"`
	NIK              *string       `json:"nik,omitempty" db:"nik"`
	NamaLengkap      string        `json:"nama_lengkap" db:"nama_lengkap"`
	NamaPanggilan    string        `json:"nama_panggilan" db:"nama_panggilan"`
	TempatLahir      string        `json:"tempat_lahir" db:"tempat_lahir"`
	TanggalLahir     time.Time     `json:"tanggal_lahir" db:"tanggal_lahir"`
	JenisKelamin     JenisKelamin  `json:"jenis_kelamin" db:"jenis_kelamin"`
	AlamatJalan      string        `json:"alamat_jalan" db:"alamat_jalan"`
	RT               string        `json:"rt" db:"rt"`
	RW               string        `json:"rw" db:"rw"`
	DesaKelurahan    string        `json:"desa_kelurahan" db:"desa_kelurahan"`
	Kecamatan        string        `json:"kecamatan" db:"kecamatan"`
	TanggalMasuk     time.Time     `json:"tanggal_masuk" db:"tanggal_masuk"`
	StatusAnak       StatusAnak    `json:"status_anak" db:"status_anak"`
	StatusAktif      StatusAktif   `json:"status_aktif" db:"status_aktif"`
	NamaWali         string        `json:"nama_wali" db:"nama_wali"`
	KontakWali       string        `json:"kontak_wali" db:"kontak_wali"`
	HubunganWali     string        `json:"hubungan_wali" db:"hubungan_wali"`
	JenjangPendidikan string        `json:"jenjang_pendidikan" db:"jenjang_pendidikan"`
	NamaSekolah      string        `json:"nama_sekolah" db:"nama_sekolah"`
	Kelas            string        `json:"kelas" db:"kelas"`
	KondisiKesehatan string        `json:"kondisi_kesehatan" db:"kondisi_kesehatan"`
	CatatanKhusus    string        `json:"catatan_khusus" db:"catatan_khusus"`
	FotoProfilURL    *string       `json:"foto_profil_url,omitempty" db:"foto_profil_url"`
}

// Value implements driver.Valuer interface
func (j JenisKelamin) Value() (driver.Value, error) {
	return string(j), nil
}

// Scan implements sql.Scanner interface
func (j *JenisKelamin) Scan(value interface{}) error {
	if value == nil {
		*j = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*j = JenisKelamin(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into JenisKelamin", value)
}

// Value implements driver.Valuer interface
func (s StatusAnak) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements sql.Scanner interface
func (s *StatusAnak) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*s = StatusAnak(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into StatusAnak", value)
}

// Value implements driver.Valuer interface
func (s StatusAktif) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements sql.Scanner interface
func (s *StatusAktif) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*s = StatusAktif(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into StatusAktif", value)
}