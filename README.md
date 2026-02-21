# Panti App - Sistem Informasi Manajemen & Portal Web Lembaga Sosial

Aplikasi terpadu untuk manajemen lembaga sosial dengan fitur manajemen data internal, sistem antrean otomatis "Jumat Berkah", dan portal web publik.

## 🚀 Fitur Utama

### 📊 Modul Data Internal
- **Manajemen Anak Asuh**: Data lengkap identitas, status, jenjang sekolah, dan pemetaan alamat
- **Manajemen Donatur**: Pencatatan data donatur dengan kategori (Individu/Kelompok/Instansi)
- **Manajemen Keuangan**: Pencatatan pemasukan dan pengeluaran dengan laporan transparan
- **Manajemen Pengurus**: Sistem otentikasi dengan role-based access control

### 📢 Modul Web Publik & CMS
- **Landing Page**: Profil lembaga dengan hero image dan widget transparansi keuangan
- **Manajemen Artikel**: Editor teks dengan pembuatan slug otomatis dan meta description
- **Pengaturan Sistem**: Konfigurasi kontak, alamat, dan link media sosial

### 🕌 Modul "Jumat Berkah"
- **Form Pendaftaran Publik**: Pemilihan RT/RW dengan daftar anak asuh yang valid
- **Manajemen Kuota**: Auto-lock jika kuota mingguan terpenuhi
- **Notifikasi WhatsApp**: Pengiriman pesan persetujuan massal via OneSender API

## 🛠️ Tech Stack

- **Backend**: Go dengan Echo framework
- **Database**: SQLite (file-based, mudah backup)
- **Frontend**: Server-side rendering dengan Go html/template
- **CSS**: Tailwind CSS (mobile-first)
- **WhatsApp API**: OneSender API Gateway

## 📁 Struktur Proyek

```
puriyatim-app/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/                  # Konfigurasi aplikasi
│   ├── database/                # Setup koneksi SQLite
│   ├── models/                  # Struct models (ERD)
│   ├── repository/              # Query SQL
│   ├── services/                # Logika bisnis
│   ├── handlers/                # Echo HTTP controllers
│   └── middleware/              # Middleware otentikasi & RBAC
├── pkg/
│   └── onesender/               # Client WhatsApp API
├── templates/                   # HTML templates
│   ├── layouts/                 # Base layout
│   ├── public/                  # Halaman publik
│   └── admin/                   # Halaman dashboard
├── static/                      # File statis
│   ├── css/                     # CSS output
│   ├── js/                      # JavaScript
│   └── uploads/                 # File uploads
├── db/
│   ├── migrations/              # SQL migrations
│   └── puriyatim.db             # Database SQLite
├── .env                         # Konfigurasi rahasia
├── go.mod                       # Dependencies Go
├── tailwind.config.js           # Konfigurasi Tailwind
└── Makefile                     # Command shortcuts
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js & npm (untuk Tailwind CSS)
- Git

### Installation

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd puriyatim-app
   ```

2. **Install dependencies**
   ```bash
   # Install Go dependencies
   make install
   
   # Install Node.js dependencies
   npm install
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env dengan konfigurasi Anda
   ```

4. **Build CSS**
   ```bash
   make css-prod
   ```

5. **Run database migrations**
   ```bash
   make migrate
   ```

6. **Run application**
   ```bash
   make dev
   ```

7. **Access application**
   - Public website: http://localhost:8080
   - Admin dashboard: http://localhost:8080/admin/login
   - Default credentials: admin@puriyatim.com / admin123

## 📋 Available Commands

```bash
# Build aplikasi
make build

# Run di development mode
make dev

# Run di production mode
make run

# Install dependencies
make install

# Run database migrations
make migrate

# Build CSS development
make css

# Build CSS production
make css-prod

# Clean build artifacts
make clean
```

