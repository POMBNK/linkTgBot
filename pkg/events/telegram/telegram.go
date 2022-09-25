package telegram

import (
	"context"
	"errors"

	"github.com/POMBNK/linktgBot/logging"
	"github.com/POMBNK/linktgBot/pkg/clients/telegram"
	"github.com/POMBNK/linktgBot/pkg/e"
	"github.com/POMBNK/linktgBot/pkg/events"
	"github.com/POMBNK/linktgBot/pkg/storage"
)

// Processor struct.
type Dispatcher struct {
	logger  *logging.Logger
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

// New конструктор для стуктуры Dispatcher клиента telegram.
func New(logger *logging.Logger, tg *telegram.Client, storage storage.Storage) *Dispatcher {
	return &Dispatcher{
		logger:  logger,
		tg:      tg,
		storage: storage,
	}
}

type Meta struct {
	ChatID   int
	Username string
}

func (d *Dispatcher) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	//Получаем обновления
	updates, err := d.tg.Updates(d.offset, limit)
	if err != nil {
		d.logger.Error("Не удалось получить обновления")
		return nil, err
	}
	// Если обновлений нет, сразу возвращаем nil
	if len(updates) == 0 {
		return nil, nil
	}
	// Создаем переменную и алоцируем память для результата заранее
	// так как заранее известна величина updates
	res := make([]events.Event, 0, len(updates))

	// Преобразование всех []updates в новый тип []event
	// Таким образом мы создавая новую сущность events можно быстро изменить бота под другой месенджер
	// Так как updates структура с уникальными полями и сильно связана с телеграмом.
	// Events не имеет данного минуса, чтобы не пришло на вход он преобразует в свой тип.
	for _, u := range updates {
		res = append(res, event(u))
	}

	// Обноваляем параметр смещения, чтобы получать новые обновления с последнего, а не с первого
	d.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

// event преобразует сущность update в сущность event
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}

	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message

}

// Process
func (d *Dispatcher) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return d.processMessage(ctx, event)
	default:
		e.Wrap("Не удалось обработать сообщение", errors.New("неизвестный eventType"))
	}
	return nil
}

// processMessage обрабатывает полученное сообщение
func (d *Dispatcher) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return err
	}
	if err := d.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		d.logger.Error("Не удалось обработать сообщение")
		return err
	}
	return nil
}

// meta производит приведение типа пустого интерфейса meta к структуре Meta
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("не удалось выполнить приведение типа Meta", errors.New("неизвестный eventType"))
	}
	return res, nil
}
