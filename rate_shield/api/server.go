package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/limiter"
	redisClient "github.com/x-sushant-x/RateShield/redis"
	"github.com/x-sushant-x/RateShield/service"
)

type Server struct {
	port    int
	limiter *limiter.Limiter
}

func NewServer(limiter *limiter.Limiter) Server {
	return Server{
		port:    getPort(),
		limiter: limiter,
	}
}

func (s Server) StartServer() error {
	log.Info().Msgf("Setting Up API endpoints in port: %d ✅", s.port)
	mux := http.NewServeMux()

	s.rulesRoutes(mux)
	s.auditRoutes(mux)
	s.registerRateLimiterRoutes(mux)
	s.setupHome(mux)

	corsMux := s.setupCORS(mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: corsMux,
	}

	log.Info().Msg("Rate Shield running on port: " + fmt.Sprintf("%d", s.port) + " ✅")

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

	// Create audit client and service
	auditClient := redisClient.NewAuditClient(redisRuleClient.(redisClient.RedisRules).GetClient())
	auditSvc := service.NewAuditService(auditClient)

	// Create rules service with audit service
	rulesSvc := service.NewRedisRulesService(redisRuleClient, auditSvc)
	rulesHandler := NewRulesAPIHandler(rulesSvc)

	mux.HandleFunc("/rule/list", rulesHandler.ListAllRules)
	mux.HandleFunc("/rule/add", rulesHandler.CreateOrUpdateRule)
	mux.HandleFunc("/rule/delete", rulesHandler.DeleteRule)
	mux.HandleFunc("/rule/search", rulesHandler.SearchRules)
}

func (s Server) auditRoutes(mux *http.ServeMux) {
	redisRuleClient, err := redisClient.NewRulesClient()
	if err != nil {
		log.Err(err).Msg("unable to setup new redis rules client for audit")
		log.Fatal()
	}

	// Create audit client and service
	auditClient := redisClient.NewAuditClient(redisRuleClient.(redisClient.RedisRules).GetClient())
	auditSvc := service.NewAuditService(auditClient)
	auditHandler := NewAuditAPIHandler(auditSvc)

	mux.HandleFunc("/audit/logs", auditHandler.ListAuditLogs)
}

func (s Server) registerRateLimiterRoutes(mux *http.ServeMux) {
	rateLimiterHandler := NewRateLimitHandler(s.limiter)
	mux.HandleFunc("/check-limit", rateLimiterHandler.CheckRateLimit)
}

func (s Server) setupHome(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		wd, wdError := os.Getwd()

		homepage, err := os.ReadFile(wd + "/static/" + "index.html")
		if err != nil || wdError != nil {
			fmt.Println(err)
			w.Write([]byte("Rate Shield is running. Open frontend client on port 5173. If it does not work make sure react application is running."))
		}

		fmt.Fprint(w, string(homepage))
	})
}

func getPort() int {
	port := os.Getenv("RATE_SHIELD_PORT")
	if len(port) == 0 {
		log.Fatal().Msg("RATE_SHIELD_PORT environment variable not provided in docker run command.")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal().Msg("Invalid port number provided.")
	}

	return portInt
}