## 🔧 Konfigurasi

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port server | `8080` |
| `ENV` | Environment (development/production) | `development` |
| `DB_PATH` | Path ke database SQLite | `./db/puriyatim.db` |
| `JWT_SECRET` | Secret key untuk JWT | `default_secret_key` |
| `ONESENDER_API_URL` | URL OneSender API | `https://api.onesender.com/v1/message` |
| `ONESENDER_API_KEY` | API Key OneSender | - |
| `ONESENDER_GROUP_ID` | Group ID WhatsApp | - |

### Database Setup

Database akan otomatis dibuat saat pertama kali aplikasi dijalankan. Untuk menjalankan migration manual:

```bash
sqlite3 db/puriyatim.db < db/migrations/001_create_tables.sql
sqlite3 db/puriyatim.db < db/migrations/002_seed_data.sql
```

## 👥 Role & Permissions

### Superadmin
- Akses penuh ke semua modul
- Manajemen user dan pengaturan sistem
- Laporan dan statistik lengkap

### Keuangan
- Input pemasukan dan pengeluaran
- Cetak laporan keuangan
- Akses data donatur

### Penulis Berita
- Manajemen artikel dan kategori
- Publikasi konten web
- Update kegiatan puri yatim

## 📱 API Endpoints

### Authentication
- `POST /api/auth/login` - Login pengurus
- `POST /api/auth/logout` - Logout
- `GET /api/auth/profile` - Get profile
- `POST /api/auth/change-password` - Ubah password

### Anak Asuh
- `GET /api/anak-asuh` - List anak asuh
- `POST /api/anak-asuh` - Tambah anak asuh
- `GET /api/anak-asuh/:id` - Detail anak asuh
- `PUT /api/anak-asuh/:id` - Update anak asuh
- `DELETE /api/anak-asuh/:id` - Hapus anak asuh

### Donatur
- `GET /api/donatur` - List donatur
- `POST /api/donatur` - Tambah donatur
- `GET /api/donatur/:id` - Detail donatur
- `PUT /api/donatur/:id` - Update donatur
- `DELETE /api/donatur/:id` - Hapus donatur

### Keuangan
- `GET /api/keuangan/pemasukan` - List pemasukan
- `POST /api/keuangan/pemasukan` - Tambah pemasukan
- `GET /api/keuangan/pengeluaran` - List pengeluaran
- `POST /api/keuangan/pengeluaran` - Tambah pengeluaran
- `GET /api/keuangan/laporan` - Laporan keuangan

### Jumat Berkah
- `GET /api/jumat-berkah/kegiatan` - List kegiatan
- `POST /api/jumat-berkah/kegiatan` - Tambah kegiatan
- `POST /api/jumat-berkah/daftar` - Pendaftaran publik
- `POST /api/jumat-berkah/approve` - Approve pendaftar
- `POST /api/jumat-berkah/announce` - Kirim pengumuman WA

## 🎨 Customization

### CSS & Styling
- Edit `tailwind.config.js` untuk konfigurasi Tailwind
- Edit `static/css/app.css` untuk custom CSS
- Run `make css` untuk rebuild CSS development
- Run `make css-prod` untuk build CSS production

### Templates
- Edit files di `templates/` untuk mengubah tampilan
- Base layout: `templates/layouts/base.html`
- Admin pages: `templates/admin/`
- Public pages: `templates/public/`

## 🚀 Deployment

### Build untuk Production
```bash
make build
```

### Environment Setup
1. Set `ENV=production` di environment
2. Build CSS dengan `make css-prod`
3. Setup reverse proxy (nginx/Apache)
4. Setup SSL certificate

### Docker Deployment
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/puriyatim-app .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
CMD ["./puriyatim-app"]
```

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

Untuk bantuan atau pertanyaan:
- Email: support@puriyatim-app.com
- Documentation: [Wiki](https://github.com/username/puriyatim-app/wiki)
- Issues: [GitHub Issues](https://github.com/username/puriyatim-app/issues)

## 🙏 Acknowledgments

- [Echo Framework](https://echo.labstack.com/) - High performance Go web framework
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework
- [SQLite](https://sqlite.org/) - Self-contained database engine
- [OneSender](https://onesender.com/) - WhatsApp API Gateway