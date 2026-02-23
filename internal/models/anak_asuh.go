package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type JenisKelamin string

const (
	JenisKelaminLakiLaki  JenisKelamin = "L"
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
	ID                string       `json:"id_anak" db:"id_anak"`
	NIK               *string      `json:"nik,omitempty" db:"nik"`
	NamaLengkap       string       `json:"nama_lengkap" db:"nama_lengkap"`
	NamaPanggilan     string       `json:"nama_panggilan" db:"nama_panggilan"`
	TempatLahir       string       `json:"tempat_lahir" db:"tempat_lahir"`
	TanggalLahir      time.Time    `json:"tanggal_lahir" db:"tanggal_lahir"`
	JenisKelamin      JenisKelamin `json:"jenis_kelamin" db:"jenis_kelamin"`
	AlamatJalan       string       `json:"alamat_jalan" db:"alamat_jalan"`
	RT                string       `json:"rt" db:"rt"`
	RW                string       `json:"rw" db:"rw"`
	DesaKelurahan     string       `json:"desa_kelurahan" db:"desa_kelurahan"`
	Kecamatan         string       `json:"kecamatan" db:"kecamatan"`
	Kota              string       `json:"kota" db:"kota"`
	TanggalMasuk      time.Time    `json:"tanggal_masuk" db:"tanggal_masuk"`
	StatusAnak        StatusAnak   `json:"status_anak" db:"status_anak"`
	StatusAktif       StatusAktif  `json:"status_aktif" db:"status_aktif"`
	NamaWali          string       `json:"nama_wali" db:"nama_wali"`
	KontakWali        string       `json:"kontak_wali" db:"kontak_wali"`
	HubunganWali      string       `json:"hubungan_wali" db:"hubungan_wali"`
	JenjangPendidikan string       `json:"jenjang_pendidikan" db:"jenjang_pendidikan"`
	NamaSekolah       string       `json:"nama_sekolah" db:"nama_sekolah"`
	Kelas             string       `json:"kelas" db:"kelas"`
	KondisiKesehatan  string       `json:"kondisi_kesehatan" db:"kondisi_kesehatan"`
	CatatanKhusus     string       `json:"catatan_khusus" db:"catatan_khusus"`
	FotoProfilURL     *string      `json:"foto_profil_url,omitempty" db:"foto_profil_url"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" db:"updated_at"`
}

func (a *AnakAsuh) GetInitials() string {
	if len(a.NamaLengkap) < 2 {
		return a.NamaLengkap
	}
	parts := strings.Fields(a.NamaLengkap)
	if len(parts) >= 2 {
		return strings.ToUpper(string(parts[0][0]) + string(parts[1][0]))
	}
	return strings.ToUpper(a.NamaLengkap[:2])
}

func (a *AnakAsuh) GetWilayah() string {
	return fmt.Sprintf("RW %s", a.RW)
}

func (a *AnakAsuh) GetDomisili() string {
	return fmt.Sprintf("RT %s / RW %s", a.RT, a.RW)
}

func (a *AnakAsuh) GetStatusStyle() (bg, text, border, dot string) {
	switch a.StatusAnak {
	case StatusAnakYatim:
		return "bg-blue-50", "text-blue-700", "border-blue-100", "bg-blue-500"
	case StatusAnakPiatu:
		return "bg-cyan-50", "text-cyan-700", "border-cyan-100", "bg-cyan-500"
	case StatusAnakYatimPiatu:
		return "bg-rose-50", "text-rose-700", "border-rose-100", "bg-rose-500"
	case StatusAnakDhuafa:
		return "bg-amber-50", "text-amber-700", "border-amber-100", "bg-amber-500"
	default:
		return "bg-gray-50", "text-gray-700", "border-gray-100", "bg-gray-500"
	}
}

func (a *AnakAsuh) GetAvatarStyle() (bg, text string) {
	initials := a.GetInitials()
	colors := []struct{ bg, text string }{
		{"bg-emerald-100", "text-emerald-600"},
		{"bg-purple-100", "text-purple-600"},
		{"bg-indigo-100", "text-indigo-600"},
		{"bg-pink-100", "text-pink-600"},
		{"bg-teal-100", "text-teal-600"},
		{"bg-blue-100", "text-blue-600"},
	}
	sum := 0
	for _, c := range initials {
		sum += int(c)
	}
	style := colors[sum%len(colors)]
	return style.bg, style.text
}

func (a *AnakAsuh) GetJenjangShort() string {
	switch a.JenjangPendidikan {
	case "SD":
		return "SD"
	case "SMP":
		return "SMP"
	case "SMA/SMK":
		return "SMK"
	default:
		return a.JenjangPendidikan
	}
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
