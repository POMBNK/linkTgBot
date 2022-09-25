package eventconsumer

import (
	"context"
	"time"

	"github.com/POMBNK/linktgBot/logging"
	"github.com/POMBNK/linktgBot/pkg/events"
)

type Consumer struct {
	logger    *logging.Logger
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(logger *logging.Logger, fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		logger:    logger,
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(context.Background(), c.batchSize)
		if err != nil {
			c.logger.Error("consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.handleEvents(context.Background(), gotEvents); err != nil {
			c.logger.Info(err)
			continue
		}

	}

}

/*Возможные проблемы:
1.Потеря данных при нестабильном соединении. Решения: ретрай, вернуть обратно в сторадж,
  фолбэк(сохранение в ОЗУ,например), подтверждение для Fetcher (Fecther не делает сдвиг,
  пока не обработал всю пачку данных или мы сами передаем ему сколько обработать каждый раз)
2.Обработка всей пачки при проблеме. Т.е мы раз за разом пытаемся обработать данные, хотя уже N раз получили ошибку
  Решения: Останавливаться после первой ошибки, либо вести счетчик и после скольки-то ошибок останавливать обработку,
  Последнее решение sync.WaitGroup()
*/

func (c Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {

		if err := c.processor.Process(ctx, event); err != nil {
			c.logger.Error("Processor: Не удалось обработать событие: %s", err.Error())
			continue
		}
	}
	return nil
}
