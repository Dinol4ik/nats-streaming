package storage

import (
	"context"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"wbL0/internal/models"
)

type Repo interface {
	SaveOrder(ctx context.Context, uuid uuid.UUID, data []byte, orderUid string) error
	GetListOrders(ctx context.Context) ([]models.GetOrder, error)
	GetOrder(ctx *fiber.Ctx, orderUid string) (models.GetOrder, error)
}
