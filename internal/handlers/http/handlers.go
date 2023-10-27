package http

import (
	"github.com/gofiber/fiber/v2"
	"wbL0/internal/models/params"
	"wbL0/internal/models/response"
)

func (s *Server) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
func (s *Server) PullOrder(ctx *fiber.Ctx) error {
	text := params.PublishOrderForDocker{}
	err := ctx.BodyParser(&text)
	err = s.experiment.PublishOrder(text.Text)
	if err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}
func (s *Server) getOrder(ctx *fiber.Ctx) error {
	order := params.OrderUid{}
	err := ctx.BodyParser(&order)
	if err != nil {
		return err
	}
	result, err := s.experiment.GetOrder(ctx, order)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ошибка в получении заказа",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(
		response.GetOrder{
			OrderUid:  result.OrderUid,
			OrderInfo: result.OrderInfo,
		})
}
