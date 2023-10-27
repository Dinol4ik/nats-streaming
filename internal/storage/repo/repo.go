package repo

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"wbL0/internal/models"
)

type Repo struct {
	db *sqlx.DB
}

func NewConnect(db *sqlx.DB) Repo {
	return Repo{db}
}

const addSegmentSql = `INSERT INTO public."order"(id, order_uid, order_info) VALUES ($1,$2,$3)`

func (r *Repo) SaveOrder(ctx context.Context, uuid uuid.UUID, data []byte, orderUid string) error {
	_, err := r.db.ExecContext(ctx, addSegmentSql, uuid, orderUid, data)
	if err != nil {
		return err
	}
	return nil
}

const getOrders = `SELECT order_uid,order_info FROM "order"`

func (r *Repo) GetListOrders(ctx context.Context) ([]models.GetOrder, error) {
	var orders []models.GetOrder
	err := r.db.SelectContext(ctx, &orders, getOrders)
	if err != nil {
		return []models.GetOrder{}, err
	}
	return orders, nil
}

const getOrder = `SELECT order_uid,order_info FROM "order" WHERE order_uid = $1`

func (r *Repo) GetOrder(ctx *fiber.Ctx, orderUid string) (models.GetOrder, error) {
	var order models.GetOrder
	err := r.db.GetContext(ctx.Context(), &order, getOrder, orderUid)
	if err != nil {
		return models.GetOrder{}, err
	}
	return order, nil
}
