package cache

import (
	"context"
	"encoding/json"
	lru "github.com/hashicorp/golang-lru/v2"
	"go.uber.org/zap"
	"time"
	"wbL0/internal/storage"
)

type Cache struct {
	cache  *lru.Cache[string, json.RawMessage]
	size   int
	repo   storage.Repo
	logger *zap.Logger
}

func NewCache(size int, repo storage.Repo, logger *zap.Logger) *Cache {
	c, err := lru.New[string, json.RawMessage](size)
	if err != nil {
		logger.Fatal("can`t create cache ")
	}
	return &Cache{
		cache:  c,
		size:   size,
		repo:   repo,
		logger: logger,
	}
}

func (c *Cache) Recover(timeout int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	orders, err := c.repo.GetListOrders(ctx)
	if err != nil {
		return err
	}
	for _, o := range orders {
		c.cache.ContainsOrAdd(o.OrderUid, o.OrderInfo)
	}
	return nil
}
func (c *Cache) GetValue(uuid string) (json.RawMessage, bool) {
	order, ok := c.cache.Get(uuid)
	if !ok {
		return json.RawMessage{}, ok
	}
	return order, ok
}
