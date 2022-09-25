package telegram

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/POMBNK/linktgBot/pkg/storage"
)

const (
	RndCmd         = "/rnd"
	HelpCmd        = "/help"
	StartCmd       = "/start"
	UnknownCommand = "Извини, я не знаю такой команды"
)

func (d *Dispatcher) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	d.logger.Infof("got new command: `%s` from %s", text, username)

	if isSavingUrl(text) {
		return d.savePage(ctx, chatID, text, username)
	}

	switch text {
	case RndCmd:
		d.sendRandom(ctx, chatID, username)
	case HelpCmd:
		d.help(chatID)
	case StartCmd:
		d.hello(chatID)
	default:
		d.tg.SendMessage(chatID, UnknownCommand)
	}
	return nil
}

func isSavingUrl(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

func (d *Dispatcher) savePage(ctx context.Context, chatID int, pageURL string, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := d.storage.IsExist(ctx, page)
	if err != nil {
		return err
	}

	if isExist {
		return d.tg.SendMessage(chatID, AlreadySavedMsg)
	}

	if err := d.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := d.tg.SendMessage(chatID, SuccessMsg); err != nil {
		return err
	}

	return nil
}

func (d *Dispatcher) sendRandom(ctx context.Context, chatID int, username string) error {
	page, err := d.storage.PickRandom(ctx, username)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return d.tg.SendMessage(chatID, "Не удалось взять случайную статью из списка, так как он пуст")
	}

	if err = d.tg.SendMessage(chatID, page.URL); err != nil {
		d.logger.Error("Не удалось отправить случайную статью")
		return err
	}
	if err = d.storage.Remove(ctx, page); err != nil {
		return err
	}
	return nil
}

func (d *Dispatcher) help(chatID int) error {
	if err := d.tg.SendMessage(chatID, HelpMsg); err != nil {
		return err
	}
	return nil
}

func (d *Dispatcher) hello(chatID int) error {
	if err := d.tg.SendMessage(chatID, HelloMsg); err != nil {
		return err
	}
	return nil
}
