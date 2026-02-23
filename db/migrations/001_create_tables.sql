-- Create tables for Panti App
-- Migration: 001_create_tables.sql

-- Enable foreign key support
PRAGMA foreign_keys = ON;

-- Create PENGURUS table
CREATE TABLE IF NOT EXISTS PENGURUS (
    id_pengurus TEXT PRIMARY KEY,
    nama_lengkap TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    peran TEXT NOT NULL CHECK (peran IN ('Superadmin', 'Keuangan', 'Penulis Berita')),
    status TEXT NOT NULL DEFAULT 'Aktif' CHECK (status IN ('Aktif', 'Nonaktif')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create PENGATURAN_WEB table
CREATE TABLE IF NOT EXISTS PENGATURAN_WEB (
    id_pengaturan INTEGER PRIMARY KEY CHECK (id_pengaturan = 1),
    nama_lembaga TEXT NOT NULL DEFAULT 'Panti Asuhan',
    deskripsi_tentang_kami TEXT,
    logo_url TEXT,
    hero_image_url TEXT,
    nomor_wa TEXT,
    email_lembaga TEXT,
    alamat_lengkap TEXT,
    link_instagram TEXT,
    link_youtube TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insert default settings
INSERT OR IGNORE INTO PENGATURAN_WEB (id_pengaturan, nama_lembaga, deskripsi_tentang_kami, nomor_wa, email_lembaga, alamat_lengkap)
VALUES (1, 'Puri Yatim', 'Rumah kembali para yatim - Lembaga sosial yang berdedikasi untuk membantu anak-anak yatim dan dhuafa', '+628123456789', 'info@puriyatim.com', 'Jl. Contoh No. 123, Kota, Provinsi');

-- Create KATEGORI_KONTEN table
CREATE TABLE IF NOT EXISTS KATEGORI_KONTEN (
    id_kategori INTEGER PRIMARY KEY AUTOINCREMENT,
    nama_kategori TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insert default categories
INSERT OR IGNORE INTO KATEGORI_KONTEN (nama_kategori, slug) VALUES
('Berita', 'berita'),
('Kegiatan', 'kegiatan'),
('Program', 'program'),
('Pengumuman', 'pengumuman');

-- Create ANAK_ASUH table
CREATE TABLE IF NOT EXISTS ANAK_ASUH (
    id_anak TEXT PRIMARY KEY,
    nik TEXT UNIQUE,
    nama_lengkap TEXT NOT NULL,
    nama_panggilan TEXT NOT NULL,
    tempat_lahir TEXT NOT NULL,
    tanggal_lahir DATE NOT NULL,
    jenis_kelamin TEXT NOT NULL CHECK (jenis_kelamin IN ('L', 'P')),
    alamat_jalan TEXT NOT NULL,
    rt TEXT NOT NULL,
    rw TEXT NOT NULL,
    desa_kelurahan TEXT NOT NULL,
    kecamatan TEXT NOT NULL,
    kota TEXT NOT NULL DEFAULT 'Kota Bogor',
    tanggal_masuk DATE NOT NULL,
    status_anak TEXT NOT NULL CHECK (status_anak IN ('Yatim', 'Piatu', 'Yatim Piatu', 'Dhuafa')),
    status_aktif TEXT NOT NULL DEFAULT 'Aktif' CHECK (status_aktif IN ('Aktif', 'Lulus', 'Keluar')),
    nama_wali TEXT NOT NULL,
    kontak_wali TEXT NOT NULL,
    hubungan_wali TEXT NOT NULL,
    jenjang_pendidikan TEXT NOT NULL,
    nama_sekolah TEXT NOT NULL,
    kelas TEXT NOT NULL,
    kondisi_kesehatan TEXT,
    catatan_khusus TEXT,
    foto_profil_url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create DONATUR table
CREATE TABLE IF NOT EXISTS DONATUR (
    id_donatur TEXT PRIMARY KEY,
    nama_donatur TEXT NOT NULL,
    tipe_donatur TEXT NOT NULL CHECK (tipe_donatur IN ('Individu', 'Kelompok', 'Instansi')),
    no_telepon TEXT NOT NULL,
    email TEXT,
    alamat TEXT,
    catatan_khusus TEXT,
    tanggal_bergabung DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create ARTIKEL table
CREATE TABLE IF NOT EXISTS ARTIKEL (
    id_artikel TEXT PRIMARY KEY,
    id_pengurus TEXT NOT NULL,
    id_kategori INTEGER NOT NULL,
    judul TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    konten_html_markdown TEXT NOT NULL,
    gambar_thumbnail_url TEXT,
    meta_deskripsi TEXT,
    status_publikasi TEXT NOT NULL DEFAULT 'Draft' CHECK (status_publikasi IN ('Draft', 'Terbit', 'Arsip')),
    tanggal_terbit DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_pengurus) REFERENCES PENGURUS(id_pengurus) ON DELETE CASCADE,
    FOREIGN KEY (id_kategori) REFERENCES KATEGORI_KONTEN(id_kategori) ON DELETE RESTRICT
);

-- Create PEMASUKAN_DONASI table
CREATE TABLE IF NOT EXISTS PEMASUKAN_DONASI (
    id_pemasukan TEXT PRIMARY KEY,
    nama_donatur TEXT DEFAULT 'Hamba Allah',
    tanggal_donasi DATE NOT NULL,
    nominal REAL NOT NULL,
    kategori_dana TEXT NOT NULL CHECK (kategori_dana IN ('Infaq', 'Sedekah', 'Wakaf', 'Zakat', 'Lainnya')),
    catatan TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create PENGELUARAN table
CREATE TABLE IF NOT EXISTS PENGELUARAN (
    id_pengeluaran TEXT PRIMARY KEY,
    tanggal_pengeluaran DATE NOT NULL,
    nominal REAL NOT NULL,
    id_anak TEXT,
    keterangan TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_anak) REFERENCES ANAK_ASUH(id_anak) ON DELETE SET NULL
);

-- Create KEGIATAN_JUMAT_BERKAH table
CREATE TABLE IF NOT EXISTS KEGIATAN_JUMAT_BERKAH (
    id_kegiatan TEXT PRIMARY KEY,
    tanggal_kegiatan DATE NOT NULL,
    kuota_maksimal INTEGER NOT NULL DEFAULT 50,
    total_terdaftar INTEGER NOT NULL DEFAULT 0,
    status_kegiatan TEXT NOT NULL DEFAULT 'Dibuka' CHECK (status_kegiatan IN ('Dibuka', 'Ditutup', 'Selesai')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create PENDAFTAR_JUMAT_BERKAH table
CREATE TABLE IF NOT EXISTS PENDAFTAR_JUMAT_BERKAH (
    id_pendaftaran TEXT PRIMARY KEY,
    id_kegiatan TEXT NOT NULL,
    id_anak TEXT NOT NULL,
    waktu_submit DATETIME DEFAULT CURRENT_TIMESTAMP,
    status_approval TEXT NOT NULL DEFAULT 'Menunggu' CHECK (status_approval IN ('Menunggu', 'Disetujui', 'Ditolak')),
    catatan TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_kegiatan) REFERENCES KEGIATAN_JUMAT_BERKAH(id_kegiatan) ON DELETE CASCADE,
    FOREIGN KEY (id_anak) REFERENCES ANAK_ASUH(id_anak) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_anak_asuh_status_aktif ON ANAK_ASUH(status_aktif);
CREATE INDEX IF NOT EXISTS idx_anak_asuh_rt_rw ON ANAK_ASUH(rt, rw);
CREATE INDEX IF NOT EXISTS idx_artikel_status ON ARTIKEL(status_publikasi);
CREATE INDEX IF NOT EXISTS idx_artikel_kategori ON ARTIKEL(id_kategori);
CREATE INDEX IF NOT EXISTS idx_pemasukan_tanggal ON PEMASUKAN_DONASI(tanggal_donasi);
CREATE INDEX IF NOT EXISTS idx_pengeluaran_tanggal ON PENGELUARAN(tanggal_pengeluaran);
CREATE INDEX IF NOT EXISTS idx_kegiatan_tanggal ON KEGIATAN_JUMAT_BERKAH(tanggal_kegiatan);
CREATE INDEX IF NOT EXISTS idx_pendaftaran_kegiatan ON PENDAFTAR_JUMAT_BERKAH(id_kegiatan);
CREATE INDEX IF NOT EXISTS idx_pendaftaran_status ON PENDAFTAR_JUMAT_BERKAH(status_approval);

-- Create triggers for updated_at timestamps
CREATE TRIGGER IF NOT EXISTS update_pengurus_updated_at
    AFTER UPDATE ON PENGURUS
    FOR EACH ROW
BEGIN
    UPDATE PENGURUS SET updated_at = CURRENT_TIMESTAMP WHERE id_pengurus = NEW.id_pengurus;
END;

CREATE TRIGGER IF NOT EXISTS update_pengaturan_web_updated_at
    AFTER UPDATE ON PENGATURAN_WEB
    FOR EACH ROW
BEGIN
    UPDATE PENGATURAN_WEB SET updated_at = CURRENT_TIMESTAMP WHERE id_pengaturan = NEW.id_pengaturan;
END;

CREATE TRIGGER IF NOT EXISTS update_kategori_konten_updated_at
    AFTER UPDATE ON KATEGORI_KONTEN
    FOR EACH ROW
BEGIN
    UPDATE KATEGORI_KONTEN SET updated_at = CURRENT_TIMESTAMP WHERE id_kategori = NEW.id_kategori;
END;

CREATE TRIGGER IF NOT EXISTS update_anak_asuh_updated_at
    AFTER UPDATE ON ANAK_ASUH
    FOR EACH ROW
BEGIN
    UPDATE ANAK_ASUH SET updated_at = CURRENT_TIMESTAMP WHERE id_anak = NEW.id_anak;
END;

CREATE TRIGGER IF NOT EXISTS update_donatur_updated_at
    AFTER UPDATE ON DONATUR
    FOR EACH ROW
BEGIN
    UPDATE DONATUR SET updated_at = CURRENT_TIMESTAMP WHERE id_donatur = NEW.id_donatur;
END;

CREATE TRIGGER IF NOT EXISTS update_artikel_updated_at
    AFTER UPDATE ON ARTIKEL
    FOR EACH ROW
BEGIN
    UPDATE ARTIKEL SET updated_at = CURRENT_TIMESTAMP WHERE id_artikel = NEW.id_artikel;
END;

CREATE TRIGGER IF NOT EXISTS update_pemasukan_donasi_updated_at
    AFTER UPDATE ON PEMASUKAN_DONASI
    FOR EACH ROW
BEGIN
    UPDATE PEMASUKAN_DONASI SET updated_at = CURRENT_TIMESTAMP WHERE id_pemasukan = NEW.id_pemasukan;
END;

CREATE TRIGGER IF NOT EXISTS update_pengeluaran_updated_at
    AFTER UPDATE ON PENGELUARAN
    FOR EACH ROW
BEGIN
    UPDATE PENGELUARAN SET updated_at = CURRENT_TIMESTAMP WHERE id_pengeluaran = NEW.id_pengeluaran;
END;

CREATE TRIGGER IF NOT EXISTS update_kegiatan_jumat_berkah_updated_at
    AFTER UPDATE ON KEGIATAN_JUMAT_BERKAH
    FOR EACH ROW
BEGIN
    UPDATE KEGIATAN_JUMAT_BERKAH SET updated_at = CURRENT_TIMESTAMP WHERE id_kegiatan = NEW.id_kegiatan;
END;

CREATE TRIGGER IF NOT EXISTS update_pendaftar_jumat_berkah_updated_at
    AFTER UPDATE ON PENDAFTAR_JUMAT_BERKAH
    FOR EACH ROW
BEGIN
    UPDATE PENDAFTAR_JUMAT_BERKAH SET updated_at = CURRENT_TIMESTAMP WHERE id_pendaftaran = NEW.id_pendaftaran;
END;