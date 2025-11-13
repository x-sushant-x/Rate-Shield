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

	if as.doesAccountExist(ctx, DEFAULT_USER_KEY_REDIS) {
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

func (as *AuthService) doesAccountExist(ctx context.Context, accountKey string) bool {
	found, err := as.redis.Exists(ctx, accountKey).Result()

	if err != nil {
		log.Err(err).Msgf("failed to check existing creds: %v", err)
		return false
	}

	if found > 0 {
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

	email := "default"
	password := string(pass)

	hashPass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)

	if err != nil {
		log.Err(err).Msgf("failed to hash default creds: %v", err)
		return "", "", nil, err
	}

	return email, password, hashPass, err
}

func (as *AuthService) LoginUser(email, pass string) error {
	ctx := context.Background()

	redisKey := "user:" + email

	if !as.doesAccountExist(ctx, redisKey) {
		return fmt.Errorf("account does not exists with email: %s", email)
	}

	res, err := as.redis.HMGet(ctx, redisKey, "password").Result()
	if err != nil || res == nil || len(res) == 0 {
		return fmt.Errorf("unable to verify your account")
	}

	password := res[0].(string)

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(pass))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	// Serve JWT Token

	return nil
}

func (as *AuthService) CreateUser(email, pass string) error {
	ctx := context.Background()

	redisKey := "user:" + email

	if as.doesAccountExist(ctx, redisKey) {
		return fmt.Errorf("account already exists with email: %s", email)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		log.Err(err).Msgf("failed to hash default creds: %v", err)
		return fmt.Errorf("unable to create user account")
	}

	err = as.redis.HSet(ctx, DEFAULT_USER_KEY_REDIS, map[string]string{
		"email":    email,
		"password": string(hashedPass),
	}).Err()

	if err != nil {
		log.Err(err).Msgf("failed to store default creds: %v", err)
		return fmt.Errorf("unable to create user account")
	}

	return nil
}
