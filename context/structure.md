Struktur *folder* (direktori) untuk aplikasi Golang sangat fleksibel, namun untuk skala MVP dengan arsitektur monolitik (menyatukan *backend*, CMS, dan *frontend* internal/publik), pendekatan terbaik adalah menggunakan **Standard Go Project Layout** yang dipadukan dengan pola **Layered Architecture** (*Handler, Service, Repository*).

Pola ini membuat kode sangat rapi, mudah di- *debug*, dan siap dikembangkan jika fitur bertambah.

Berikut adalah rancangan struktur *folder* untuk sistem lembaga sosial Anda:

```text
puriyatim-app/
├── cmd/
│   └── server/
│       └── main.go              # Titik masuk (Entry point) utama aplikasi
├── internal/
│   ├── config/                  # Membaca file .env (Port, OneSender Token)
│   ├── database/                # Setup koneksi dan inisialisasi SQLite
│   ├── models/                  # Struct Golang yang merepresentasikan ERD
│   ├── repository/              # Urusan query SQL (Select, Insert, Update ke SQLite)
│   ├── services/                # Logika bisnis (Cek kuota Jumat Berkah, rekap kas)
│   ├── handlers/                # Echo HTTP controllers (Terima request, render HTML/JSON)
│   └── middleware/              # Pengecekan sesi login, Hak Akses (RBAC) Admin/Humas
├── pkg/
│   └── onesender/               # Client khusus untuk API Whatsapp OneSender
├── templates/                   # File HTML (Go html/template)
│   ├── layouts/                 # Base layout (header.html, footer.html, sidebar.html)
│   ├── public/                  # Halaman publik (index.html, form-jumat-berkah.html)
│   └── admin/                   # Halaman dashboard (dashboard.html, form-anak.html)
├── static/                      # File statis yang bisa diakses langsung oleh browser
│   ├── css/                     # File CSS output dari Tailwind (misal: app.css)
│   ├── js/                      # Script interaksi UI sederhana
│   └── uploads/                 # Folder simpan bukti transfer & foto (jika lokal)
├── db/
│   ├── migrations/              # Kumpulan file .sql untuk generate tabel (ERD)
│   └── puriyatim.db             # File database SQLite (Otomatis terbuat, abaikan di Git)
├── .env                         # Konfigurasi rahasia
├── tailwind.config.js           # Konfigurasi Tailwind CSS (scan folder templates/)
├── go.mod                       # Daftar dependensi Golang
└── Makefile                     # (Opsional) Kumpulan command cepat (run, build, migrate)

```

### Penjelasan Alur Kerja (*Workflow*) Folder Tersebut:

1. **`internal/models/`**: Di sini Anda mendefinisikan bentuk datanya. Misalnya membuat `type AnakAsuh struct { ID string, Nama string, RT string ... }`.
2. **`internal/repository/`**: Tempat Anda menaruh semua *query* murni ke SQLite. Misalnya fungsi `GetAnakByRT(rt string)`. Ini memisahkan *query* database dari logika web.
3. **`internal/services/`**: Tempat otak aplikasi berada. Misalnya ada fungsi `SubmitJumatBerkah()`. Fungsi ini akan memanggil repositori untuk mengecek apakah kuota masih ada. Jika habis, fungsi ini mengembalikan *error* "Kuota Penuh" ke *handler*.
4. **`internal/handlers/`**: Di sinilah Echo *Framework* bekerja. *Handler* akan menerima *request* dari *browser*, memanggil *Service* untuk meminta data, lalu melemparkan data tersebut ke dalam file HTML yang ada di folder `templates/` untuk di-*render* menjadi halaman web utuh.
5. **`pkg/onesender/`**: Kita pisahkan di luar `internal` agar kode pengiriman WA ini bersifat modular. Isinya murni fungsi HTTP POST Golang yang menembak *endpoint* OneSender. Jika suatu saat Anda ganti vendor API WA, Anda cukup mengubah file di dalam *folder* ini tanpa merusak logika aplikasi lainnya.
6. **`templates/` & `static/**`: Tailwind CSS akan diatur untuk memantau perubahan *class* di dalam semua file HTML di `templates/`. Hasil *compile* Tailwind-nya akan diletakkan di `static/css/app.css` yang kemudian dipanggil oleh file HTML tersebut.

Struktur ini akan membuat proses kompilasi Golang menghasilkan satu *file binary* yang siap dijalankan di VPS Anda (Nginx/Docker), berdampingan dengan *folder* `templates` dan `static`.

