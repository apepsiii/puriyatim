package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"puriyatim-app/internal/config"
	"puriyatim-app/internal/database"
	"puriyatim-app/internal/handlers"
	authmw "puriyatim-app/internal/middleware"
	"puriyatim-app/internal/repository"
	"puriyatim-app/internal/services"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadConfig()

	e := echo.New()

	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.Secure())

	// Custom HTTP error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// Jangan handle jika response sudah dikirim
		if c.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		data := map[string]interface{}{
			"Year": time.Now().Year(),
		}

		// Render template sesuai kode error
		switch code {
		case http.StatusNotFound:
			if renderErr := c.Render(http.StatusNotFound, "public/404.html", data); renderErr != nil {
				c.String(http.StatusNotFound, "404 - Halaman tidak ditemukan")
			}
		case http.StatusForbidden, http.StatusUnauthorized:
			if renderErr := c.Render(http.StatusForbidden, "public/403.html", data); renderErr != nil {
				c.String(http.StatusForbidden, "403 - Akses ditolak")
			}
		default:
			// Untuk error lain, gunakan handler bawaan Echo
			e.DefaultHTTPErrorHandler(err, c)
		}
	}

	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeImageURL": func(s string) template.URL {
			url := strings.TrimSpace(s)
			if strings.HasPrefix(url, "data:image/") ||
				strings.HasPrefix(url, "/static/") ||
				strings.HasPrefix(url, "http://") ||
				strings.HasPrefix(url, "https://") {
				return template.URL(url)
			}
			return template.URL("")
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"slice": func(s string, start int) string {
			if start >= len(s) {
				return ""
			}
			return s[start:]
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"lower": func(s string) string {
			result := make([]byte, len(s))
			for i, c := range s {
				if c >= 'A' && c <= 'Z' {
					result[i] = byte(c + 32)
				} else {
					result[i] = byte(c)
				}
			}
			return string(result)
		},
		"add": func(a, b int) int {
			return a + b
		},
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/**/*.html"))
	e.Renderer = &TemplateRenderer{
		templates: tmpl,
	}

	e.Static("/static", "static")

	// PWA routes — harus dapat diakses dari root scope
	e.File("/manifest.json", "static/manifest.json")
	e.File("/sw.js", "static/sw.js")
	e.File("/favicon.ico", "static/images/icons/favicon-32x32.png")

	db, err := database.NewDB(cfg.DBPath)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Running without database connection - using mock data")
	}

	var pengurusRepo *repository.PengurusRepository
	var anakAsuhRepo *repository.AnakAsuhRepository
	var keuanganRepo *repository.KeuanganRepository
	var jumatBerkahRepo *repository.JumatBerkahRepository
	var artikelRepo *repository.ArtikelRepository
	var pengaturanRepo *repository.PengaturanRepository
	var galeriRepo *repository.GaleriRepository
	var rekeningRepo *repository.RekeningDonasiRepository

	if db != nil {
		pengurusRepo = repository.NewPengurusRepository(db.DB)
		anakAsuhRepo = repository.NewAnakAsuhRepository(db.DB)
		keuanganRepo = repository.NewKeuanganRepository(db.DB)
		jumatBerkahRepo = repository.NewJumatBerkahRepository(db.DB)
		artikelRepo = repository.NewArtikelRepository(db.DB)
		pengaturanRepo = repository.NewPengaturanRepository(db.DB)
		galeriRepo = repository.NewGaleriRepository(db.DB)
		rekeningRepo = repository.NewRekeningDonasiRepository(db.DB)
		defer db.Close()
	}

	authService := services.NewAuthService(pengurusRepo, cfg.JWTSecret)
	anakAsuhService := services.NewAnakAsuhService(anakAsuhRepo)
	keuanganService := services.NewKeuanganService(keuanganRepo)
	jumatBerkahService := services.NewJumatBerkahService(jumatBerkahRepo, anakAsuhRepo)
	artikelService := services.NewArtikelService(artikelRepo)
	pengaturanService := services.NewPengaturanService(pengaturanRepo)
	galeriService := services.NewGaleriService(galeriRepo)
	rekeningService := services.NewRekeningDonasiService(rekeningRepo)
	exportImportService := services.NewExportImportService(anakAsuhService)

	dashboardHandler := handlers.NewDashboardHandler(cfg, anakAsuhService, keuanganService, jumatBerkahService)
	publicHandler := handlers.NewPublicHandler(jumatBerkahService, anakAsuhService, artikelService, keuanganService, pengaturanService)
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := authmw.NewAuthMiddleware(authService)
	anakAsuhHandler := handlers.NewAnakAsuhHandler(anakAsuhService, keuanganService, jumatBerkahService, exportImportService)
	keuanganHandler := handlers.NewKeuanganHandler(keuanganService, anakAsuhService)
	artikelHandler := handlers.NewArtikelHandler(artikelService)
	pengaturanHandler := handlers.NewPengaturanHandler(pengaturanService, rekeningService)
	jumatBerkahHandler := handlers.NewJumatBerkahHandler(jumatBerkahService)
	galeriHandler := handlers.NewGaleriHandler(galeriService, pengaturanService)

	e.GET("/", publicHandler.LandingPage)
	e.GET("/tentang", publicHandler.AboutPage)
	e.GET("/jumat-berkah", publicHandler.JumatBerkahForm)
	e.GET("/program-donasi", publicHandler.ProgramDonasiPage)
	e.GET("/program-donasi/konfirmasi", publicHandler.ProgramDonasiConfirmationPage)
	e.GET("/doa-harian", publicHandler.DoaHarianPage)
	e.GET("/dzikir", publicHandler.DzikirPage)
	e.POST("/api/jumat-berkah/register", publicHandler.SubmitJumatBerkahRegistration)
	e.POST("/api/program-donasi", publicHandler.SubmitProgramDonasi)
	e.POST("/api/program-donasi/confirm", publicHandler.SubmitProgramDonasiConfirmation)
	e.GET("/api/program-donasi/status/:id", publicHandler.GetProgramDonasiStatus)
	e.GET("/api/program-donasi/history", publicHandler.GetProgramDonasiHistory)
	e.GET("/api/jumat-berkah/anak", publicHandler.GetJumatBerkahData)
	e.GET("/zakat", publicHandler.ZakatCalculator)
	e.GET("/zakat/payment", publicHandler.ZakatPayment)
	e.POST("/api/zakat/payment", publicHandler.SubmitZakatPayment)
	e.GET("/zakat/success", publicHandler.ZakatSuccess)
	e.GET("/berita", publicHandler.NewsList)
	e.GET("/berita/:id", publicHandler.NewsDetail)
	e.GET("/galeri", galeriHandler.PublicPage)
	e.POST("/api/newsletter", publicHandler.SubscribeNewsletter)
	e.GET("/offline", publicHandler.OfflinePage)

	e.GET("/admin/login", authHandler.LoginPage)
	e.POST("/admin/login", authHandler.Login)

	admin := e.Group("/admin")
	admin.Use(authMiddleware.RequireAdminSession)
	{
		admin.GET("/dashboard", dashboardHandler.Dashboard)
		admin.GET("/logout", authHandler.Logout)
		admin.GET("/session-info", authmw.SessionInfo)

		admin.GET("/jumat-berkah", jumatBerkahHandler.List)
		admin.POST("/jumat-berkah/:id/approve", jumatBerkahHandler.Approve)
		admin.POST("/jumat-berkah/:id/reject", jumatBerkahHandler.Reject)
		admin.POST("/jumat-berkah/bulk-approve", jumatBerkahHandler.BulkApprove)
		admin.POST("/jumat-berkah/bulk-reject", jumatBerkahHandler.BulkReject)
		admin.POST("/jumat-berkah/approve-all", jumatBerkahHandler.ApproveAll)
		admin.POST("/jumat-berkah/quota", jumatBerkahHandler.UpdateQuota)
		admin.POST("/jumat-berkah/toggle-form", jumatBerkahHandler.ToggleForm)
		admin.POST("/jumat-berkah/manual-register", jumatBerkahHandler.ManualRegister)

		admin.GET("/anak-asuh", anakAsuhHandler.List)
		admin.GET("/anak-asuh/tambah", anakAsuhHandler.Form)
		admin.POST("/anak-asuh", anakAsuhHandler.Create)
		admin.GET("/anak-asuh/:id", anakAsuhHandler.Detail)
		admin.GET("/anak-asuh/:id/edit", anakAsuhHandler.EditForm)
		admin.POST("/anak-asuh/:id", anakAsuhHandler.Update)
		admin.DELETE("/anak-asuh/:id", anakAsuhHandler.Delete)
		admin.GET("/anak-asuh/export/excel", anakAsuhHandler.ExportExcel)
		admin.GET("/anak-asuh/export/csv", anakAsuhHandler.ExportCSV)
		admin.GET("/anak-asuh/export/template", anakAsuhHandler.DownloadTemplate)
		admin.POST("/anak-asuh/import", anakAsuhHandler.ImportData)

		admin.GET("/keuangan", keuanganHandler.BukuKas)
		admin.GET("/keuangan/pemasukan", keuanganHandler.CatatPemasukan)
		admin.POST("/keuangan/pemasukan", keuanganHandler.SavePemasukan)
		admin.POST("/keuangan/pemasukan/:id/verify", keuanganHandler.VerifyPemasukan)
		admin.GET("/keuangan/pengeluaran", keuanganHandler.CatatPengeluaran)
		admin.POST("/keuangan/pengeluaran", keuanganHandler.SavePengeluaran)
		admin.POST("/keuangan/donatur", keuanganHandler.CreateDonatur)
		admin.GET("/keuangan/transaksi/:id", keuanganHandler.GetTransactionDetail)
		admin.PUT("/keuangan/transaksi/:id", keuanganHandler.UpdateTransaction)
		admin.DELETE("/keuangan/transaksi/:id", keuanganHandler.DeleteTransaction)
		admin.GET("/keuangan/edit-form-data", keuanganHandler.GetEditFormData)
		admin.GET("/keuangan/export/csv", keuanganHandler.ExportCSV)
		admin.GET("/keuangan/export/pdf", keuanganHandler.ExportPDF)

		admin.GET("/artikel", artikelHandler.List)
		admin.GET("/artikel/tambah", artikelHandler.Form)
		admin.POST("/artikel", artikelHandler.Create)
		admin.GET("/artikel/:id/edit", artikelHandler.EditForm)
		admin.POST("/artikel/:id", artikelHandler.Update)
		admin.DELETE("/artikel/:id", artikelHandler.Delete)
		admin.POST("/artikel/:id/publish", artikelHandler.Publish)

		admin.GET("/galeri", galeriHandler.AdminPage)
		admin.POST("/galeri", galeriHandler.Upload)
		admin.PUT("/galeri/:id", galeriHandler.Update)
		admin.DELETE("/galeri/:id", galeriHandler.Delete)

		admin.GET("/pengaturan", pengaturanHandler.Page)
		admin.POST("/pengaturan", pengaturanHandler.Save)
		admin.GET("/pengaturan/rekening", pengaturanHandler.ListRekening)
		admin.POST("/pengaturan/rekening", pengaturanHandler.CreateRekening)
		admin.DELETE("/pengaturan/rekening/:id", pengaturanHandler.DeleteRekening)
	}

	port := ":" + cfg.Port
	log.Printf("Server starting on %s", port)
	e.Logger.Fatal(e.Start(port))
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
