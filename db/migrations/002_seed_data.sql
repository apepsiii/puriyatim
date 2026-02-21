-- Seed data for Panti App
-- Migration: 002_seed_data.sql

-- Insert default superadmin user
-- Password: admin123 (hashed with bcrypt)
INSERT OR IGNORE INTO PENGURUS (
    id_pengurus, 
    nama_lengkap, 
    email, 
    password_hash, 
    peran, 
    status
) VALUES (
    'admin-001',
    'Superadmin',
    'admin@puriyatim.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- admin123
    'Superadmin',
    'Aktif'
);

-- Insert sample pengurus users
INSERT OR IGNORE INTO PENGURUS (
    id_pengurus, 
    nama_lengkap, 
    email, 
    password_hash, 
    peran, 
    status
) VALUES 
(
    'keuangan-001',
    'Ahmad Keuangan',
    'keuangan@puriyatim.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- admin123
    'Keuangan',
    'Aktif'
),
(
    'humas-001',
    'Siti Humas',
    'humas@puriyatim.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- admin123
    'Penulis Berita',
    'Aktif'
);

-- Insert sample anak asuh data
INSERT OR IGNORE INTO ANAK_ASUH (
    id_anak, nik, nama_lengkap, nama_panggilan, tempat_lahir, tanggal_lahir,
    jenis_kelamin, alamat_jalan, rt, rw, desa_kelurahan, kecamatan,
    tanggal_masuk, status_anak, status_aktif, nama_wali, kontak_wali,
    hubungan_wali, jenjang_pendidikan, nama_sekolah, kelas,
    kondisi_kesehatan, catatan_khusus
) VALUES 
(
    'anak-001',
    '3201011234560001',
    'Ahmad Fadillah',
    'Ahmad',
    'Jakarta',
    '2010-05-15',
    'L',
    'Jl. Merdeka No. 123',
    '01',
    '02',
    'Mulyaharja',
    'Bogor Selatan',
    '2015-01-01',
    'Yatim',
    'Aktif',
    'Ibu Siti',
    '081234567890',
    'Ibu',
    'SD',
    'SDN Mulyaharja 01',
    '5',
    'Sehat',
    'Anak yang rajin dan pandai'
),
(
    'anak-002',
    '3201011234560002',
    'Siti Nurhaliza',
    'Siti',
    'Bogor',
    '2011-08-20',
    'P',
    'Jl. Pahlawan No. 456',
    '02',
    '03',
    'Mulyaharja',
    'Bogor Selatan',
    '2016-06-15',
    'Yatim Piatu',
    'Aktif',
    'Nenek Rohmah',
    '082345678901',
    'Nenek',
    'SD',
    'SDN Mulyaharja 02',
    '4',
    'Sehat',
    'Anak yang ceria dan aktif'
),
(
    'anak-003',
    '3201011234560003',
    'Muhammad Rizki',
    'Rizki',
    'Jakarta',
    '2009-12-10',
    'L',
    'Jl. Sudirman No. 789',
    '03',
    '01',
    'Mulyaharja',
    'Bogor Selatan',
    '2014-03-20',
    'Dhuafa',
    'Aktif',
    'Bapak Budi',
    '083456789012',
    'Ayah',
    'SMP',
    'SMPN 1 Bogor',
    '7',
    'Sehat',
    'Anak yang berprestasi di bidang olahraga'
);

-- Insert sample donatur data
INSERT OR IGNORE INTO DONATUR (
    id_donatur, nama_donatur, tipe_donatur, no_telepon, email, alamat, catatan_khusus, tanggal_bergabung
) VALUES 
(
    'donatur-001',
    'H. Ahmad Wijaya',
    'Individu',
    '081234567891',
    'ahmad.wijaya@email.com',
    'Jl. Gatot Subroto No. 100, Jakarta',
    'Donatur rutin bulanan',
    '2023-01-15'
),
(
    'donatur-002',
    'PT. Sejahtera Bersama',
    'Instansi',
    '02112345678',
    'csr@sejahtera.com',
    'Jl. Sudirman No. 200, Jakarta',
    'Program CSR perusahaan',
    '2023-03-20'
),
(
    'donatur-003',
    'Komunitas Peduli',
    'Kelompok',
    '082345678913',
    'info@peduli.com',
    'Jl. Thamrin No. 300, Jakarta',
    'Donatur dari komunitas sosial',
    '2023-06-10'
);

