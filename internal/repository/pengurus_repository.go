package repository

import (
	"database/sql"
	"fmt"
	"puriyatim-app/internal/models"
)

type PengurusRepository struct {
	db *sql.DB
}

func NewPengurusRepository(db *sql.DB) *PengurusRepository {
	return &PengurusRepository{db: db}
}

func (r *PengurusRepository) Create(pengurus *models.Pengurus) error {
	query := `
		INSERT INTO PENGURUS (
			id_pengurus, nama_lengkap, email, password_hash, peran, status
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.Exec(query,
		pengurus.ID, pengurus.NamaLengkap, pengurus.Email, pengurus.PasswordHash,
		pengurus.Peran, pengurus.Status,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create pengurus: %w", err)
	}
	
	return nil
}

func (r *PengurusRepository) GetByID(id string) (*models.Pengurus, error) {
	query := `
		SELECT id_pengurus, nama_lengkap, email, password_hash, peran, status
		FROM PENGURUS
		WHERE id_pengurus = ?
	`
	
	var pengurus models.Pengurus
	
	err := r.db.QueryRow(query, id).Scan(
		&pengurus.ID, &pengurus.NamaLengkap, &pengurus.Email, &pengurus.PasswordHash,
		&pengurus.Peran, &pengurus.Status,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pengurus with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get pengurus: %w", err)
	}
	
	return &pengurus, nil
}

func (r *PengurusRepository) GetByEmail(email string) (*models.Pengurus, error) {
	query := `
		SELECT id_pengurus, nama_lengkap, email, password_hash, peran, status
		FROM PENGURUS
		WHERE email = ?
	`
	
	var pengurus models.Pengurus
	
	err := r.db.QueryRow(query, email).Scan(
		&pengurus.ID, &pengurus.NamaLengkap, &pengurus.Email, &pengurus.PasswordHash,
		&pengurus.Peran, &pengurus.Status,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pengurus with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get pengurus by email: %w", err)
	}
	
	return &pengurus, nil
}

func (r *PengurusRepository) GetAll() ([]*models.Pengurus, error) {
	query := `
		SELECT id_pengurus, nama_lengkap, email, password_hash, peran, status
		FROM PENGURUS
		ORDER BY nama_lengkap
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pengurus: %w", err)
	}
	defer rows.Close()
	
	var pengurusList []*models.Pengurus
	
	for rows.Next() {
		var pengurus models.Pengurus
		
		err := rows.Scan(
			&pengurus.ID, &pengurus.NamaLengkap, &pengurus.Email, &pengurus.PasswordHash,
			&pengurus.Peran, &pengurus.Status,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan pengurus: %w", err)
		}
		
		pengurusList = append(pengurusList, &pengurus)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pengurus rows: %w", err)
	}
	
	return pengurusList, nil
}

func (r *PengurusRepository) Update(pengurus *models.Pengurus) error {
	query := `
		UPDATE PENGURUS SET
			nama_lengkap = ?, email = ?, password_hash = ?, peran = ?, status = ?
		WHERE id_pengurus = ?
	`
	
	_, err := r.db.Exec(query,
		pengurus.NamaLengkap, pengurus.Email, pengurus.PasswordHash,
		pengurus.Peran, pengurus.Status, pengurus.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update pengurus: %w", err)
	}
	
	return nil
}

func (r *PengurusRepository) UpdatePassword(id string, passwordHash string) error {
	query := `UPDATE PENGURUS SET password_hash = ? WHERE id_pengurus = ?`
	
	_, err := r.db.Exec(query, passwordHash, id)
	if err != nil {
		return fmt.Errorf("failed to update pengurus password: %w", err)
	}
	
	return nil
}

func (r *PengurusRepository) Delete(id string) error {
	query := `DELETE FROM PENGURUS WHERE id_pengurus = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pengurus: %w", err)
	}
	
	return nil
}