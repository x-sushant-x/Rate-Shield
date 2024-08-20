package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x-sushant-x/RateShield/service"
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
	return nil
}
