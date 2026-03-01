package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"puriyatim-app/internal/models"
	"puriyatim-app/internal/repository"
	"strings"
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
	secret := strings.TrimSpace(jwtSecret)
	if secret == "" || secret == "default_secret_key" {
		secret = generateRuntimeSecret()
		log.Printf("warning: JWT secret default/kosong, menggunakan secret runtime sementara")
	}
	return &AuthService{
		pengurusRepo: pengurusRepo,
		jwtSecret:    secret,
	}
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	if s.pengurusRepo == nil {
		return nil, errors.New("layanan autentikasi belum tersedia")
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)
	if email == "" || password == "" {
		return nil, errors.New("email atau password salah")
	}

	pengurus, err := s.pengurusRepo.GetByEmail(email)
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute).UTC()),
			Issuer:    "puri-yatim-admin",
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
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.Issuer != "puri-yatim-admin" {
			return nil, errors.New("invalid issuer")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func generateRuntimeSecret() string {
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(random)
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