-- Insert sample kegiatan jumat berkah
INSERT OR IGNORE INTO KEGIATAN_JUMAT_BERKAH (
    id_kegiatan, tanggal_kegiatan, kuota_maksimal, total_terdaftar, status_kegiatan
) VALUES 
(
    'jumat-001',
    date('now', '+7 days'),
    50,
    0,
    'Dibuka'
),
(
    'jumat-002',
    date('now', '+14 days'),
    50,
    0,
    'Dibuka'
);

-- Insert sample artikel
INSERT OR IGNORE INTO ARTIKEL (
    id_artikel, id_pengurus, id_kategori, judul, slug, konten_html_markdown, 
    meta_deskripsi, status_publikasi, tanggal_terbit
) VALUES 
(
    'artikel-001',
    'humas-001',
    1,
    'Program Beasiswa Pendidikan 2024',
    'program-beasiswa-pendidikan-2024',
    '<h2>Program Beasiswa Pendidikan 2024</h2><p>Panti Asuhan kembali membuka program beasiswa pendidikan untuk tahun ajaran 2024/2025. Program ini ditujukan untuk anak-anak yatim dan dhuafa yang membutuhkan bantuan biaya pendidikan.</p><p>Persyaratan:</p><ul><li>Anak yatim/dhuafa</li><li>Berprestasi</li><li>Kurang mampu</li></ul>',
    'Program beasiswa pendidikan untuk anak yatim dan dhuafa tahun 2024',
    'Terbit',
    datetime('now')
),
(
    'artikel-002',
    'humas-001',
    2,
    'Kegiatan Buka Bersama Anak Asuh',
    'kegiatan-buka-bersama-anak-asuh',
    '<h2>Buka Bersama Anak Asuh</h2><p>Alhamdulillah, pada tanggal 15 April 2024 telah dilaksanakan kegiatan buka bersama dengan seluruh anak asuh. Kegiatan ini dihadiri oleh para donatur dan wali murid.</p>',
    'Dokumentasi kegiatan buka bersama anak asuh',
    'Terbit',
    datetime('now', '-1 day')
);

-- Insert sample pemasukan donasi
INSERT OR IGNORE INTO PEMASUKAN_DONASI (
    id_pemasukan, id_donatur, tanggal_donasi, nominal, kategori_dana, catatan
) VALUES 
(
    'pemasukan-001',
    'donatur-001',
    date('now', '-30 days'),
    1000000,
    'Infaq',
    'Donasi rutin bulanan'
),
(
    'pemasukan-002',
    'donatur-002',
    date('now', '-15 days'),
    5000000,
    'Sedekah',
    'Program CSR perusahaan'
),
(
    'pemasukan-003',
    'donatur-003',
    date('now', '-7 days'),
    2500000,
    'Zakat',
    'Donasi komunitas'
);

-- Insert sample pengeluaran
INSERT OR IGNORE INTO PENGELUARAN (
    id_pengeluaran, tanggal_pengeluaran, nominal, id_anak, keterangan
) VALUES 
(
    'pengeluaran-001',
    date('now', '-25 days'),
    500000,
    'anak-001',
    'Biaya sekolah dan seragam'
),
(
    'pengeluaran-002',
    date('now', '-20 days'),
    750000,
    'anak-002',
    'Biaya sekolah dan buku'
),
(
    'pengeluaran-003',
    date('now', '-10 days'),
    2000000,
    NULL,
    'Biaya makan bulanan'
);