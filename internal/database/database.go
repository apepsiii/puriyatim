package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Konfigurasi connection pool untuk mencegah "database is locked"
	db.SetMaxOpenConns(1)  // SQLite tidak mendukung multiple writers
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err = applyPragmas(db); err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return &DB{db}, nil
}

// applyPragmas mengatur PRAGMA SQLite untuk keamanan, performa, dan konsistensi data.
func applyPragmas(db *sql.DB) error {
	pragmas := []string{
		// Aktifkan foreign key constraints
		"PRAGMA foreign_keys = ON",
		// WAL mode: lebih baik untuk concurrent reads + satu writer
		"PRAGMA journal_mode = WAL",
		// Tunggu hingga 5 detik jika database terkunci (mencegah error langsung)
		"PRAGMA busy_timeout = 5000",
		// Sinkronisasi normal: aman dan lebih cepat dari FULL
		"PRAGMA synchronous = NORMAL",
		// Cache 64MB di memori (nilai negatif = kibibytes)
		"PRAGMA cache_size = -65536",
		// Simpan tabel temp di memori
		"PRAGMA temp_store = MEMORY",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return fmt.Errorf("failed to apply pragma [%s]: %w", p, err)
		}
	}
	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// InitializeTables dipanggil sekali saat startup untuk memastikan PRAGMA aktif.
// Deprecated: PRAGMA sekarang diapply di NewDB via applyPragmas.
func (db *DB) InitializeTables() error {
	log.Println("Database initialized successfully")
	return nil
}
