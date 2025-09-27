package jwt

import (
	"time"

	"github.com/cctw-zed/wonder/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

// TokenService provides JWT token management
type TokenService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	GetSigningKey() []byte
}

// Claims represents JWT token claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTService implements TokenService
type JWTService struct {
	signingKey []byte
	expiry     time.Duration
}

// NewTokenService creates a new JWT token service
func NewTokenService(signingKey string, expiry time.Duration) TokenService {
	return &JWTService{
		signingKey: []byte(signingKey),
		expiry:     expiry,
	}
}

// GenerateToken generates a JWT token for the given user ID
func (j *JWTService) GenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.NewRequiredFieldError("user_id", userID)
	}

	// Create claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wonder-api",
			Subject:   userID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(j.signingKey)
	if err != nil {
		return "", errors.NewBusinessLogicError("token_generation", "failed to sign JWT token")
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.NewRequiredFieldError("token", tokenString)
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewUnauthorizedError("token_validation", "", "invalid signing method")
		}
		return j.signingKey, nil
	})

	if err != nil {
		return nil, errors.NewUnauthorizedError("token_validation", "", "invalid token")
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.NewUnauthorizedError("token_validation", "", "invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.NewUnauthorizedError("token_validation", "", "invalid token claims")
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.NewUnauthorizedError("token_validation", "", "token expired")
	}

	return claims, nil
}

// GetSigningKey returns the signing key (for testing purposes)
func (j *JWTService) GetSigningKey() []byte {
	return j.signingKey
}
