package main

import (
	"context"

	"github.com/POMBNK/linktgBot/config"
	"github.com/POMBNK/linktgBot/logging"
	tgClient "github.com/POMBNK/linktgBot/pkg/clients/telegram"
	eventconsumer "github.com/POMBNK/linktgBot/pkg/consumer/eventConsumer"
	"github.com/POMBNK/linktgBot/pkg/events/telegram"
	"github.com/POMBNK/linktgBot/pkg/storage/sqlite"
)

func main() {
	logger := logging.GetLogger("trace")
	cfg := config.GetCfg()
	storage, err := sqlite.New(cfg.DbPath)
	if err != nil {
		logger.Fatal("Не удалось подключиться к БД", err)
	}
	logger.Info("Инициализация БД...")
	storage.Init(context.Background())

	// Fetcher отправляет запросы, чтобы получать новые события
	// Processor обрабатывает события получаемые из Fetcher и отправляет сообщения в Client(отправляет ссылку или ошибку)
	// Consumer реализует старт программы получая на вход 2 объекта в виде Fetcher, Processor
	//Токен будет храниться в .yml/.env Потом через конфиг его подтянем
	eventsProcessor := telegram.New(logger, tgClient.New(logger, cfg.Host, cfg.Token), storage)
	logger.Info("Сервис запущен...")
	consumer := eventconsumer.New(logger, eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		logger.Fatal(err)
	}

}
