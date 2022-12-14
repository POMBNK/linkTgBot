package storage

import (
	"context"
	"errors"
)

// Storage интерфейс для взаимодействия с хранилищем.
// Реализация структуры хранилища `Storage` должна удовлетворять данному интерфейсу `Storage`.
type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExist(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
	//Created time.Time
}
