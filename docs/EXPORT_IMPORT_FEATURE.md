# Fitur Export dan Import Data Anak Asuh

## Deskripsi
Fitur ini memungkinkan admin untuk mengekspor data anak asuh ke format Excel atau CSV, dan mengimport data dari file Excel atau CSV untuk menambah atau mengupdate data secara massal.

## Fitur yang Tersedia

### 1. Export Data
Admin dapat mengekspor seluruh data anak asuh dalam format:
- **Excel (.xlsx)** - Format spreadsheet dengan styling dan formatting
- **CSV (.csv)** - Format teks sederhana yang kompatibel dengan berbagai aplikasi

#### Kolom yang Diekspor:
- ID
- NIK
- Nama Lengkap
- Nama Panggilan
- Tempat Lahir
- Tanggal Lahir
- Jenis Kelamin
- Alamat Jalan
- RT
- RW
- Desa/Kelurahan
- Kecamatan
- Kota
- Tanggal Masuk
- Status Anak (Yatim/Piatu/Yatim Piatu/Dhuafa)
- Status Aktif (Aktif/Lulus/Keluar)
- Nama Wali
- Kontak Wali
- Hubungan Wali
- Jenjang Pendidikan
- Nama Sekolah
- Kelas
- Kondisi Kesehatan
- Catatan Khusus

### 2. Download Template Import
Admin dapat mendownload template Excel yang sudah berisi:
- Header kolom yang sesuai
- Contoh data untuk referensi
- Sheet petunjuk penggunaan

### 3. Import Data
Admin dapat mengimport data dari file Excel atau CSV dengan fitur:
- **Validasi format file** - Hanya menerima .xlsx, .xls, atau .csv
- **Validasi ukuran file** - Maksimal 10MB
- **Validasi data** - Memvalidasi format tanggal, field wajib, dll
- **Update atau Create** - Jika ID ada, data akan diupdate; jika tidak, data baru akan dibuat
- **Error reporting** - Menampilkan detail error untuk setiap baris yang gagal
- **Success summary** - Menampilkan jumlah data yang berhasil dan gagal diimport

## Cara Penggunaan

### Export Data

1. Buka halaman **Data Anak Asuh** di admin panel
2. Klik tombol **Export** di pojok kanan atas
3. Pilih format yang diinginkan:
   - **Export ke Excel** - Download dalam format .xlsx
   - **Export ke CSV** - Download dalam format .csv
4. File akan otomatis terdownload dengan nama `data_anak_asuh_YYYYMMDD_HHMMSS.xlsx` atau `.csv`

### Import Data

1. Buka halaman **Data Anak Asuh** di admin panel
2. Klik tombol **Export** → **Download Template** untuk mendapatkan template
3. Isi data pada template sesuai petunjuk:
   - **ID**: Kosongkan untuk data baru, isi dengan ID yang ada untuk update
   - **Tanggal**: Gunakan format YYYY-MM-DD (contoh: 2010-05-15)
   - **Jenis Kelamin**: Isi "Laki-laki" atau "Perempuan"
   - **Status Anak**: Pilih "Yatim", "Piatu", "Yatim Piatu", atau "Dhuafa"
   - **Status Aktif**: Pilih "Aktif", "Lulus", atau "Keluar"
4. Simpan file dalam format .xlsx atau .csv
5. Klik tombol **Export** → **Import Data**
6. Pilih file yang sudah diisi
7. Klik **Import**
8. Tunggu proses selesai dan lihat hasilnya

## Validasi Import

### Format File
- Ekstensi: .xlsx, .xls, atau .csv
- Ukuran maksimal: 10MB

### Field Wajib
- Nama Lengkap
- Nama Panggilan
- Tanggal Lahir (format: YYYY-MM-DD)
- Tanggal Masuk (format: YYYY-MM-DD)

