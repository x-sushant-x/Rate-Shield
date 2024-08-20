package api

import (
	"github.com/gofiber/fiber/v2"
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

func (h RulesAPIHandler) ListAllRules(c *fiber.Ctx) error {
	rules, err := h.rulesSvc.GetAllRules()
	if err != nil {
		return utils.SendInternalError(c)
	}

	return c.Status(200).JSON(map[string]interface{}{
		"status": "sucess",
		"data":   rules,
	})
}

func (h RulesAPIHandler) UpdateRule(c *fiber.Ctx) error {
	var updateReq models.Rule
	if err := c.BodyParser(&updateReq); err != nil {
		return utils.SendBadRequestError(c)
	}

	err := h.rulesSvc.CreateOrUpdateRule(updateReq)
	if err != nil {
		return utils.SendInternalError(c)
	}

	return c.Status(200).JSON(map[string]interface{}{
		"status": "sucess",
		"data":   "Rule Created Successfully",
	})
}

func (h RulesAPIHandler) DeleteRule(c *fiber.Ctx) error {
	var deleteReq models.DeleteRuleDTO

	if err := c.BodyParser(&deleteReq); err != nil {
		return utils.SendBadRequestError(c)
	}

	err := h.rulesSvc.DeleteRule(deleteReq.RuleKey)
	if err != nil {
		return utils.SendInternalError(c)
	}

	return c.Status(200).JSON(map[string]interface{}{
		"status": "sucess",
		"data":   "Rule Deleted Successfully",
	})
}
