package repository

import (
	"database/sql"
	"fmt"
	"puriyatim-app/internal/models"
)

type AnakAsuhRepository struct {
	db *sql.DB
}

func NewAnakAsuhRepository(db *sql.DB) *AnakAsuhRepository {
	return &AnakAsuhRepository{db: db}
}

func (r *AnakAsuhRepository) Create(anak *models.AnakAsuh) error {
	query := `
		INSERT INTO ANAK_ASUH (
			id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
			jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan,
			tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			kondisi_kesehatan, catatan_khusus, foto_profil_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query,
		anak.ID, anak.NIK, anak.NamaLengkap, anak.NamaPanggilan, anak.TempatLahir, anak.TanggalLahir,
		anak.JenisKelamin, anak.AlamatJalan, anak.RT, anak.RW, anak.DesaKelurahan, anak.Kecamatan,
		anak.TanggalMasuk, anak.StatusAnak, anak.StatusAktif, anak.NamaWali, anak.KontakWali,
		anak.HubunganWali, anak.JenjangPendidikan, anak.NamaSekolah, anak.Kelas,
		anak.KondisiKesehatan, anak.CatatanKhusus, anak.FotoProfilURL,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create anak asuh: %w", err)
	}
	
	return nil
}

func (r *AnakAsuhRepository) GetByID(id string) (*models.AnakAsuh, error) {
	query := `
		SELECT id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
			   jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan,
			   tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			   hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			   kondisi_kesehatan, catatan_khusus, foto_profil_url
		FROM ANAK_ASUH
		WHERE id_anak = ?
	`
	
	var anak models.AnakAsuh
	var nik, fotoProfilURL sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&anak.ID, &nik, &anak.NamaLengkap, &anak.NamaPanggilan, &anak.TempatLahir, &anak.TanggalLahir,
		&anak.JenisKelamin, &anak.AlamatJalan, &anak.RT, &anak.RW, &anak.DesaKelurahan, &anak.Kecamatan,
		&anak.TanggalMasuk, &anak.StatusAnak, &anak.StatusAktif, &anak.NamaWali, &anak.KontakWali,
		&anak.HubunganWali, &anak.JenjangPendidikan, &anak.NamaSekolah, &anak.Kelas,
		&anak.KondisiKesehatan, &anak.CatatanKhusus, &fotoProfilURL,
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
	
	return &anak, nil
}

func (r *AnakAsuhRepository) GetAll() ([]*models.AnakAsuh, error) {
	query := `
		SELECT id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
			   jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan,
			   tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			   hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			   kondisi_kesehatan, catatan_khusus, foto_profil_url
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
		
		err := rows.Scan(
			&anak.ID, &nik, &anak.NamaLengkap, &anak.NamaPanggilan, &anak.TempatLahir, &anak.TanggalLahir,
			&anak.JenisKelamin, &anak.AlamatJalan, &anak.RT, &anak.RW, &anak.DesaKelurahan, &anak.Kecamatan,
			&anak.TanggalMasuk, &anak.StatusAnak, &anak.StatusAktif, &anak.NamaWali, &anak.KontakWali,
			&anak.HubunganWali, &anak.JenjangPendidikan, &anak.NamaSekolah, &anak.Kelas,
			&anak.KondisiKesehatan, &anak.CatatanKhusus, &fotoProfilURL,
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
			   jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan,
			   tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
			   hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
			   kondisi_kesehatan, catatan_khusus, foto_profil_url
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
		
		err := rows.Scan(
			&anak.ID, &nik, &anak.NamaLengkap, &anak.NamaPanggilan, &anak.TempatLahir, &anak.TanggalLahir,
			&anak.JenisKelamin, &anak.AlamatJalan, &anak.RT, &anak.RW, &anak.DesaKelurahan, &anak.Kecamatan,
			&anak.TanggalMasuk, &anak.StatusAnak, &anak.StatusAktif, &anak.NamaWali, &anak.KontakWali,
			&anak.HubunganWali, &anak.JenjangPendidikan, &anak.NamaSekolah, &anak.Kelas,
			&anak.KondisiKesehatan, &anak.CatatanKhusus, &fotoProfilURL,
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
		
		anakList = append(anakList, &anak)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anak asuh rows: %w", err)
	}
	
	return anakList, nil
}

func (r *AnakAsuhRepository) Update(anak *models.AnakAsuh) error {
	query := `
		UPDATE ANAK_ASUH SET
			nik = ?, nama_lengkap = ?, nama_panggilan = ?, tempat_lahir = ?, tanggal_lahir = ?,
			jenis_kelamin = ?, alamat_jalan = ?, rt = ?, rw = ?, desa_kelurahan = ?, kecamatan = ?,
			tanggal_masuk = ?, status_anak = ?, status_aktif = ?, nama_wali = ?, kontak_wali = ?,
			hubungan_wali = ?, jenjang_pendidikan = ?, nama_sekolah = ?, kelas = ?,
			kondisi_kesehatan = ?, catatan_khusus = ?, foto_profil_url = ?
		WHERE id_anak = ?
	`
	
	_, err := r.db.Exec(query,
		anak.NIK, anak.NamaLengkap, anak.NamaPanggilan, anak.TempatLahir, anak.TanggalLahir,
		anak.JenisKelamin, anak.AlamatJalan, anak.RT, anak.RW, anak.DesaKelurahan, anak.Kecamatan,
		anak.TanggalMasuk, anak.StatusAnak, anak.StatusAktif, anak.NamaWali, anak.KontakWali,
		anak.HubunganWali, anak.JenjangPendidikan, anak.NamaSekolah, anak.Kelas,
		anak.KondisiKesehatan, anak.CatatanKhusus, anak.FotoProfilURL, anak.ID,
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