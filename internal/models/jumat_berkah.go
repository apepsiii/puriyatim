package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type StatusKegiatan string

const (
	StatusKegiatanDibuka     StatusKegiatan = "Dibuka"
	StatusKegiatanDitutup    StatusKegiatan = "Ditutup"
	StatusKegiatanSelesai    StatusKegiatan = "Selesai"
)

type StatusApproval string

const (
	StatusApprovalMenunggu StatusApproval = "Menunggu"
	StatusApprovalDisetujui StatusApproval = "Disetujui"
	StatusApprovalDitolak  StatusApproval = "Ditolak"
)

type KegiatanJumatBerkah struct {
	ID            string         `json:"id_kegiatan" db:"id_kegiatan"`
	TanggalKegiatan time.Time    `json:"tanggal_kegiatan" db:"tanggal_kegiatan"`
	KuotaMaksimal int            `json:"kuota_maksimal" db:"kuota_maksimal"`
	TotalTerdaftar int           `json:"total_terdaftar" db:"total_terdaftar"`
	StatusKegiatan StatusKegiatan `json:"status_kegiatan" db:"status_kegiatan"`
}

type PendaftarJumatBerkah struct {
ID           string         `json:"id_pendaftaran" db:"id_pendaftaran"`
	IDKegiatan  string         `json:"id_kegiatan" db:"id_kegiatan"`
	IDAnak      string         `json:"id_anak" db:"id_anak"`
	WaktuSubmit time.Time      `json:"waktu_submit" db:"waktu_submit"`
	StatusApproval StatusApproval `json:"status_approval" db:"status_approval"`

// Join fields
Kegiatan     *KegiatanJumatBerkah `json:"kegiatan,omitempty"`
Anak         *AnakAsuh            `json:"anak,omitempty"`
}

// Value implements driver.Valuer interface
func (s StatusKegiatan) Value() (driver.Value, error) {
return string(s), nil
}

// Scan implements sql.Scanner interface
func (s *StatusKegiatan) Scan(value interface{}) error {
if value == nil {
*s = ""
return nil
}
if str, ok := value.(string); ok {
*s = StatusKegiatan(str)
return nil
}
return fmt.Errorf("cannot scan %T into StatusKegiatan", value)
}

// Value implements driver.Valuer interface
func (s StatusApproval) Value() (driver.Value, error) {
return string(s), nil
}

// Scan implements sql.Scanner interface
func (s *StatusApproval) Scan(value interface{}) error {
if value == nil {
*s = ""
return nil
}
if str, ok := value.(string); ok {
*s = StatusApproval(str)
return nil
}
return fmt.Errorf("cannot scan %T into StatusApproval", value)
}