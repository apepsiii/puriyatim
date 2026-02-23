package handlers

import (
	"net/http"
	"puriyatim-app/internal/services"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
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
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Redirect(http.StatusFound, "/admin/login?error=Email dan password harus diisi")
	}

	req := &services.LoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := h.authService.Login(req)
	if err != nil {
		return c.Redirect(http.StatusFound, "/admin/login?error=Email atau password salah")
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    response.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	userEmail := c.Get("user_email").(string)
	userRole := c.Get("user_role").(string)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    userID,
		"email": userEmail,
		"role":  userRole,
	})
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID := c.Get("user_id").(string)

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

	if len(req.NewPassword) < 6 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password baru minimal 6 karakter",
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
	// In a stateless JWT implementation, logout is typically handled client-side
	// by simply discarding the token. However, we can provide a response to confirm
	// the logout action.
	return c.JSON(http.StatusOK, map[string]string{
		"message": "logout berhasil",
	})
}
