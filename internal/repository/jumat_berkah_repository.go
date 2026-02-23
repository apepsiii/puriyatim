package repository

import (
	"database/sql"
	"fmt"
	"time"

	"puriyatim-app/internal/models"

	"github.com/google/uuid"
)

type JumatBerkahRepository struct {
	db *sql.DB
}

func NewJumatBerkahRepository(db *sql.DB) *JumatBerkahRepository {
	return &JumatBerkahRepository{db: db}
}

func (r *JumatBerkahRepository) generateID() string {
	return uuid.New().String()[:8]
}

func (r *JumatBerkahRepository) CreateKegiatan(k *models.KegiatanJumatBerkah) error {
	if k.ID == "" {
		k.ID = r.generateID()
	}
	now := time.Now()
	k.CreatedAt = now
	k.UpdatedAt = now
	if k.StatusKegiatan == "" {
		k.StatusKegiatan = models.StatusKegiatanDibuka
	}

	query := `
		INSERT INTO KEGIATAN_JUMAT_BERKAH (id_kegiatan, tanggal_kegiatan, kuota_maksimal, total_terdaftar, status_kegiatan, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, k.ID, k.TanggalKegiatan, k.KuotaMaksimal, k.TotalTerdaftar, k.StatusKegiatan, k.CreatedAt, k.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create kegiatan: %w", err)
	}
	return nil
}

func (r *JumatBerkahRepository) GetKegiatanByID(id string) (*models.KegiatanJumatBerkah, error) {
	query := `
		SELECT id_kegiatan, tanggal_kegiatan, kuota_maksimal, total_terdaftar, status_kegiatan, created_at, updated_at
		FROM KEGIATAN_JUMAT_BERKAH WHERE id_kegiatan = ?
	`

	var k models.KegiatanJumatBerkah
	err := r.db.QueryRow(query, id).Scan(
		&k.ID, &k.TanggalKegiatan, &k.KuotaMaksimal, &k.TotalTerdaftar, &k.StatusKegiatan, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("kegiatan not found")
		}
		return nil, fmt.Errorf("failed to get kegiatan: %w", err)
	}
	return &k, nil
}

func (r *JumatBerkahRepository) GetCurrentKegiatan() (*models.KegiatanJumatBerkah, error) {
	query := `
		SELECT id_kegiatan, tanggal_kegiatan, kuota_maksimal, total_terdaftar, status_kegiatan, created_at, updated_at
		FROM KEGIATAN_JUMAT_BERKAH 
		WHERE status_kegiatan = 'Dibuka' AND tanggal_kegiatan >= date('now')
		ORDER BY tanggal_kegiatan ASC
		LIMIT 1
	`

	var k models.KegiatanJumatBerkah
	err := r.db.QueryRow(query).Scan(
		&k.ID, &k.TanggalKegiatan, &k.KuotaMaksimal, &k.TotalTerdaftar, &k.StatusKegiatan, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get current kegiatan: %w", err)
	}
	return &k, nil
}

func (r *JumatBerkahRepository) UpdateKegiatan(k *models.KegiatanJumatBerkah) error {
	k.UpdatedAt = time.Now()

	query := `
		UPDATE KEGIATAN_JUMAT_BERKAH SET
			tanggal_kegiatan = ?, kuota_maksimal = ?, total_terdaftar = ?, status_kegiatan = ?, updated_at = ?
		WHERE id_kegiatan = ?
	`

	_, err := r.db.Exec(query, k.TanggalKegiatan, k.KuotaMaksimal, k.TotalTerdaftar, k.StatusKegiatan, k.UpdatedAt, k.ID)
	if err != nil {
		return fmt.Errorf("failed to update kegiatan: %w", err)
	}
	return nil
}

func (r *JumatBerkahRepository) CreatePendaftar(p *models.PendaftarJumatBerkah) error {
	if p.ID == "" {
		p.ID = r.generateID()
	}
	now := time.Now()
	p.WaktuSubmit = now
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.StatusApproval == "" {
		p.StatusApproval = models.StatusApprovalMenunggu
	}

	query := `
		INSERT INTO PENDAFTAR_JUMAT_BERKAH (id_pendaftaran, id_kegiatan, id_anak, waktu_submit, status_approval, catatan, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, p.ID, p.IDKegiatan, p.IDAnak, p.WaktuSubmit, p.StatusApproval, p.Catatan, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create pendaftar: %w", err)
	}
	return nil
}

