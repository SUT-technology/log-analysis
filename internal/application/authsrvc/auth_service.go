package authsrvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SUT-technology/log-analysis/internal/domain/models"
	cockroachdb "github.com/SUT-technology/log-analysis/internal/infrastructure/cockroachDB"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthSrvc struct {
	cockroachdb *cockroachdb.CockroachDBClient
	secret_key  string
}

func NewAuthSrvc(cockroachdb *cockroachdb.CockroachDBClient, secret_key string) AuthSrvc {
	return AuthSrvc{
		cockroachdb: cockroachdb,
		secret_key:  secret_key,
	}
}

func generateToken(userID uuid.UUID, username string, secret_key string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &models.JWTClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret_key))
}

func (c AuthSrvc) Login(ctx context.Context, username, password string) (string, error) {
	user, err := c.cockroachdb.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("get user fom db: %w", err)
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	// Generate token
	token, err := generateToken(user.ID, user.Username, c.secret_key)
	if err != nil {
		return "", err
	}

	return token, nil
}
