# Product Requirements Document (PRD) - Update

**Nama Proyek**: Sistem Informasi & Portal Publik Puri Yatim  
**Versi**: 1.1 (Updated with Zakat Module & Mobile-First UI)  
**Tech Stack**: Golang (Echo Framework), SQLite, Go html/template, Tailwind CSS

## 1. Ringkasan Proyek (Executive Summary)

Membangun sebuah sistem informasi berbasis web (Monolithic) yang berfungsi ganda:

1. **Portal Publik (Front-end)**: Menjadi wajah digital Puri Yatim bagi masyarakat dan donatur dengan pendekatan antarmuka Mobile-First (mirip aplikasi HP).
2. **Sistem Manajemen (Back-end/Admin)**: Memudahkan pengurus Puri Yatim dalam mengelola data anak asuh, transparansi kas, pendaftaran program "Jumat Berkah", dan pengelolaan konten web.

## 2. Tujuan Bisnis & Proyek

### Transparansi: 
Menampilkan ringkasan aliran dana masuk dan keluar secara langsung (real-time) kepada publik.

### Efisiensi: 
Mengurangi pencatatan manual (kertas) untuk buku kas dan data anak asuh.

### Otomatisasi: 
Mempermudah persetujuan program "Jumat Berkah" dengan fitur Bulk Approve yang terintegrasi dengan WhatsApp (OneSender API).

### Peningkatan Donasi: 
Mempermudah Muzakki (donatur) menghitung dan menyalurkan zakat profesi/penghasilan melalui fitur Kalkulator Zakat terintegrasi.

## 3. Target Pengguna (User Persona)

### Pengurus / Admin Puri Yatim
Membutuhkan dashboard yang bersih, cepat, dan mudah dioperasikan di laptop maupun HP untuk manajemen data harian.

### Masyarakat Warga (Pendaftar)
Warga sekitar yang ingin mendaftarkan anak asuh/yatim di wilayahnya untuk mendapatkan bantuan (Program Jumat Berkah).

### Donatur / Muzakki
Dermawan yang ingin melihat kegiatan Puri Yatim, mengecek transparansi dana, menghitung zakat, dan melakukan donasi.

## 4. Ruang Lingkup Fitur (Feature Scope)

### A. Sisi Publik (Web Portal - Mobile First)

#### Landing Page (index.html)
- Hero Banner ajakan donasi
- Widget Transparansi Keuangan Bulan Ini (Masuk vs Keluar)
- Quick Menu: Jumat Berkah, Zakat, Anak Asuh, Laporan
- Daftar Berita Terbaru (Horizontal Scroll)

#### Modul Jumat Berkah (jumat_berkah_form.html)
- Form pendaftaran bantuan berdasarkan pemilihan RT/RW
- Validasi kuota otomatis
- Konfirmasi pendaftaran via WhatsApp

#### Modul Zakat
- Kalkulator Zakat Online
- Form pembayaran Zakat
- Riwayat pembayaran

#### Modul Anak Asuh
- Profil anak asuh (dengan foto)
- Kebutuhan dan status pendidikan
- Form adopsi/sponsorship

#### Modul Berita & Kegiatan
- Daftar artikel terbaru
- Detail artikel dengan galeri foto
- Kategori berita (Pendidikan, Kesehatan, dll)

### B. Sisi Admin (Dashboard Management)

#### Dashboard Utama
- Statistik real-time (jumlah anak, kas, pendaftar, dll)
- Grafik aliran dana (masuk/keluar)
- Quick actions (tambah kas, approve pendaftar)
- Notifikasi aktivitas terbaru

#### Manajemen Anak Asuh
- CRUD data lengkap anak asuh
- Upload foto dan dokumen
- Tracking pendidikan dan kesehatan
- Pemetaan alamat dengan Google Maps

#### Manajemen Keuangan & Kas
- Input transaksi masuk/keluar
- Kategori otomatis (SPP, Makan, Kesehatan, dll)
- Export laporan (PDF/Excel)
- Rekap bulanan dan tahunan

