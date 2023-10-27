package consumer

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"wbL0/internal/models"
	"wbL0/internal/storage"
)

type Consumer struct {
	sc        stan.Conn
	repo      storage.Repo
	logger    *zap.Logger
	validator *validator.Validate
}

func NewConsumer(sc stan.Conn, repo storage.Repo, logger *zap.Logger, validator *validator.Validate) *Consumer {
	return &Consumer{
		sc:        sc,
		repo:      repo,
		logger:    logger,
		validator: validator,
	}
}
func (c *Consumer) ListeningMessages() error {
	_, err := c.sc.Subscribe("orders", func(msg *stan.Msg) {
		err := c.SaveMessage(msg.Data)
		if err != nil {
			c.logger.Info("consumer: can`t save order in DB")
		}
	}, stan.DeliverAllAvailable(), stan.DurableName("my-name"))
	if err != nil {
		c.logger.Info("can`t subscribe consumer")
	}
	return nil
}
func (c *Consumer) SaveMessage(data []byte) error {

	order := &models.Order{}
	err := json.Unmarshal(data, order)
	if err != nil {
		return err
	}
	err = c.validator.Struct(order)
	if err != nil {
		c.logger.Error("message in not validate for struct ")
		return err
	}
	uidGenerate := uuid.NewV4()
	err = c.repo.SaveOrder(context.Background(), uidGenerate, data, order.OrderUID)
	if err != nil {
		return err
	}
	c.logger.Info("new order save in data base")
	return nil
}
