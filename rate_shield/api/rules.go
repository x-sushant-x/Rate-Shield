package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/service"
	"github.com/x-sushant-x/RateShield/utils"
)

type RulesAPIHandler struct {
	rulesSvc service.RulesService
}

func NewRulesAPIHandler(svc service.RulesService) RulesAPIHandler {
	return RulesAPIHandler{
		rulesSvc: svc,
	}
}

// extractActorInfo extracts actor information from request headers
// Priority: X-User-ID > Authorization (parsed) > "anonymous"
func extractActorInfo(r *http.Request) string {
	// Check for X-User-ID header first
	userID := r.Header.Get("X-User-ID")
	if userID != "" {
		return userID
	}

	// Check Authorization header and try to extract user info
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Simple extraction: remove "Bearer " prefix if present
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			// Return first 8 characters of token as identifier
			if len(token) > 8 {
				return "token:" + token[:8] + "..."
			}
			return "token:" + token
		}
		return "auth:provided"
	}

	// Default to anonymous
	return "anonymous"
}

// extractIPAddress extracts the client IP address from the request
func extractIPAddress(r *http.Request) string {
	// Check for X-Forwarded-For header (common with proxies/load balancers)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check for X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	// RemoteAddr includes port, so we need to strip it
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

func (h RulesAPIHandler) ListAllRules(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	items := r.URL.Query().Get("items")

	if page != "" && items != "" {
		pageInt, pageIntErr := strconv.Atoi(page)
		itemsInt, itemsIntErr := strconv.Atoi(items)

		if pageIntErr != nil || itemsIntErr != nil {
			utils.BadRequestError(w)
			return
		}

		rules, err := h.rulesSvc.GetPaginatedRules(pageInt, itemsInt)
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}
		utils.SuccessResponse(rules, w)
	} else {
		rules, err := h.rulesSvc.GetAllRules()
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}
		utils.SuccessResponse(rules, w)
	}
}

func (h RulesAPIHandler) SearchRules(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	searchText := q.Get("endpoint")
	if len(searchText) == 0 {
		utils.BadRequestError(w)
		return
	}

	rules, err := h.rulesSvc.SearchRule(searchText)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.SuccessResponse(rules, w)
}

func (h RulesAPIHandler) CreateOrUpdateRule(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// Preflight
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method == http.MethodPost {
		updateReq, err := utils.ParseAPIBody[models.Rule](r)
		if err != nil {
			utils.BadRequestError(w)
			return
		}

		// Extract audit information
		actor := extractActorInfo(r)
		ipAddress := extractIPAddress(r)
		userAgent := r.UserAgent()

		err = h.rulesSvc.CreateOrUpdateRule(updateReq, actor, ipAddress, userAgent)
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}

		utils.SuccessResponse("Rule Created Successfully", w)
	} else {
		utils.MethodNotAllowedError(w)
	}
}

func (h RulesAPIHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		deleteReq, err := utils.ParseAPIBody[models.DeleteRuleDTO](r)
		if err != nil {
			utils.BadRequestError(w)
			return
		}

		// Extract audit information
		actor := extractActorInfo(r)
		ipAddress := extractIPAddress(r)
		userAgent := r.UserAgent()

		err = h.rulesSvc.DeleteRule(deleteReq.RuleKey, actor, ipAddress, userAgent)
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}

		utils.SuccessResponse("Rule Deleted Successfully", w)
	} else {
		utils.MethodNotAllowedError(w)
	}
}
