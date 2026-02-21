# PRODUCT REQUIREMENTS DOCUMENT (PRD)

**Sistem Informasi Manajemen & Portal Web Lembaga Sosial (MVP V1.0)**

## 1. Ringkasan Eksekutif

Proyek ini bertujuan membangun *platform* terpadu dan berkinerja tinggi untuk lembaga sosial berskala lokal. Sistem mengintegrasikan manajemen data internal (anak asuh, donatur, keuangan), sistem antrean otomatis "Jumat Berkah", dan portal web publik. Menggunakan arsitektur monolitik yang ringan, MVP ini dirancang untuk kecepatan eksekusi, kemudahan operasional pengurus, dan kemudahan *deployment* ke *server*.

## 2. Objektif & Metrik Keberhasilan

* **Objektif Utama:** Mendigitalisasi data operasional, mengeliminasi *race condition* pendaftaran bantuan di grup WhatsApp, dan menyediakan portal publik yang SEO-friendly.
* **Metrik Keberhasilan Operasional:**
* Pengurus dapat mencetak laporan keuangan bulanan dalam waktu di bawah 1 menit.
* Tidak ada lagi data pendaftar ganda pada program "Jumat Berkah".


* **Metrik Keberhasilan Teknis:**
* Aplikasi dapat berjalan optimal pada VPS spesifikasi rendah (1GB RAM) berkat efisiensi Golang.
* *Backup database* dapat dilakukan secara instan karena menggunakan *file* tunggal.



## 3. Peran Pengguna (Role-Based Access)

* **Superadmin:** Akses penuh ke manajemen *user*, pengaturan profil web, dan seluruh data.
* **Admin Keuangan:** Akses khusus ke *input* pemasukan, pengeluaran, dan cetak laporan.
* **Humas / Penulis:** Akses khusus ke modul CMS (Kategori dan Artikel) untuk *update* kegiatan puri yatim.
* **Warga / Wali:** Pengguna tanpa akun yang mengisi form pendaftaran Jumat Berkah via *smartphone*.
* **Publik:** Pengunjung anonim yang membaca artikel berita dan melihat ringkasan transparansi kas.

## 4. Ruang Lingkup Fitur (MVP)

**Modul Web Publik & CMS**

* **Landing Page:** Menampilkan profil, *hero image*, program, dan *widget* agregat keuangan transparan.
* **Manajemen Artikel:** Editor teks sederhana, pembuatan *slug* otomatis, dan pengaturan *meta description*.
* **Pengaturan Sistem:** Form *single-row* untuk memperbarui kontak WA, alamat, dan *link* sosial media.

**Modul "Jumat Berkah" (Otomatisasi WA)**

* **Form Publik:** Pemilihan RT/RW (misal: area Kelurahan Mulyaharja) yang memicu munculnya daftar anak asuh yang valid.
* **Manajemen Kuota:** Fitur *auto-lock* jika kuota mingguan terpenuhi.
* **Notifikasi Pintar:** Pengiriman pesan persetujuan massal ke grup warga menggunakan layanan pihak ketiga.

**Modul Data Inti (Core)**

* **Manajemen Anak Asuh:** Data lengkap identitas, status (Yatim/Dhuafa), jenjang sekolah, dan pemetaan alamat.
* **Manajemen Keuangan:** Pencatatan donasi masuk (terikat ke data Donatur) dan kas keluar (bisa terikat ke anak asuh spesifik).

## 5. Arsitektur Teknis & Tech Stack

Sistem ini menggunakan pendekatan Monolitik berkinerja tinggi dengan spesifikasi berikut:

* **Bahasa Pemrograman & Framework:** **Golang** dengan framework **Echo**. Dipilih karena *routing* yang sangat cepat, konsumsi memori yang minim, dan kompilasi menjadi satu *binary execution*.
* **Database:** **SQLite**. Sangat ideal untuk skala puri yatim/lembaga lokal. Tidak memerlukan instalasi *service database* terpisah, cukup satu *file* `.db` yang mudah di-*backup* atau dipindahkan.
* **Frontend UI/UX:** Menggunakan **Go `html/template**` yang di-*render* di sisi *server* (SSR), dikombinasikan dengan **Tailwind CSS** (via CDN atau proses *build* sederhana) untuk desain *mobile-first* yang cepat tanpa perlu *framework JavaScript* yang berat.
* **WhatsApp API Gateway:** **OneSender API**. Eksekusi pengiriman pesan dilakukan dari sisi *backend* Golang menggunakan standar HTTP Request (`net/http` bawaan Go) menuju *endpoint* OneSender setelah admin melakukan *Bulk Approve* pendaftar Jumat Berkah.

