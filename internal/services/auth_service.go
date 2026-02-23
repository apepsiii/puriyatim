package services

import (
	"errors"
	"fmt"
	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	pengurusRepo *repository.PengurusRepository
	jwtSecret    string
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token    string           `json:"token"`
	Pengurus *models.Pengurus `json:"pengurus"`
}

type Claims struct {
	ID    string               `json:"id"`
	Email string               `json:"email"`
	Peran models.PeranPengurus `json:"peran"`
	jwt.RegisteredClaims
}

func NewAuthService(pengurusRepo *repository.PengurusRepository, jwtSecret string) *AuthService {
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
	}
	return &AuthService{
		pengurusRepo: pengurusRepo,
		jwtSecret:    jwtSecret,
	}
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	if s.pengurusRepo == nil {
		if req.Email == "admin@puriyatim.org" && req.Password == "admin123" {
			pengurus := &models.Pengurus{
				ID:          "1",
				NamaLengkap: "Budi Admin",
				Email:       "admin@puriyatim.org",
				Peran:       models.PeranSuperadmin,
				Status:      models.StatusPengurusAktif,
			}
			token, err := s.generateToken(pengurus)
			if err != nil {
				return nil, fmt.Errorf("failed to generate token: %w", err)
			}
			return &LoginResponse{
				Token:    token,
				Pengurus: pengurus,
			}, nil
		}
		return nil, errors.New("email atau password salah")
	}

	pengurus, err := s.pengurusRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if pengurus.Status != models.StatusPengurusAktif {
		return nil, errors.New("akun tidak aktif")
	}

	err = bcrypt.CompareHashAndPassword([]byte(pengurus.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	token, err := s.generateToken(pengurus)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	pengurus.PasswordHash = ""

	return &LoginResponse{
		Token:    token,
		Pengurus: pengurus,
	}, nil
}

func (s *AuthService) generateToken(pengurus *models.Pengurus) (string, error) {
	// Create claims
	claims := &Claims{
		ID:    pengurus.ID,
		Email: pengurus.Email,
		Peran: pengurus.Peran,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *AuthService) CreatePengurus(pengurus *models.Pengurus, password string) error {
	// Hash password
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Set password hash
	pengurus.PasswordHash = hashedPassword

	// Create pengurus
	err = s.pengurusRepo.Create(pengurus)
	if err != nil {
		return fmt.Errorf("failed to create pengurus: %w", err)
	}

	return nil
}

func (s *AuthService) ChangePassword(pengurusID, currentPassword, newPassword string) error {
	// Get pengurus
	pengurus, err := s.pengurusRepo.GetByID(pengurusID)
	if err != nil {
		return fmt.Errorf("pengurus not found: %w", err)
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(pengurus.PasswordHash), []byte(currentPassword))
	if err != nil {
		return errors.New("password saat ini salah")
	}

	// Hash new password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	err = s.pengurusRepo.UpdatePassword(pengurusID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
