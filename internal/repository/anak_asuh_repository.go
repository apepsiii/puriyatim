package repository

import (
	"database/sql"
	"fmt"
	"time"

	"puriyatim-app/internal/models"

	"github.com/google/uuid"
)

type AnakAsuhRepository struct {
	db *sql.DB
}

func NewAnakAsuhRepository(db *sql.DB) *AnakAsuhRepository {
	return &AnakAsuhRepository{db: db}
}

func (r *AnakAsuhRepository) generateID() string {
	year := time.Now().Year()
	uniqueID := uuid.New().String()[:8]
	return fmt.Sprintf("P-%d-%s", year, uniqueID)
}

func (r *AnakAsuhRepository) Create(anak *models.AnakAsuh) error {
	if anak.ID == "" {
		anak.ID = r.generateID()
	}

	now := time.Now()
	anak.CreatedAt = now
	anak.UpdatedAt = now

	query := `
		INSERT INTO ANAK_ASUH (
			id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
			jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan, kota,
			tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			kondisi_kesehatan, catatan_khusus, foto_profil_url, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var nik, fotoProfilURL interface{}
	if anak.NIK != nil {
		nik = *anak.NIK
	}
	if anak.FotoProfilURL != nil {
		fotoProfilURL = *anak.FotoProfilURL
	}

	_, err := r.db.Exec(query,
		anak.ID, nik, anak.NamaLengkap, anak.NamaPanggilan, anak.TempatLahir, anak.TanggalLahir,
		anak.JenisKelamin, anak.AlamatJalan, anak.RT, anak.RW, anak.DesaKelurahan, anak.Kecamatan, anak.Kota,
		anak.TanggalMasuk, anak.StatusAnak, anak.StatusAktif, anak.NamaWali, anak.KontakWali,
		anak.HubunganWali, anak.JenjangPendidikan, anak.NamaSekolah, anak.Kelas,
		anak.KondisiKesehatan, anak.CatatanKhusus, fotoProfilURL, anak.CreatedAt, anak.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create anak asuh: %w", err)
	}

	return nil
}

func (r *AnakAsuhRepository) GetByID(id string) (*models.AnakAsuh, error) {
	query := `
		SELECT id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
			   jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan, kota,
			   tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			   hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			   kondisi_kesehatan, catatan_khusus, foto_profil_url, created_at, updated_at
		FROM ANAK_ASUH
		WHERE id_anak = ?
	`

	var anak models.AnakAsuh
	var nik, fotoProfilURL sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&anak.ID, &nik, &anak.NamaLengkap, &anak.NamaPanggilan, &anak.TempatLahir, &anak.TanggalLahir,
		&anak.JenisKelamin, &anak.AlamatJalan, &anak.RT, &anak.RW, &anak.DesaKelurahan, &anak.Kecamatan, &anak.Kota,
		&anak.TanggalMasuk, &anak.StatusAnak, &anak.StatusAktif, &anak.NamaWali, &anak.KontakWali,
		&anak.HubunganWali, &anak.JenjangPendidikan, &anak.NamaSekolah, &anak.Kelas,
		&anak.KondisiKesehatan, &anak.CatatanKhusus, &fotoProfilURL, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("anak asuh with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get anak asuh: %w", err)
	}

	if nik.Valid {
		anak.NIK = &nik.String
	}
	if fotoProfilURL.Valid {
		anak.FotoProfilURL = &fotoProfilURL.String
	}
	if createdAt.Valid {
		anak.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		anak.UpdatedAt = updatedAt.Time
	}

	return &anak, nil
}

func (r *AnakAsuhRepository) GetAll() ([]*models.AnakAsuh, error) {
	query := `
		SELECT id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
			   jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan, kota,
			   tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			   hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			   kondisi_kesehatan, catatan_khusus, foto_profil_url, created_at, updated_at
		FROM ANAK_ASUH
		ORDER BY nama_lengkap
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query anak asuh: %w", err)
	}
	defer rows.Close()

	var anakList []*models.AnakAsuh

	for rows.Next() {
		var anak models.AnakAsuh
		var nik, fotoProfilURL sql.NullString
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&anak.ID, &nik, &anak.NamaLengkap, &anak.NamaPanggilan, &anak.TempatLahir, &anak.TanggalLahir,
			&anak.JenisKelamin, &anak.AlamatJalan, &anak.RT, &anak.RW, &anak.DesaKelurahan, &anak.Kecamatan, &anak.Kota,
			&anak.TanggalMasuk, &anak.StatusAnak, &anak.StatusAktif, &anak.NamaWali, &anak.KontakWali,
			&anak.HubunganWali, &anak.JenjangPendidikan, &anak.NamaSekolah, &anak.Kelas,
			&anak.KondisiKesehatan, &anak.CatatanKhusus, &fotoProfilURL, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan anak asuh: %w", err)
		}

		if nik.Valid {
			anak.NIK = &nik.String
		}
		if fotoProfilURL.Valid {
			anak.FotoProfilURL = &fotoProfilURL.String
		}
		if createdAt.Valid {
			anak.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			anak.UpdatedAt = updatedAt.Time
		}

		anakList = append(anakList, &anak)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anak asuh rows: %w", err)
	}

	return anakList, nil
}

