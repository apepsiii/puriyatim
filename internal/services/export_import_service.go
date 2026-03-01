package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"puriyatim-app/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExportImportService struct {
	anakAsuhService *AnakAsuhService
}

func NewExportImportService(anakAsuhService *AnakAsuhService) *ExportImportService {
	return &ExportImportService{
		anakAsuhService: anakAsuhService,
	}
}

// ExportToExcel exports anak asuh data to Excel format
func (s *ExportImportService) ExportToExcel() (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Data Anak Asuh"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Set headers
	headers := []string{
		"ID", "NIK", "Nama Lengkap", "Nama Panggilan", "Tempat Lahir", "Tanggal Lahir",
		"Jenis Kelamin", "Alamat Jalan", "RT", "RW", "Desa/Kelurahan", "Kecamatan", "Kota",
		"Tanggal Masuk", "Status Anak", "Status Aktif", "Nama Wali", "Kontak Wali",
		"Hubungan Wali", "Jenjang Pendidikan", "Nama Sekolah", "Kelas",
		"Kondisi Kesehatan", "Catatan Khusus",
	}

	// Style for header
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Get all data
	anakAsuhList, err := s.anakAsuhService.GetAll()
	if err != nil {
		return nil, err
	}

	// Write data
	for rowIdx, anak := range anakAsuhList {
		row := rowIdx + 2 // Start from row 2 (after header)

		nik := ""
		if anak.NIK != nil {
			nik = *anak.NIK
		}

		jenisKelamin := "Laki-laki"
		if anak.JenisKelamin == "P" {
			jenisKelamin = "Perempuan"
		}

		data := []interface{}{
			anak.ID,
			nik,
			anak.NamaLengkap,
			anak.NamaPanggilan,
			anak.TempatLahir,
			anak.TanggalLahir.Format("2006-01-02"),
			jenisKelamin,
			anak.AlamatJalan,
			anak.RT,
			anak.RW,
			anak.DesaKelurahan,
			anak.Kecamatan,
			anak.Kota,
			anak.TanggalMasuk.Format("2006-01-02"),
			string(anak.StatusAnak),
			string(anak.StatusAktif),
			anak.NamaWali,
			anak.KontakWali,
			anak.HubunganWali,
			anak.JenjangPendidikan,
			anak.NamaSekolah,
			anak.Kelas,
			anak.KondisiKesehatan,
			anak.CatatanKhusus,
		}

		for colIdx, value := range data {
			cell := fmt.Sprintf("%s%d", string(rune('A'+colIdx)), row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Auto-fit columns
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

// ExportToCSV exports anak asuh data to CSV format
func (s *ExportImportService) ExportToCSV(w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write headers
	headers := []string{
		"ID", "NIK", "Nama Lengkap", "Nama Panggilan", "Tempat Lahir", "Tanggal Lahir",
		"Jenis Kelamin", "Alamat Jalan", "RT", "RW", "Desa/Kelurahan", "Kecamatan", "Kota",
		"Tanggal Masuk", "Status Anak", "Status Aktif", "Nama Wali", "Kontak Wali",
		"Hubungan Wali", "Jenjang Pendidikan", "Nama Sekolah", "Kelas",
		"Kondisi Kesehatan", "Catatan Khusus",
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Get all data
	anakAsuhList, err := s.anakAsuhService.GetAll()
	if err != nil {
		return err
	}

	// Write data
	for _, anak := range anakAsuhList {
		nik := ""
		if anak.NIK != nil {
			nik = *anak.NIK
		}

		jenisKelamin := "Laki-laki"
		if anak.JenisKelamin == "P" {
			jenisKelamin = "Perempuan"
		}

		record := []string{
			anak.ID,
			nik,
			anak.NamaLengkap,
			anak.NamaPanggilan,
			anak.TempatLahir,
			anak.TanggalLahir.Format("2006-01-02"),
			jenisKelamin,
			anak.AlamatJalan,
			anak.RT,
			anak.RW,
			anak.DesaKelurahan,
			anak.Kecamatan,
			anak.Kota,
			anak.TanggalMasuk.Format("2006-01-02"),
			string(anak.StatusAnak),
			string(anak.StatusAktif),
			anak.NamaWali,
			anak.KontakWali,
			anak.HubunganWali,
			anak.JenjangPendidikan,
			anak.NamaSekolah,
			anak.Kelas,
			anak.KondisiKesehatan,
			anak.CatatanKhusus,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// ImportFromExcel imports anak asuh data from Excel file
func (s *ExportImportService) ImportFromExcel(file *multipart.FileHeader) (int, []string, error) {
	src, err := file.Open()
	if err != nil {
		return 0, nil, err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return 0, nil, fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return 0, nil, err
	}

	if len(rows) < 2 {
		return 0, nil, fmt.Errorf("file tidak memiliki data")
	}

	successCount := 0
	var errors []string

	// Skip header row
	for i, row := range rows[1:] {
		rowNum := i + 2 // Actual row number in Excel

		if len(row) < 15 { // Minimum required columns
			errors = append(errors, fmt.Sprintf("Baris %d: Data tidak lengkap", rowNum))
			continue
		}

		// Parse data
		anak, err := s.parseRowToAnakAsuh(row, rowNum)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Baris %d: %s", rowNum, err.Error()))
			continue
		}

		// Check if ID exists (for update) or create new
		if anak.ID != "" {
			existing, err := s.anakAsuhService.GetByID(anak.ID)
			if err == nil && existing != nil {
				// Update existing
				anak.CreatedAt = existing.CreatedAt
				if err := s.anakAsuhService.Update(anak); err != nil {
					errors = append(errors, fmt.Sprintf("Baris %d: Gagal update - %s", rowNum, err.Error()))
					continue
				}
			} else {
				// Create new
				if err := s.anakAsuhService.Create(anak); err != nil {
					errors = append(errors, fmt.Sprintf("Baris %d: Gagal create - %s", rowNum, err.Error()))
					continue
				}
			}
		} else {
			// Create new without ID
			if err := s.anakAsuhService.Create(anak); err != nil {
				errors = append(errors, fmt.Sprintf("Baris %d: Gagal create - %s", rowNum, err.Error()))
				continue
			}
		}

		successCount++
	}

	return successCount, errors, nil
}

// ImportFromCSV imports anak asuh data from CSV file
func (s *ExportImportService) ImportFromCSV(file *multipart.FileHeader) (int, []string, error) {
	src, err := file.Open()
	if err != nil {
		return 0, nil, err
	}
	defer src.Close()

	reader := csv.NewReader(src)
	rows, err := reader.ReadAll()
	if err != nil {
		return 0, nil, err
	}

	if len(rows) < 2 {
		return 0, nil, fmt.Errorf("file tidak memiliki data")
	}

	successCount := 0
	var errors []string

	// Skip header row
	for i, row := range rows[1:] {
		rowNum := i + 2

		if len(row) < 15 {
			errors = append(errors, fmt.Sprintf("Baris %d: Data tidak lengkap", rowNum))
			continue
		}

		anak, err := s.parseRowToAnakAsuh(row, rowNum)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Baris %d: %s", rowNum, err.Error()))
			continue
		}

		// Check if ID exists (for update) or create new
		if anak.ID != "" {
			existing, err := s.anakAsuhService.GetByID(anak.ID)
			if err == nil && existing != nil {
				// Update existing
				anak.CreatedAt = existing.CreatedAt
				if err := s.anakAsuhService.Update(anak); err != nil {
					errors = append(errors, fmt.Sprintf("Baris %d: Gagal update - %s", rowNum, err.Error()))
					continue
				}
			} else {
				// Create new
				if err := s.anakAsuhService.Create(anak); err != nil {
					errors = append(errors, fmt.Sprintf("Baris %d: Gagal create - %s", rowNum, err.Error()))
					continue
				}
			}
		} else {
			// Create new without ID
			if err := s.anakAsuhService.Create(anak); err != nil {
				errors = append(errors, fmt.Sprintf("Baris %d: Gagal create - %s", rowNum, err.Error()))
				continue
			}
		}

		successCount++
	}

	return successCount, errors, nil
}

// parseRowToAnakAsuh converts a row of data to AnakAsuh model
func (s *ExportImportService) parseRowToAnakAsuh(row []string, rowNum int) (*models.AnakAsuh, error) {
	// Helper function to get value or empty string
	getValue := func(index int) string {
		if index < len(row) {
			return strings.TrimSpace(row[index])
		}
		return ""
	}

	// Parse dates
	tanggalLahir, err := time.Parse("2006-01-02", getValue(5))
	if err != nil {
		return nil, fmt.Errorf("format tanggal lahir tidak valid (gunakan YYYY-MM-DD)")
	}

	tanggalMasuk, err := time.Parse("2006-01-02", getValue(13))
	if err != nil {
		return nil, fmt.Errorf("format tanggal masuk tidak valid (gunakan YYYY-MM-DD)")
	}

	// Parse jenis kelamin
	jenisKelamin := models.JenisKelaminLakiLaki
	jkStr := strings.ToLower(getValue(6))
	if jkStr == "perempuan" || jkStr == "p" {
		jenisKelamin = models.JenisKelaminPerempuan
	}

	// Parse status anak
	statusAnak := models.StatusAnak(getValue(14))
	if statusAnak == "" {
		statusAnak = models.StatusAnakYatim
	}

	// Parse status aktif
	statusAktif := models.StatusAktif(getValue(15))
	if statusAktif == "" {
		statusAktif = models.StatusAktifAktif
	}

	// Validate required fields
	namaLengkap := getValue(2)
	if namaLengkap == "" {
		return nil, fmt.Errorf("nama lengkap wajib diisi")
	}

	namaPanggilan := getValue(3)
	if namaPanggilan == "" {
		return nil, fmt.Errorf("nama panggilan wajib diisi")
	}

	anak := &models.AnakAsuh{
		ID:                getValue(0),
		NamaLengkap:       namaLengkap,
		NamaPanggilan:     namaPanggilan,
		TempatLahir:       getValue(4),
		TanggalLahir:      tanggalLahir,
		JenisKelamin:      jenisKelamin,
		AlamatJalan:       getValue(7),
		RT:                getValue(8),
		RW:                getValue(9),
		DesaKelurahan:     getValue(10),
		Kecamatan:         getValue(11),
		Kota:              getValue(12),
		TanggalMasuk:      tanggalMasuk,
		StatusAnak:        statusAnak,
		StatusAktif:       statusAktif,
		NamaWali:          getValue(16),
		KontakWali:        getValue(17),
		HubunganWali:      getValue(18),
		JenjangPendidikan: getValue(19),
		NamaSekolah:       getValue(20),
		Kelas:             getValue(21),
		KondisiKesehatan:  getValue(22),
		CatatanKhusus:     getValue(23),
	}

	// Handle NIK (optional)
	nik := getValue(1)
	if nik != "" {
		anak.NIK = &nik
	}

	return anak, nil
}

// GetImportTemplate generates a template Excel file for import
func (s *ExportImportService) GetImportTemplate() (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Template Import"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Set headers
	headers := []string{
		"ID", "NIK", "Nama Lengkap", "Nama Panggilan", "Tempat Lahir", "Tanggal Lahir",
		"Jenis Kelamin", "Alamat Jalan", "RT", "RW", "Desa/Kelurahan", "Kecamatan", "Kota",
		"Tanggal Masuk", "Status Anak", "Status Aktif", "Nama Wali", "Kontak Wali",
		"Hubungan Wali", "Jenjang Pendidikan", "Nama Sekolah", "Kelas",
		"Kondisi Kesehatan", "Catatan Khusus",
	}

	// Style for header
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Add example data
	exampleData := []interface{}{
		"", // ID (kosongkan untuk data baru)
		"3201234567890123",
		"Ahmad Fauzi",
		"Fauzi",
		"Bandung",
		"2010-05-15",
		"Laki-laki",
		"Jl. Merdeka No. 123",
		"001",
		"005",
		"Sukajadi",
		"Sukajadi",
		"Bandung",
		"2020-01-10",
		"Yatim",
		"Aktif",
		"Siti Aminah",
		"081234567890",
		"Ibu",
		"SD",
		"SDN 1 Sukajadi",
		"6",
		"Sehat",
		"Anak rajin dan berprestasi",
	}

	for colIdx, value := range exampleData {
		cell := fmt.Sprintf("%s2", string(rune('A'+colIdx)))
		f.SetCellValue(sheetName, cell, value)
	}

	// Add instructions sheet
	instructionSheet := "Petunjuk"
	f.NewSheet(instructionSheet)
	
	instructions := []string{
		"PETUNJUK IMPORT DATA ANAK ASUH",
		"",
		"1. Kolom ID: Kosongkan untuk data baru, isi dengan ID yang ada untuk update data",
		"2. Tanggal Lahir & Tanggal Masuk: Format YYYY-MM-DD (contoh: 2010-05-15)",
		"3. Jenis Kelamin: Isi dengan 'Laki-laki' atau 'Perempuan'",
		"4. Status Anak: Pilih salah satu - Yatim, Piatu, Yatim Piatu, atau Dhuafa",
		"5. Status Aktif: Pilih salah satu - Aktif, Lulus, atau Keluar",
		"6. Kolom yang wajib diisi: Nama Lengkap, Nama Panggilan, Tanggal Lahir, Tanggal Masuk",
		"",
		"TIPS:",
		"- Gunakan sheet 'Template Import' untuk mengisi data",
		"- Hapus baris contoh sebelum mengisi data Anda",
		"- Pastikan format tanggal sesuai (YYYY-MM-DD)",
		"- Simpan file dalam format .xlsx atau .csv",
	}

	for i, instruction := range instructions {
		cell := fmt.Sprintf("A%d", i+1)
		f.SetCellValue(instructionSheet, cell, instruction)
	}

	// Auto-fit columns
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	f.SetColWidth(instructionSheet, "A", "A", 80)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

// ValidateImportFile validates the structure of import file
func (s *ExportImportService) ValidateImportFile(file *multipart.FileHeader) error {
	// Check file extension
	filename := file.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") &&
		!strings.HasSuffix(strings.ToLower(filename), ".xls") &&
		!strings.HasSuffix(strings.ToLower(filename), ".csv") {
		return fmt.Errorf("format file tidak didukung. Gunakan .xlsx, .xls, atau .csv")
	}

	// Check file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return fmt.Errorf("ukuran file terlalu besar. Maksimal 10MB")
	}

	return nil
}

// GetImportStats returns statistics about the import
func (s *ExportImportService) GetImportStats(successCount int, errors []string) map[string]interface{} {
	return map[string]interface{}{
		"success_count": successCount,
		"error_count":   len(errors),
		"total":         successCount + len(errors),
		"errors":        errors,
	}
}

// Helper function to parse int safely
func parseInt(s string) int {
	val, _ := strconv.Atoi(strings.TrimSpace(s))
	return val
}
