-- Migration: 004_donasi_online.sql
-- Tabel transaksi pembayaran online via Pakasir

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS DONASI_ONLINE (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id        TEXT UNIQUE NOT NULL,         -- ID unik transaksi, format: DON-{timestamp}-{rand}
    jenis           TEXT NOT NULL CHECK (jenis IN ('donasi', 'zakat', 'jumat_berkah')),
    nama_donatur    TEXT NOT NULL DEFAULT 'Hamba Allah',
    nominal         REAL NOT NULL,
    payment_method  TEXT NOT NULL,                -- qris, bri_va, bni_va, dll
    status          TEXT NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'completed', 'expired', 'cancelled')),
    qr_string       TEXT,                         -- QR string untuk QRIS
    va_number       TEXT,                         -- Nomor Virtual Account
    total_payment   REAL,                         -- Total bayar termasuk fee
    fee             REAL,                         -- Fee dari Pakasir
    expired_at      DATETIME,                     -- Waktu kadaluarsa transaksi
    completed_at    DATETIME,                     -- Waktu pembayaran berhasil
    catatan         TEXT,                         -- Catatan tambahan (program, dll)
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_donasi_online_order_id  ON DONASI_ONLINE(order_id);
CREATE INDEX IF NOT EXISTS idx_donasi_online_status    ON DONASI_ONLINE(status);
CREATE INDEX IF NOT EXISTS idx_donasi_online_created   ON DONASI_ONLINE(created_at);

CREATE TRIGGER IF NOT EXISTS update_donasi_online_updated_at
    AFTER UPDATE ON DONASI_ONLINE
    FOR EACH ROW
BEGIN
    UPDATE DONASI_ONLINE SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
