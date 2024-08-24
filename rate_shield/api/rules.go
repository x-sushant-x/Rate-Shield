package api

import (
	"net/http"

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

func (h RulesAPIHandler) ListAllRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.rulesSvc.GetAllRules()
	if err != nil {
		utils.InternalError(w)
	}
	utils.SuccessResponse(rules, w)
}

func (h RulesAPIHandler) CreateOrUpdateRule(w http.ResponseWriter, r *http.Request) {
	updateReq, err := utils.ParseAPIBody[models.Rule](r)
	if err != nil {
		utils.BadRequestError(w)
	}

	err = h.rulesSvc.CreateOrUpdateRule(updateReq)
	if err != nil {
		utils.InternalError(w)
	}

	utils.SuccessResponse("Rule Created Successfully", w)
}

func (h RulesAPIHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	deleteReq, err := utils.ParseAPIBody[models.DeleteRuleDTO](r)
	if err != nil {
		utils.BadRequestError(w)
	}

	err = h.rulesSvc.DeleteRule(deleteReq.RuleKey)
	if err != nil {
		utils.InternalError(w)
	}

	utils.SuccessResponse("Rule Deleted Successfully", w)
}
