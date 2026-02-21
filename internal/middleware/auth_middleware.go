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
		// Get token from header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "authorization header required",
			})
		}

		// Check if token format is correct
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid authorization header format",
			})
		}

		// Validate token
		claims, err := m.authService.ValidateToken(tokenParts[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
		}

		// Set user context
		c.Set("user_id", claims.ID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Peran)

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