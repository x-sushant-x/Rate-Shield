package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	DEFAULT_USER_KEY_REDIS = "user:default"
	letterBytes            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type AuthService struct {
	redis *redis.Client
}

func NewAuthService(instance *redis.Client) *AuthService {
	return &AuthService{
		redis: instance,
	}
}

// Initialize default account for login and store them in redis.
// This function only create default account if this is 1st time we are running rate sheild.
// IMPORTANT: default email and password must be changed for security reasons.
func (as *AuthService) InitializeDefaultCreds() {
	ctx := context.Background()

	if as.isAccountExist(ctx) {
		return
	}

	email, pass, hashedPass, err := as.generateDefaultCreds()
	if err != nil {
		return
	}

	err = as.redis.HSet(ctx, DEFAULT_USER_KEY_REDIS, map[string]string{
		"email":    email,
		"password": string(hashedPass),
	}).Err()

	if err != nil {
		log.Err(err).Msgf("failed to store default creds: %v", err)
		return
	}

	credsLog := fmt.Sprintf("✅ IMPORTANT -- Default Credentials -- Email: %s & Password: %s", email, pass)
	log.Info().Msg(credsLog)
}

func (as *AuthService) isAccountExist(ctx context.Context) bool {
	found, err := as.redis.Exists(ctx, DEFAULT_USER_KEY_REDIS).Result()

	if err != nil {
		log.Err(err).Msgf("failed to check existing creds: %v", err)
		return false
	}

	if found > 0 {
		log.Info().Msg("Dashboard credentials already exists. Use them to login.")
		return true
	}

	return false
}

func (as *AuthService) generateDefaultCreds() (string, string, []byte, error) {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))

	pass := make([]byte, 10)
	for i := range pass {
		pass[i] = letterBytes[src.Intn(len(letterBytes))]
	}

	email := "admin@rate-shield.io"
	password := string(pass)

	hashPass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

	if err != nil {
		log.Err(err).Msgf("failed to hash default creds: %v", err)
		return "", "", nil, err
	}

	return email, password, hashPass, err
}
