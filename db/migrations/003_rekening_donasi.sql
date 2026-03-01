-- Migration: 003_rekening_donasi.sql
-- Menambah tabel REKENING_DONASI untuk mendukung banyak rekening dengan logo bank

CREATE TABLE IF NOT EXISTS REKENING_DONASI (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama_bank TEXT NOT NULL,
    logo_bank TEXT NOT NULL DEFAULT '',
    nomor_rekening TEXT NOT NULL,
    atas_nama TEXT NOT NULL,
    urutan INTEGER NOT NULL DEFAULT 0,
    aktif INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Seed dari data rekening lama jika ada
INSERT OR IGNORE INTO REKENING_DONASI (nama_bank, logo_bank, nomor_rekening, atas_nama, urutan)
SELECT 'BSI', 'bsi', rekening_bsi, COALESCE(nama_pemilik_rekening, 'Puri Yatim'), 1
FROM PENGATURAN_WEB
WHERE rekening_bsi IS NOT NULL AND rekening_bsi != ''
LIMIT 1;

INSERT OR IGNORE INTO REKENING_DONASI (nama_bank, logo_bank, nomor_rekening, atas_nama, urutan)
SELECT 'Mandiri', 'mandiri', rekening_mandiri, COALESCE(nama_pemilik_rekening, 'Puri Yatim'), 2
FROM PENGATURAN_WEB
WHERE rekening_mandiri IS NOT NULL AND rekening_mandiri != ''
LIMIT 1;