## 6. Skenario Integrasi WhatsApp (OneSender)

1. Admin menekan tombol "Setujui & Umumkan" di halaman *dashboard* Jumat Berkah.
2. *Handler* Golang memperbarui status pendaftar menjadi "Disetujui" di SQLite.
3. Golang merangkai *string template* berisi daftar nama penerima bantuan.
4. Fungsi `http.Post` di Golang mengirim *payload* JSON berisi pesan tersebut ke *endpoint* OneSender API.
5. Pesan langsung mendarat di Grup WhatsApp warga.
---

## 7. Tutorial Instalasi & Konfigurasi

### 7.1 Persyaratan Sistem

* **Go 1.19+** - Untuk kompilasi aplikasi
* **SQLite3** - Sudah terintegrasi dalam aplikasi
* **Git** - Untuk cloning repository (opsional)

### 7.2 Instalasi Aplikasi

1. **Clone Repository**
   ```bash
   git clone <repository-url>
   cd puriyatim-app
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Konfigurasi Environment**
   - Salin file `.env.example` ke `.env` (jika ada)
   - Edit file `.env` sesuai kebutuhan:
     ```env
     # Server Configuration
     PORT=8083
     ENV=development
     
     # Database Configuration
     DB_PATH=./db/puriyatim.db
     
     # JWT Configuration
     JWT_SECRET=fufufafaabadi
     
     # Application Configuration
     APP_NAME=Puri Yatim
     APP_VERSION=1.0.0
     ```

4. **Setup Database**
   ```bash
   sqlite3 db/puriyatim.db < db/migrations/001_create_tables.sql
   sqlite3 db/puriyatim.db < db/migrations/002_seed_data.sql
   ```

5. **Build & Run Aplikasi**
   ```bash
   # Build aplikasi
   make build
   
   # Atau jalankan langsung
   make run
   ```

### 7.3 Manajemen Port

#### Cara Menjalankan di Port Tertentu

1. **Edit file .env**
   ```env
   PORT=8083
   ```

2. **Restart aplikasi** untuk menerapkan perubahan

#### Cara Menghentikan Proses di Port Tertentu (Windows)

1. **Cari proses yang menggunakan port**
   ```bash
   netstat -ano | findstr :8083
   ```

2. **Hentikan proses dengan PID yang ditemukan**
   ```bash
   taskkill /PID [PID_NUMBER] /F
   ```

3. **Atau gunakan perintah satu baris**
   ```bash
   for /f "tokens=5" %a in ('netstat -ano ^| findstr :8083') do taskkill /PID %a /F
   ```

#### Cara Menghentikan Proses di Port Tertentu (Linux/Mac)

1. **Cari proses yang menggunakan port**
   ```bash
   lsof -i :8083
   ```

2. **Hentikan proses dengan PID yang ditemukan**
   ```bash
   kill -9 [PID_NUMBER]
   ```

3. **Atau gunakan perintah satu baris**
   ```bash
   kill -9 $(lsof -t -i:8083)
   ```

### 7.4 Akses Aplikasi

Setelah berhasil dijalankan, aplikasi dapat diakses melalui:

* **URL Utama:** http://localhost:8083
* **Admin Dashboard:** http://localhost:8083/admin/login
* **Default Credentials:** admin@puriyatim.com / admin123

### 7.5 Troubleshooting

#### Port Sudah Digunakan
Jika muncul error "bind: address already in use", ikuti langkah di bagian "Manajemen Port" untuk menghentikan proses yang menggunakan port tersebut.

#### Database Not Found
Jika muncul error database tidak ditemukan, pastikan:
1. File `db/puriyatim.db` sudah ada
2. Jalankan migrasi database jika belum dilakukan
3. Periksa path di file `.env`

#### Permission Denied
Jika muncul error permission, pastikan:
1. User memiliki hak akses ke folder project
2. Folder `db/` memiliki permission untuk menulis file


