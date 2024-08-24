package api

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/limiter"
	"github.com/x-sushant-x/RateShield/service"
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
	mux := http.NewServeMux()

	s.registerRulesRoutes(mux)
	s.registerRateLimiterRoutes(mux)

	corsMux := s.setupCORS(mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: corsMux,
	}

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

		h.ServeHTTP(w, r)
	})
}

func (s Server) registerRulesRoutes(mux *http.ServeMux) {
	rulesSvc := service.RulesServiceRedis{}
	rulesHandler := NewRulesAPIHandler(rulesSvc)

	mux.HandleFunc("/rule/list", rulesHandler.ListAllRules)
	mux.HandleFunc("/rule/add", rulesHandler.CreateOrUpdateRule)
	mux.HandleFunc("/rule/delete", rulesHandler.DeleteRule)
}

func (s Server) registerRateLimiterRoutes(mux *http.ServeMux) {
	tokenBucketSvc := limiter.NewTokenBucketService()
	limiter := limiter.NewRateLimiterService(tokenBucketSvc)
	rateLimiterHandler := NewRateLimitHandler(limiter)

	mux.HandleFunc("/check-limit", rateLimiterHandler.CheckRateLimit)
}
