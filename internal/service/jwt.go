//go:generate mockgen -source=jwt.go -package service -destination=jwt_mock.go
package service

import (
	"fmt"
	"mathly/internal/models"
	"time"

	jwtLib "github.com/golang-jwt/jwt"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

func (t TokenType) toString() string {
	switch t {
	case Access:
		return "access"
	case Refresh:
		return "refresh"
	}
	return "undefined"
}

type JwtConfig struct {
	Secret string `json:"secret"`
}

type JWT interface {
	GenerateToken(user *models.User, tokenType TokenType) (string, error)
	ValidateToken(tokenString string, tokenType TokenType) (*models.Claims, error)
}

type jwt struct {
	config JwtConfig
}

func newJWT(config JwtConfig) JWT {
	return &jwt{
		config: config,
	}
}

func (j *jwt) GenerateToken(user *models.User, tokenType TokenType) (string, error) {
	var expirationTime int64
	var subject string

	if tokenType == Access {
		expirationTime = time.Now().Add(time.Hour * 1).Unix() // expire in 1 hour
	} else if tokenType == Refresh {
		expirationTime = time.Now().Add(time.Hour * 24 * 7).Unix() // expire in 7 days
		subject = "refresh"
	}

	claims := &models.Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		StandardClaims: jwtLib.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    "mathly",
			Subject:   subject,
		},
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(j.config.Secret))

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (j *jwt) ValidateToken(tokenString string, tokenType TokenType) (*models.Claims, error) {
	token, err := jwtLib.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwtLib.Token) (any, error) {
		if _, ok := token.Method.(*jwtLib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(j.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if tokenType == Refresh && claims.Subject != string(Refresh) {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
