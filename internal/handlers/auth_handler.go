package handlers

import (
	"net/http"
	"puriyatim-app/internal/services"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
	limiter     *loginAttemptLimiter
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		limiter:     newLoginAttemptLimiter(),
	}
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

func (h *AuthHandler) LoginPage(c echo.Context) error {
	data := map[string]interface{}{
		"Title": "Login Admin",
		"Error": c.QueryParam("error"),
	}
	return c.Render(http.StatusOK, "admin/login.html", data)
}

func (h *AuthHandler) Login(c echo.Context) error {
	email := strings.ToLower(strings.TrimSpace(c.FormValue("email")))
	password := strings.TrimSpace(c.FormValue("password"))
	clientIP := strings.TrimSpace(c.RealIP())
	limiterKey := clientIP + "|" + email

	if email == "" || password == "" {
		return c.Redirect(http.StatusFound, "/admin/login?error=Email dan password harus diisi")
	}
	if !h.limiter.allow(limiterKey) {
		return c.Redirect(http.StatusFound, "/admin/login?error=Terlalu banyak percobaan login. Coba lagi beberapa menit.")
	}

	req := &services.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := h.authService.Login(req)
	if err != nil {
		h.limiter.fail(limiterKey)
		return c.Redirect(http.StatusFound, "/admin/login?error=Email atau password salah")
	}
	h.limiter.success(limiterKey)

	remember := c.FormValue("remember") == "on"
	cookieDuration := 8 * time.Hour
	if remember {
		cookieDuration = 7 * 24 * time.Hour
		// Buat ulang token dengan durasi 7 hari
		longToken, err := h.authService.GenerateTokenWithDuration(response.Pengurus, cookieDuration)
		if err == nil {
			response.Token = longToken
		}
	}

	SetAuthCookie(c, response.Token, cookieDuration)
	c.Response().Header().Set("Cache-Control", "no-store")

	return c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID, _ := c.Get("user_id").(string)
	userEmail, _ := c.Get("user_email").(string)
	userRole, _ := c.Get("user_role").(string)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    userID,
		"email": userEmail,
		"role":  userRole,
	})
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	// Validate request
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password saat ini dan password baru harus diisi",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password baru minimal 8 karakter",
		})
	}

	// Change password
	err := h.authService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "password berhasil diubah",
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	ClearAuthCookie(c)

	if strings.Contains(strings.ToLower(c.Request().Header.Get("Accept")), "application/json") {
		return JSONOk(c, "logout berhasil")
	}
	return c.Redirect(http.StatusFound, "/admin/login")
}

type loginAttemptLimiter struct {
	mu       sync.Mutex
	attempts map[string]attemptState
}

type attemptState struct {
	Count        int
	FirstAttempt time.Time
	BlockedUntil time.Time
}

func newLoginAttemptLimiter() *loginAttemptLimiter {
	return &loginAttemptLimiter{
		attempts: make(map[string]attemptState),
	}
}

func (l *loginAttemptLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.cleanup(now)

	state, ok := l.attempts[key]
	if !ok {
		return true
	}
	if state.BlockedUntil.After(now) {
		return false
	}
	return true
}

func (l *loginAttemptLimiter) fail(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	state := l.attempts[key]
	if state.FirstAttempt.IsZero() || now.Sub(state.FirstAttempt) > 15*time.Minute {
		state = attemptState{
			Count:        1,
			FirstAttempt: now,
		}
		l.attempts[key] = state
		return
	}

	state.Count++
	if state.Count >= 5 {
		state.BlockedUntil = now.Add(15 * time.Minute)
	}
	l.attempts[key] = state
}

func (l *loginAttemptLimiter) success(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
}

func (l *loginAttemptLimiter) cleanup(now time.Time) {
	for key, state := range l.attempts {
		if state.BlockedUntil.IsZero() && now.Sub(state.FirstAttempt) > 30*time.Minute {
			delete(l.attempts, key)
			continue
		}
		if !state.BlockedUntil.IsZero() && now.After(state.BlockedUntil.Add(30*time.Minute)) {
			delete(l.attempts, key)
		}
	}
}

// isSecureRequest adalah shim ke IsSecureRequest (helpers.go).
func isSecureRequest(c echo.Context) bool { return IsSecureRequest(c) }
