package api

import (
	"net/http"
	"strconv"

	"github.com/x-sushant-x/RateShield/service"
	"github.com/x-sushant-x/RateShield/utils"
)

type AuditAPIHandler struct {
	auditSvc service.AuditService
}

func NewAuditAPIHandler(svc service.AuditService) AuditAPIHandler {
	return AuditAPIHandler{
		auditSvc: svc,
	}
}

// ListAuditLogs handles GET /audit/logs
// Supports pagination: ?page=1&items=10
// Supports filtering: ?endpoint=/api/v1/test&actor=user@example.com&action=CREATE
func (h AuditAPIHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.MethodNotAllowedError(w)
		return
	}

	query := r.URL.Query()

	// Check for filters
	endpoint := query.Get("endpoint")
	actor := query.Get("actor")
	action := query.Get("action")

	// If filters are provided, return filtered results
	if endpoint != "" {
		logs, err := h.auditSvc.GetAuditLogsByEndpoint(endpoint)
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}
		utils.SuccessResponse(logs, w)
		return
	}

	if actor != "" {
		logs, err := h.auditSvc.GetAuditLogsByActor(actor)
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}
		utils.SuccessResponse(logs, w)
		return
	}

	if action != "" {
		logs, err := h.auditSvc.GetAuditLogsByAction(action)
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}
		utils.SuccessResponse(logs, w)
		return
	}

	// Check for pagination parameters
	page := query.Get("page")
	items := query.Get("items")

	if page != "" && items != "" {
		pageInt, pageErr := strconv.Atoi(page)
		itemsInt, itemsErr := strconv.Atoi(items)

		if pageErr != nil || itemsErr != nil {
			utils.BadRequestError(w)
			return
		}

		logs, err := h.auditSvc.GetAuditLogs(pageInt, itemsInt)
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}
		utils.SuccessResponse(logs, w)
		return
	}

	// No pagination or filters - return all logs
	logs, err := h.auditSvc.GetAllAuditLogs()
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.SuccessResponse(logs, w)
}