func (r *AnakAsuhRepository) GetByRT(rt string) ([]*models.AnakAsuh, error) {
	query := `
		SELECT id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
			   jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan, kota,
			   tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			   hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			   kondisi_kesehatan, catatan_khusus, foto_profil_url, created_at, updated_at
		FROM ANAK_ASUH
		WHERE rt = ? AND status_aktif = 'Aktif'
		ORDER BY nama_lengkap
	`

	rows, err := r.db.Query(query, rt)
	if err != nil {
		return nil, fmt.Errorf("failed to query anak asuh by RT: %w", err)
	}
	defer rows.Close()

	var anakList []*models.AnakAsuh

	for rows.Next() {
		var anak models.AnakAsuh
		var nik, fotoProfilURL sql.NullString
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&anak.ID, &nik, &anak.NamaLengkap, &anak.NamaPanggilan, &anak.TempatLahir, &anak.TanggalLahir,
			&anak.JenisKelamin, &anak.AlamatJalan, &anak.RT, &anak.RW, &anak.DesaKelurahan, &anak.Kecamatan, &anak.Kota,
			&anak.TanggalMasuk, &anak.StatusAnak, &anak.StatusAktif, &anak.NamaWali, &anak.KontakWali,
			&anak.HubunganWali, &anak.JenjangPendidikan, &anak.NamaSekolah, &anak.Kelas,
			&anak.KondisiKesehatan, &anak.CatatanKhusus, &fotoProfilURL, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan anak asuh: %w", err)
		}

		if nik.Valid {
			anak.NIK = &nik.String
		}
		if fotoProfilURL.Valid {
			anak.FotoProfilURL = &fotoProfilURL.String
		}
		if createdAt.Valid {
			anak.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			anak.UpdatedAt = updatedAt.Time
		}

		anakList = append(anakList, &anak)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anak asuh rows: %w", err)
	}

	return anakList, nil
}

func (r *AnakAsuhRepository) GetByRTRW(rt, rw string) ([]*models.AnakAsuh, error) {
	query := `
		SELECT id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
			   jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan, kota,
			   tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			   hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			   kondisi_kesehatan, catatan_khusus, foto_profil_url, created_at, updated_at
		FROM ANAK_ASUH
		WHERE rt = ? AND rw = ? AND status_aktif = 'Aktif'
		ORDER BY nama_lengkap
	`

	rows, err := r.db.Query(query, rt, rw)
	if err != nil {
		return nil, fmt.Errorf("failed to query anak asuh by RT/RW: %w", err)
	}
	defer rows.Close()

	var anakList []*models.AnakAsuh

	for rows.Next() {
		var anak models.AnakAsuh
		var nik, fotoProfilURL sql.NullString
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&anak.ID, &nik, &anak.NamaLengkap, &anak.NamaPanggilan, &anak.TempatLahir, &anak.TanggalLahir,
			&anak.JenisKelamin, &anak.AlamatJalan, &anak.RT, &anak.RW, &anak.DesaKelurahan, &anak.Kecamatan, &anak.Kota,
			&anak.TanggalMasuk, &anak.StatusAnak, &anak.StatusAktif, &anak.NamaWali, &anak.KontakWali,
			&anak.HubunganWali, &anak.JenjangPendidikan, &anak.NamaSekolah, &anak.Kelas,
			&anak.KondisiKesehatan, &anak.CatatanKhusus, &fotoProfilURL, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan anak asuh: %w", err)
		}

		if nik.Valid {
			anak.NIK = &nik.String
		}
		if fotoProfilURL.Valid {
			anak.FotoProfilURL = &fotoProfilURL.String
		}
		if createdAt.Valid {
			anak.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			anak.UpdatedAt = updatedAt.Time
		}

		anakList = append(anakList, &anak)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anak asuh rows: %w", err)
	}

	return anakList, nil
}

func (r *AnakAsuhRepository) Update(anak *models.AnakAsuh) error {
	anak.UpdatedAt = time.Now()

	query := `
		UPDATE ANAK_ASUH SET
			nik = ?, nama_lengkap = ?, nama_panggilan = ?, tempat_lahir = ?, tanggal_lahir = ?,
			jenis_kelamin = ?, alamat_jalan = ?, rt = ?, rw = ?, desa_kelurahan = ?, kecamatan = ?, kota = ?,
			tanggal_masuk = ?, status_anak = ?, status_aktif = ?, nama_wali = ?, kontak_wali = ?,
			hubungan_wali = ?, jenjang_pendidikan = ?, nama_sekolah = ?, kelas = ?,
			kondisi_kesehatan = ?, catatan_khusus = ?, foto_profil_url = ?, updated_at = ?
		WHERE id_anak = ?
	`

	var nik, fotoProfilURL interface{}
	if anak.NIK != nil {
		nik = *anak.NIK
	}
	if anak.FotoProfilURL != nil {
		fotoProfilURL = *anak.FotoProfilURL
	}

	_, err := r.db.Exec(query,
		nik, anak.NamaLengkap, anak.NamaPanggilan, anak.TempatLahir, anak.TanggalLahir,
		anak.JenisKelamin, anak.AlamatJalan, anak.RT, anak.RW, anak.DesaKelurahan, anak.Kecamatan, anak.Kota,
		anak.TanggalMasuk, anak.StatusAnak, anak.StatusAktif, anak.NamaWali, anak.KontakWali,
		anak.HubunganWali, anak.JenjangPendidikan, anak.NamaSekolah, anak.Kelas,
		anak.KondisiKesehatan, anak.CatatanKhusus, fotoProfilURL, anak.UpdatedAt, anak.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update anak asuh: %w", err)
	}

	return nil
}

func (r *AnakAsuhRepository) Delete(id string) error {
	query := `DELETE FROM ANAK_ASUH WHERE id_anak = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete anak asuh: %w", err)
	}

	return nil
}

func (r *AnakAsuhRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM ANAK_ASUH WHERE status_aktif = 'Aktif'`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count anak asuh: %w", err)
	}

	return count, nil
}

func (r *AnakAsuhRepository) CountAll() (int, error) {
	query := `SELECT COUNT(*) FROM ANAK_ASUH`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count all anak asuh: %w", err)
	}

	return count, nil
}
