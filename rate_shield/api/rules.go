package api

import (
	"net/http"
	"strconv"

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
	// if r.Method == http.MethodPost {
	updateReq, err := utils.ParseAPIBody[models.Rule](r)
	if err != nil {
		utils.BadRequestError(w)
		return
	}
	err = h.rulesSvc.CreateOrUpdateRule(updateReq)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.SuccessResponse("Rule Created Successfully", w)
	// }

	// else {
	// 	utils.MethodNotAllowedError(w)
	// }
}

func (h RulesAPIHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		deleteReq, err := utils.ParseAPIBody[models.DeleteRuleDTO](r)
		if err != nil {
			utils.BadRequestError(w)
			return
		}

		err = h.rulesSvc.DeleteRule(deleteReq.RuleKey)
		if err != nil {
			utils.InternalError(w, err.Error())
			return
		}

		utils.SuccessResponse("Rule Deleted Successfully", w)
	} else {
		utils.MethodNotAllowedError(w)
	}
}