### Format Data
- **Tanggal**: YYYY-MM-DD (contoh: 2010-05-15)
- **Jenis Kelamin**: "Laki-laki" atau "Perempuan" (atau "L"/"P")
- **Status Anak**: "Yatim", "Piatu", "Yatim Piatu", atau "Dhuafa"
- **Status Aktif**: "Aktif", "Lulus", atau "Keluar"

## Error Handling

Jika terjadi error saat import, sistem akan:
1. Tetap memproses baris yang valid
2. Menampilkan jumlah data yang berhasil diimport
3. Menampilkan detail error untuk setiap baris yang gagal
4. Memberikan informasi baris mana yang bermasalah

Contoh pesan error:
- "Baris 5: format tanggal lahir tidak valid (gunakan YYYY-MM-DD)"
- "Baris 8: nama lengkap wajib diisi"
- "Baris 12: Data tidak lengkap"

## Endpoints API

### Export Excel
```
GET /admin/anak-asuh/export/excel
```
Response: File Excel (.xlsx)

### Export CSV
```
GET /admin/anak-asuh/export/csv
```
Response: File CSV (.csv)

### Download Template
```
GET /admin/anak-asuh/export/template
```
Response: File Excel template (.xlsx)

### Import Data
```
POST /admin/anak-asuh/import
Content-Type: multipart/form-data
Body: file (Excel atau CSV)
```
Response:
```json
{
  "success": true,
  "message": "Berhasil mengimport 25 data, 2 data gagal",
  "stats": {
    "success_count": 25,
    "error_count": 2,
    "total": 27,
    "errors": [
      "Baris 5: format tanggal lahir tidak valid",
      "Baris 12: nama lengkap wajib diisi"
    ]
  }
}
```

## File yang Dimodifikasi

1. **internal/services/export_import_service.go** - Service untuk export/import
2. **internal/handlers/anak_asuh_handler.go** - Handler untuk endpoints
3. **cmd/server/main.go** - Routing dan inisialisasi service
4. **templates/admin/anak_asuh_list.html** - UI untuk export/import
5. **go.mod** - Dependency untuk library Excel (excelize)

## Dependencies

- **github.com/xuri/excelize/v2** - Library untuk membaca dan menulis file Excel

## Tips Penggunaan

1. **Backup Data**: Selalu backup data sebelum melakukan import massal
2. **Test dengan Data Kecil**: Coba import dengan beberapa baris data terlebih dahulu
3. **Gunakan Template**: Selalu gunakan template yang disediakan untuk menghindari error format
4. **Periksa Format Tanggal**: Pastikan format tanggal sesuai (YYYY-MM-DD)
5. **Update Data**: Untuk update data yang sudah ada, pastikan kolom ID diisi dengan ID yang benar

## Troubleshooting

### File tidak bisa diupload
- Periksa ukuran file (maksimal 10MB)
- Pastikan format file .xlsx, .xls, atau .csv
- Coba compress file jika terlalu besar

### Data tidak terimport
- Periksa format tanggal (harus YYYY-MM-DD)
- Pastikan field wajib terisi (Nama Lengkap, Nama Panggilan, Tanggal Lahir, Tanggal Masuk)
- Lihat detail error yang ditampilkan

### Error "format tanggal tidak valid"
- Gunakan format YYYY-MM-DD
- Contoh yang benar: 2010-05-15
- Contoh yang salah: 15/05/2010, 15-05-2010, 2010/05/15

## Keamanan

- Fitur ini hanya dapat diakses oleh admin yang sudah login
- Validasi file dilakukan di server untuk mencegah upload file berbahaya
- Ukuran file dibatasi untuk mencegah DoS attack
- Data divalidasi sebelum disimpan ke database

## Future Improvements

1. Import dengan preview data sebelum save
2. Export dengan filter (berdasarkan status, wilayah, dll)
3. Export ke format PDF
4. Scheduled export otomatis
5. Import dengan mapping kolom custom
6. Bulk delete via import
