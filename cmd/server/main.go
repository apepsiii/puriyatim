package main

import (
	"html/template"
	"io"
	"log"
	"puriyatim-app/internal/config"
	"puriyatim-app/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	// Template renderer
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/**/*.html")),
	}

	// Static files
	e.Static("/static", "static")

	// Initialize handlers
	dashboardHandler := handlers.NewDashboardHandler(cfg)

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to Puri Yatim!")
	})

	// Admin routes
	admin := e.Group("/admin")
	{
		admin.GET("/dashboard", dashboardHandler.Dashboard)
		admin.POST("/jumat-berkah/approve/:id", dashboardHandler.ApproveJumatBerkah)
		admin.POST("/jumat-berkah/reject/:id", dashboardHandler.RejectJumatBerkah)
		admin.POST("/jumat-berkah/approve-all", dashboardHandler.ApproveAllJumatBerkah)
	}

	// Start server
	port := ":" + cfg.Port
	log.Printf("Server starting on %s", port)
	e.Logger.Fatal(e.Start(port))
}

// TemplateRenderer is a custom renderer for Echo
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}