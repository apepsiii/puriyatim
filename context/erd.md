erDiagram
    ANAK_ASUH {
        string id_anak PK
        string NIK "Opsional"
        string nama_lengkap
        string nama_panggilan
        string tempat_lahir
        date tanggal_lahir
        enum jenis_kelamin "L / P"
        string alamat_jalan "Nama jalan, gang, atau dusun"
        string rt
        string rw
        string desa_kelurahan
        string kecamatan
        date tanggal_masuk
        enum status_anak "Yatim / Piatu / Yatim Piatu / Dhuafa"
        enum status_aktif "Aktif / Lulus / Keluar"
        string nama_wali 
        string kontak_wali
        string hubungan_wali
        string jenjang_pendidikan "SD / SMP / SMK / dll"
        string nama_sekolah
        string kelas
        text kondisi_kesehatan 
        text catatan_khusus 
        string foto_profil_url
    }

    DONATUR {
        string id_donatur PK
        string nama_donatur
        enum tipe_donatur "Individu / Kelompok / Instansi"
        string no_telepon
        string email "Opsional"
        text alamat "Opsional"
        text catatan_khusus 
        date tanggal_bergabung
    }

    %% ================= MODUL OTENTIKASI & PENGURUS =================
    PENGURUS {
        string id_pengurus PK
        string nama_lengkap
        string email
        string password_hash
        enum peran "Superadmin / Keuangan / Penulis Berita"
        enum status "Aktif / Nonaktif"
    }

    %% ================= MODUL WEB & CMS (BARU) =================
    PENGATURAN_WEB {
        int id_pengaturan PK "Hanya 1 baris (Single Row)"
        string nama_lembaga
        text deskripsi_tentang_kami
        string logo_url
        string hero_image_url "Gambar banner depan"
        string nomor_wa
        string email_lembaga
        text alamat_lengkap
        string link_instagram
        string link_youtube
    }

    KATEGORI_KONTEN {
        int id_kategori PK
        string nama_kategori "Misal: Berita, Kegiatan, Program"
        string slug "Untuk URL SEO-friendly"
    }

    ARTIKEL {
        string id_artikel PK
        string id_pengurus FK "Penulis"
        int id_kategori FK
        string judul
        string slug "Format URL ramah mesin pencari"
        text konten_html_markdown
        string gambar_thumbnail_url
        string meta_deskripsi "Cuplikan untuk SEO/Share WA"
        enum status_publikasi "Draft / Terbit / Arsip"
        datetime tanggal_terbit
    }

    %% ================= MODUL CORE INTERNAL (SEBELUMNYA) =================

    DONATUR {
        string id_donatur PK
        string nama_donatur
        string no_telepon
        text catatan_khusus
	enum jenis_donatur "Individu/ Lembaga/ Lainnya"
    }

    PEMASUKAN_DONASI {
        string id_pemasukan PK
        string id_donatur FK "Bisa Null"
        date tanggal_donasi
        decimal nominal
        enum kategori_dana
    }

    PENGELUARAN {
        string id_pengeluaran PK
        date tanggal_pengeluaran
        decimal nominal
        string id_anak FK "Bisa Null"
    }

    %% ================= MODUL JUMAT BERKAH =================
    KEGIATAN_JUMAT_BERKAH {
        string id_kegiatan PK
        date tanggal_kegiatan
        int kuota_maksimal
        int total_terdaftar
        enum status_kegiatan
    }

    PENDAFTAR_JUMAT_BERKAH {
        string id_pendaftaran PK
        string id_kegiatan FK
        string id_anak FK
        datetime waktu_submit
        enum status_approval
    }

    %% ================= RELASI =================
    PENGURUS ||--o{ ARTIKEL : "menulis"
    KATEGORI_KONTEN ||--o{ ARTIKEL : "mengelompokkan"
    DONATUR ||--o{ PEMASUKAN_DONASI : "melakukan"
    ANAK_ASUH ||--o{ PENGELUARAN : "menerima"
    ANAK_ASUH ||--o{ PENDAFTAR_JUMAT_BERKAH : "didaftarkan"
    KEGIATAN_JUMAT_BERKAH ||--o{ PENDAFTAR_JUMAT_BERKAH : "menampung"
