package usecase

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"wbL0/internal/models"
	"wbL0/internal/models/params"
)

type Experiment interface {
	PublishOrder(text json.RawMessage) error
	GetOrder(ctx *fiber.Ctx, order params.OrderUid) (models.GetOrder, error)
}