func (r *JumatBerkahRepository) GetPendaftarByID(id string) (*models.PendaftarJumatBerkah, error) {
	query := `
		SELECT id_pendaftaran, id_kegiatan, id_anak, waktu_submit, status_approval, catatan, created_at, updated_at
		FROM PENDAFTAR_JUMAT_BERKAH WHERE id_pendaftaran = ?
	`

	var p models.PendaftarJumatBerkah
	var catatan sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.IDKegiatan, &p.IDAnak, &p.WaktuSubmit, &p.StatusApproval, &catatan, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pendaftar not found")
		}
		return nil, fmt.Errorf("failed to get pendaftar: %w", err)
	}

	if catatan.Valid {
		p.Catatan = catatan.String
	}
	return &p, nil
}

func (r *JumatBerkahRepository) GetAllPendaftar() ([]*models.PendaftarJumatBerkah, error) {
	query := `
		SELECT p.id_pendaftaran, p.id_kegiatan, p.id_anak, p.waktu_submit, p.status_approval, p.catatan, p.created_at, p.updated_at,
			   a.nama_lengkap, a.nama_panggilan, a.jenjang_pendidikan, a.status_anak, a.rt, a.rw
		FROM PENDAFTAR_JUMAT_BERKAH p
		LEFT JOIN ANAK_ASUH a ON p.id_anak = a.id_anak
		ORDER BY p.waktu_submit DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pendaftar: %w", err)
	}
	defer rows.Close()

	var list []*models.PendaftarJumatBerkah
	for rows.Next() {
		var p models.PendaftarJumatBerkah
		var catatan sql.NullString
		var anak models.AnakAsuh
		var namaPanggilan, jenjang, statusAnak, rt, rw sql.NullString

		err := rows.Scan(
			&p.ID, &p.IDKegiatan, &p.IDAnak, &p.WaktuSubmit, &p.StatusApproval, &catatan, &p.CreatedAt, &p.UpdatedAt,
			&anak.NamaLengkap, &namaPanggilan, &jenjang, &statusAnak, &rt, &rw,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pendaftar: %w", err)
		}

		if catatan.Valid {
			p.Catatan = catatan.String
		}

		anak.ID = p.IDAnak
		if namaPanggilan.Valid {
			anak.NamaPanggilan = namaPanggilan.String
		}
		if jenjang.Valid {
			anak.JenjangPendidikan = jenjang.String
		}
		if statusAnak.Valid {
			anak.StatusAnak = models.StatusAnak(statusAnak.String)
		}
		if rt.Valid {
			anak.RT = rt.String
		}
		if rw.Valid {
			anak.RW = rw.String
		}

		p.Anak = &anak
		list = append(list, &p)
	}

	return list, nil
}

func (r *JumatBerkahRepository) GetPendaftarByKegiatan(kegiatanID string) ([]*models.PendaftarJumatBerkah, error) {
	query := `
		SELECT p.id_pendaftaran, p.id_kegiatan, p.id_anak, p.waktu_submit, p.status_approval, p.catatan, p.created_at, p.updated_at,
			   a.nama_lengkap, a.nama_panggilan, a.jenjang_pendidikan, a.status_anak, a.rt, a.rw
		FROM PENDAFTAR_JUMAT_BERKAH p
		LEFT JOIN ANAK_ASUH a ON p.id_anak = a.id_anak
		WHERE p.id_kegiatan = ?
		ORDER BY p.waktu_submit DESC
	`

	rows, err := r.db.Query(query, kegiatanID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pendaftar by kegiatan: %w", err)
	}
	defer rows.Close()

	var list []*models.PendaftarJumatBerkah
	for rows.Next() {
		var p models.PendaftarJumatBerkah
		var catatan sql.NullString
		var anak models.AnakAsuh
		var namaPanggilan, jenjang, statusAnak, rt, rw sql.NullString

		err := rows.Scan(
			&p.ID, &p.IDKegiatan, &p.IDAnak, &p.WaktuSubmit, &p.StatusApproval, &catatan, &p.CreatedAt, &p.UpdatedAt,
			&anak.NamaLengkap, &namaPanggilan, &jenjang, &statusAnak, &rt, &rw,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pendaftar: %w", err)
		}

		if catatan.Valid {
			p.Catatan = catatan.String
		}

		anak.ID = p.IDAnak
		if namaPanggilan.Valid {
			anak.NamaPanggilan = namaPanggilan.String
		}
		if jenjang.Valid {
			anak.JenjangPendidikan = jenjang.String
		}
		if statusAnak.Valid {
			anak.StatusAnak = models.StatusAnak(statusAnak.String)
		}
		if rt.Valid {
			anak.RT = rt.String
		}
		if rw.Valid {
			anak.RW = rw.String
		}

		p.Anak = &anak
		list = append(list, &p)
	}

	return list, nil
}

func (r *JumatBerkahRepository) GetPendaftarByStatus(kegiatanID string, status models.StatusApproval) ([]*models.PendaftarJumatBerkah, error) {
	query := `
		SELECT p.id_pendaftaran, p.id_kegiatan, p.id_anak, p.waktu_submit, p.status_approval, p.catatan, p.created_at, p.updated_at,
			   a.nama_lengkap, a.nama_panggilan, a.jenjang_pendidikan, a.status_anak, a.rt, a.rw
		FROM PENDAFTAR_JUMAT_BERKAH p
		LEFT JOIN ANAK_ASUH a ON p.id_anak = a.id_anak
		WHERE p.id_kegiatan = ? AND p.status_approval = ?
		ORDER BY p.waktu_submit DESC
	`

	rows, err := r.db.Query(query, kegiatanID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query pendaftar by status: %w", err)
	}
	defer rows.Close()

	var list []*models.PendaftarJumatBerkah
	for rows.Next() {
		var p models.PendaftarJumatBerkah
		var catatan sql.NullString
		var anak models.AnakAsuh
		var namaPanggilan, jenjang, statusAnak, rt, rw sql.NullString

		err := rows.Scan(
			&p.ID, &p.IDKegiatan, &p.IDAnak, &p.WaktuSubmit, &p.StatusApproval, &catatan, &p.CreatedAt, &p.UpdatedAt,
			&anak.NamaLengkap, &namaPanggilan, &jenjang, &statusAnak, &rt, &rw,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pendaftar: %w", err)
		}

		if catatan.Valid {
			p.Catatan = catatan.String
		}

		anak.ID = p.IDAnak
		if namaPanggilan.Valid {
			anak.NamaPanggilan = namaPanggilan.String
		}
		if jenjang.Valid {
			anak.JenjangPendidikan = jenjang.String
		}
		if statusAnak.Valid {
			anak.StatusAnak = models.StatusAnak(statusAnak.String)
		}
		if rt.Valid {
			anak.RT = rt.String
		}
		if rw.Valid {
			anak.RW = rw.String
		}

		p.Anak = &anak
		list = append(list, &p)
	}

	return list, nil
}

func (r *JumatBerkahRepository) UpdatePendaftarStatus(id string, status models.StatusApproval) error {
	now := time.Now()

	query := `UPDATE PENDAFTAR_JUMAT_BERKAH SET status_approval = ?, updated_at = ? WHERE id_pendaftaran = ?`

	_, err := r.db.Exec(query, status, now, id)
	if err != nil {
		return fmt.Errorf("failed to update pendaftar status: %w", err)
	}
	return nil
}

func (r *JumatBerkahRepository) UpdateMultiplePendaftarStatus(ids []string, status models.StatusApproval) (int, error) {
	now := time.Now()
	count := 0

	for _, id := range ids {
		query := `UPDATE PENDAFTAR_JUMAT_BERKAH SET status_approval = ?, updated_at = ? WHERE id_pendaftaran = ? AND status_approval = 'Menunggu'`
		result, err := r.db.Exec(query, status, now, id)
		if err != nil {
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			count++
		}
	}

	return count, nil
}

func (r *JumatBerkahRepository) CountPendingByKegiatan(kegiatanID string) (int, error) {
	query := `SELECT COUNT(*) FROM PENDAFTAR_JUMAT_BERKAH WHERE id_kegiatan = ? AND status_approval = 'Menunggu'`

	var count int
	err := r.db.QueryRow(query, kegiatanID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending: %w", err)
	}
	return count, nil
}

func (r *JumatBerkahRepository) CountApprovedByKegiatan(kegiatanID string) (int, error) {
	query := `SELECT COUNT(*) FROM PENDAFTAR_JUMAT_BERKAH WHERE id_kegiatan = ? AND status_approval = 'Disetujui'`

	var count int
	err := r.db.QueryRow(query, kegiatanID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count approved: %w", err)
	}
	return count, nil
}

func (r *JumatBerkahRepository) IsAnakRegistered(kegiatanID, anakID string) (bool, error) {
	query := `SELECT COUNT(*) FROM PENDAFTAR_JUMAT_BERKAH WHERE id_kegiatan = ? AND id_anak = ? AND status_approval IN ('Menunggu', 'Disetujui')`

	var count int
	err := r.db.QueryRow(query, kegiatanID, anakID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check registration: %w", err)
	}
	return count > 0, nil
}

func (r *JumatBerkahRepository) CountAllKegiatan() (int, error) {
	query := `SELECT COUNT(*) FROM KEGIATAN_JUMAT_BERKAH`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count kegiatan: %w", err)
	}
	return count, nil
}

func (r *JumatBerkahRepository) CountAllPendaftar() (int, error) {
	query := `SELECT COUNT(*) FROM PENDAFTAR_JUMAT_BERKAH`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pendaftar: %w", err)
	}
	return count, nil
}

func (r *JumatBerkahRepository) GetPendaftarByAnakID(anakID string) ([]*models.PendaftarJumatBerkah, error) {
	query := `
		SELECT p.id_pendaftaran, p.id_kegiatan, p.id_anak, p.waktu_submit, p.status_approval, p.catatan, p.created_at, p.updated_at,
			   k.tanggal_kegiatan, k.status_kegiatan
		FROM PENDAFTAR_JUMAT_BERKAH p
		LEFT JOIN KEGIATAN_JUMAT_BERKAH k ON p.id_kegiatan = k.id_kegiatan
		WHERE p.id_anak = ?
		ORDER BY p.waktu_submit DESC
	`

	rows, err := r.db.Query(query, anakID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pendaftar by anak: %w", err)
	}
	defer rows.Close()

	var list []*models.PendaftarJumatBerkah
	for rows.Next() {
		var p models.PendaftarJumatBerkah
		var catatan sql.NullString
		var kegiatan models.KegiatanJumatBerkah

		err := rows.Scan(
			&p.ID, &p.IDKegiatan, &p.IDAnak, &p.WaktuSubmit, &p.StatusApproval, &catatan, &p.CreatedAt, &p.UpdatedAt,
			&kegiatan.TanggalKegiatan, &kegiatan.StatusKegiatan,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pendaftar: %w", err)
		}

		if catatan.Valid {
			p.Catatan = catatan.String
		}

		kegiatan.ID = p.IDKegiatan
		p.Kegiatan = &kegiatan
		list = append(list, &p)
	}

	return list, nil
}