#### Manajemen Jumat Berkah
- Daftar pendaftar per wilayah
- Bulk approve/reject
- Integrasi WhatsApp OneSender
- Manajemen kuota per wilayah

#### Manajemen Konten Web (CMS)
- Editor artikel dengan WYSIWYG
- Upload gambar dan galeri
- SEO settings (meta tags, slug)
- Pengaturan menu dan halaman statis

#### Manajemen Pengguna & Hak Akses
- Multi-role (Superadmin, Keuangan, Humas)
- Permission-based access control
- Log aktivitas pengguna

#### Pengaturan Sistem
- Konfigurasi profil lembaga
- Integrasi API (WhatsApp, Payment Gateway)
- Backup & restore database

## 5. Struktur Teknis & Arsitektur

```
puriyatim-app/
├── cmd/server/main.go            # Entry point Echo framework
├── internal/
│   ├── config/                   # Setup .env, DB connection
│   ├── models/                   # Struct database (AnakAsuh, Kas, dll)
│   ├── repository/               # Logic Query SQLite (CRUD)
│   ├── services/                 # Business logic (Hitung kas, validasi kuota)
│   ├── handlers/                 # HTTP Handlers (Echo routing & render template)
│   └── middleware/               # Auth Login, Session
├── pkg/
│   ├── onesender/                # HTTP Client ke API OneSender WA
│   └── zakat/                    # HTTP Client helper ke API Zakat
├── templates/                    # Kumpulan file HTML (Semua UI yang sudah dibuat)
│   ├── layouts/
│   ├── public/
│   └── admin/
├── static/                       # File aset (CSS, JS, Uploads logo/bukti kas)
├── db/puriyatim.db             # File SQLite
├── .env                          # Konfigurasi token WA, Port, dll
└── go.mod
```

## 6. Fase Pengembangan (Roadmap)

- [x] **Fase 1: UI/UX Prototyping** - Selesai. Seluruh halaman HTML statis telah dirancang menggunakan Tailwind CSS.
- [x] **Fase 2: Backend & Database Setup** - Membuat skema database SQLite (ERD) dan struktur awal Golang (Echo + Routing).
- [ ] **Fase 3: Core CRUD Integration** - Menghubungkan form HTML (Kas, Anak Asuh, Artikel) dengan Backend Golang ke Database.
- [ ] **Fase 4: Advanced Features** - Mengaktifkan Kalkulator Zakat dengan API eksternal dan mengintegrasikan OneSender WhatsApp untuk Jumat Berkah.
- [ ] **Fase 5: Testing & Deployment** - Persiapan file binary dan deploy ke VPS ringan (misalnya dengan Nginx).

## 7. Spesifikasi Teknis

### Performance Requirements
- Load time < 2 detik untuk halaman utama
- Support 100+ concurrent users
- Database response time < 500ms

### Security Requirements
- HTTPS encryption
- Input validation dan sanitization
- Session management yang aman
- Role-based access control

### Browser Support
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

### Mobile Responsiveness
- Android 8+ (Chrome Mobile)
- iOS 12+ (Safari Mobile)
- Progressive Web App (PWA) ready

## 8. Integrasi Pihak Ketiga

### OneSender WhatsApp API
- Pengiriman notifikasi otomatis
- Template pesan yang dapat dikustomisasi
- Delivery tracking dan analytics

### Payment Gateway (Opsional)
- Midtrans/Stripe untuk donasi online
- Virtual Account untuk zakat
- Auto-receipt generation

### Google Maps API
- Geolokasi alamat anak asuh
- Radius-based search
- Embed maps di profil

## 9. Kriteria Sukses (Success Metrics)

### Technical Metrics
- Uptime > 99.5%
- Page load < 2 detik
- Mobile usability score > 90

### Business Metrics
- 50% reduction dalam waktu proses pendaftaran Jumat Berkah
- 80% pengguna menggunakan mobile access
- 30% peningkatan donasi melalui portal online

### User Satisfaction
- User feedback score > 4.5/5
- Support ticket resolution < 24 jam
- User retention rate > 85%

---

*Dokumen ini akan terus diperbarui sesuai dengan perkembangan proyek dan feedback dari stakeholder.*