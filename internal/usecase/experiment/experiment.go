package experiment

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"wbL0/internal/cache"
	"wbL0/internal/models"
	"wbL0/internal/models/params"
	"wbL0/internal/storage"
)

type UseCase struct {
	repo   storage.Repo
	logger *zap.SugaredLogger
	cache  *cache.Cache
}

func NewWildberriesCase(r storage.Repo, logger *zap.SugaredLogger, cacheServer *cache.Cache) *UseCase {
	return &UseCase{
		repo:   r,
		logger: logger,
		cache:  cacheServer,
	}
}
func (uc *UseCase) PublishOrder(text json.RawMessage) error {
	//jsonFile, err := os.Open("model.json")
	//if err != nil {
	//	return err
	//}
	//reader, err := io.ReadAll(jsonFile)
	//if err != nil {
	//	uc.logger.Info("can`t read json file")
	//}
	//jsonFile.Close()

	// use your stan.NatsUrl as in config.yaml if u want local launch, default - is name container in docker
	sc, err := stan.Connect("test-cluster", "client-132", stan.NatsURL("nats-streaming"), stan.MaxPubAcksInflight(1000))
	if err != nil {
		uc.logger.Fatalln("erorr connect for publish")
	}
	err = sc.Publish("orders", text)
	if err != nil {
		uc.logger.Info("err publish message")
	}

	sc.Close()
	return nil
}
func (uc *UseCase) GetOrder(ctx *fiber.Ctx, order params.OrderUid) (models.GetOrder, error) {
	orderCache, ok := uc.cache.GetValue(order.OrderUid)
	if ok {
		return models.GetOrder{
			OrderUid:  order.OrderUid,
			OrderInfo: orderCache,
		}, nil
	}
	getOrder, err := uc.repo.GetOrder(ctx, order.OrderUid)
	if err != nil {
		return models.GetOrder{}, err
	}
	return getOrder, nil
}
