/*
	This file needs a lot of revamping. Dependencies are getting out of hand.
*/

package api

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/limiter"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/service"
)

type Server struct {
	port    int
	limiter limiter.Limiter
}

func NewServer(port int, limiter limiter.Limiter) Server {
	return Server{
		port:    port,
		limiter: limiter,
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
	rateLimiterHandler := NewRateLimitHandler(s.limiter)

	mux.HandleFunc("/check-limit", rateLimiterHandler.CheckRateLimit)
}
