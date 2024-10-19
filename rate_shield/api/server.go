/*
	This file needs a lot of revamping. Dependencies are getting out of hand.
*/

package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/limiter"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/service"
)

var (
	slackToken     string
	slackChannelID string
)

type Server struct {
	port int
}

func NewServer(port int) Server {
	return Server{
		port: port,
	}
}

func (s Server) StartServer() error {
	log.Info().Msg("Setting Up API Endpoints ✅")
	mux := http.NewServeMux()

	s.rulesRoutes(mux)
	s.registerRateLimiterRoutes(mux)

	corsMux := s.setupCORS(mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: corsMux,
	}

	log.Info().Msg("Rate Shield Running ✅")

	err := server.ListenAndServe()
	if err != nil {
		log.Err(err).Msg("unable to start server")
		return err
	}
	return nil
}

func (s Server) setupCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (s Server) rulesRoutes(mux *http.ServeMux) {
	redisRuleClient, err := redisClient.NewRulesClient()
	if err != nil {
		log.Err(err).Msg("unable to setup new redis rules client")
		log.Fatal()
	}

	rulesSvc := service.NewRedisRulesService(redisRuleClient)
	rulesHandler := NewRulesAPIHandler(rulesSvc)

	mux.HandleFunc("/rule/list", rulesHandler.ListAllRules)
	mux.HandleFunc("/rule/add", rulesHandler.CreateOrUpdateRule)
	mux.HandleFunc("/rule/delete", rulesHandler.DeleteRule)
	mux.HandleFunc("/rule/search", rulesHandler.SearchRules)
}

func (s Server) registerRateLimiterRoutes(mux *http.ServeMux) {

	tokenBucketSvc := limiter.NewTokenBucketService(getRedisTokenBucket(), getErrorNotificationSvc())
	fixedWindowSvc := limiter.NewFixedWindowService(getRedisFixedWindowClient())

	limiter := limiter.NewRateLimiterService(&tokenBucketSvc, &fixedWindowSvc, getRedisRulesSvc())
	rateLimiterHandler := NewRateLimitHandler(limiter)

	mux.HandleFunc("/check-limit", rateLimiterHandler.CheckRateLimit)
}

func getRedisTokenBucket() redisClient.RedisTokenBucket {
	redisTokenBucket, err := redisClient.NewTokenBucketClient()
	if err != nil {
		log.Fatal().Err(err)
	}
	return redisTokenBucket
}

func getRedisFixedWindowClient() redisClient.RedisFixedWindow {
	redisFixedWindow, err := redisClient.NewFixedWindowClient()
	if err != nil {
		log.Fatal().Err(err)
	}
	return redisFixedWindow
}

func getRedisRulesSvc() service.RulesServiceRedis {
	redisRulesClient, err := redisClient.NewRulesClient()
	if err != nil {
		log.Fatal().Err(err)
	}
	return service.NewRedisRulesService(redisRulesClient)
}

func getErrorNotificationSvc() service.ErrorNotificationSVC {
	loadENVFile()
	setSlackCredentials()

	slackSvc := service.NewSlackService(slackToken, slackChannelID)

	return service.NewErrorNotificationSVC(*slackSvc)
}

func loadENVFile() {
	err := godotenv.Load()
	if err != nil {
		log.Panic().Msgf("error while loading env file: %s", err)
	}
}

func setSlackCredentials() {
	sToken := os.Getenv("SLACK_TOKEN")
	if len(sToken) == 0 {
		log.Panic().Msg("SLACK_TOKEN not available in env file")
	}
	slackToken = sToken

	sChannel := os.Getenv("SLACK_CHANNEL")
	if len(sChannel) == 0 {
		log.Panic().Msg("SLACK_CHANNEL not available in env file")
	}
	slackChannelID = sChannel
}
