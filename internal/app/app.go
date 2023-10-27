package app

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
	"wbL0/internal/app/config"
	"wbL0/internal/cache"
	"wbL0/internal/consumer"
	"wbL0/internal/handlers/http"
	"wbL0/internal/storage/repo"
	"wbL0/internal/usecase/experiment"
)

type App struct {
	cfg *config.Config
}

func (a *App) Run() error {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zapcore.Level(*a.cfg.Logger.Level))
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		os.Stdout,
		atom,
	)
	logger := zap.New(zapCore)
	logger = logger.With(zap.String("service", "experiment"))
	log := logger.Sugar()
	atom.SetLevel(zapcore.Level(*a.cfg.Logger.Level))
	log.Infof("logger initialized successfully")
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", a.cfg.Postgres.Host, a.cfg.Postgres.User, a.cfg.Postgres.Password, a.cfg.Postgres.DBName)
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("DB initialized successfully ")

	if err != nil {
		log.Fatal("can`t close stan")
	}
	repository := repo.NewConnect(db)
	sc, err := stan.Connect(a.cfg.NatsServer.ClusterId, a.cfg.NatsServer.ClientId, stan.NatsURL(a.cfg.NatsServer.NatsUrl), stan.MaxPubAcksInflight(1000))
	if err != nil {
		log.Fatal("cat't connect to stan", zap.Error(err))
	}
	msgValidator := validator.New()
	newConsumer := consumer.NewConsumer(sc, &repository, logger, msgValidator)
	err = newConsumer.ListeningMessages()
	if err != nil {
		log.Fatal("error in listening messages for consumer")
	}
	newCache := cache.NewCache(a.cfg.HTTPServer.CacheSize, &repository, logger)
	_ = newCache.Recover(a.cfg.HTTPServer.CacheFillTimeout)
	wb := experiment.NewWildberriesCase(&repository, log, newCache)

	httpServer := http.NewServer(*a.cfg, *log, wb)
	log.Info("application has started")
	go httpServer.Run()

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	log.Debug("waiting for httpServer to shut down")

	log.Info("application has been shut down")

	return nil
}

func New(cfg *config.Config) *App {
	return &App{cfg}
}
