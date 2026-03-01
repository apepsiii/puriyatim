package middleware

import (
	"net/http"
	"puriyatim-app/internal/models"
	"puriyatim-app/internal/services"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware to check if user is authenticated
func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := extractToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": err.Error(),
			})
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
		}

		setUserContext(c, claims)

		return next(c)
	}
}

func (m *AuthMiddleware) RequireAdminSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := extractToken(c)
		if err != nil {
			return unauthorized(c)
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			return unauthorized(c)
		}

		setUserContext(c, claims)
		return next(c)
	}
}

// RequireRole middleware to check if user has specific role
func (m *AuthMiddleware) RequireRole(roles ...models.PeranPengurus) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role from context
			userRole, ok := c.Get("user_role").(models.PeranPengurus)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "user role not found",
				})
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "insufficient permissions",
				})
			}

			return next(c)
		}
	}
}

// RequireSuperadmin middleware to check if user is superadmin
func (m *AuthMiddleware) RequireSuperadmin() echo.MiddlewareFunc {
	return m.RequireRole(models.PeranSuperadmin)
}

// RequireKeuangan middleware to check if user is keuangan or superadmin
func (m *AuthMiddleware) RequireKeuangan() echo.MiddlewareFunc {
	return m.RequireRole(models.PeranKeuangan, models.PeranSuperadmin)
}

// RequirePenulisBerita middleware to check if user is penulis berita or superadmin
func (m *AuthMiddleware) RequirePenulisBerita() echo.MiddlewareFunc {
	return m.RequireRole(models.PeranPenulisBerita, models.PeranSuperadmin)
}

func extractToken(c echo.Context) (string, error) {
	authHeader := strings.TrimSpace(c.Request().Header.Get("Authorization"))
	if authHeader != "" {
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
		}
		return tokenParts[1], nil
	}

	cookie, err := c.Cookie("token")
	if err == nil && cookie != nil && strings.TrimSpace(cookie.Value) != "" {
		return strings.TrimSpace(cookie.Value), nil
	}

	return "", echo.NewHTTPError(http.StatusUnauthorized, "authorization required")
}

func setUserContext(c echo.Context, claims *services.Claims) {
	c.Set("user_id", claims.ID)
	c.Set("user_email", claims.Email)
	c.Set("user_role", claims.Peran)
	c.Set("user_role_str", string(claims.Peran))
}

func unauthorized(c echo.Context) error {
	accept := strings.ToLower(c.Request().Header.Get("Accept"))
	requestedWith := strings.ToLower(c.Request().Header.Get("X-Requested-With"))
	isJSON := strings.Contains(accept, "application/json") ||
		strings.HasPrefix(c.Path(), "/api/") ||
		requestedWith == "xmlhttprequest" ||
		c.Request().Method != http.MethodGet
	if isJSON {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}
	return c.Redirect(http.StatusFound, "/admin/login?error=Sesi berakhir, silakan login kembali")
}
