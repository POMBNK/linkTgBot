package sqlite

import (
	"context"
	"database/sql"
	"github.com/POMBNK/linktgBot/pkg/e"
	"github.com/POMBNK/linktgBot/pkg/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New конструктор стуктуры Storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("Не удалось открыть БД", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("Не удалось установить соединение с БД", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT,user_name TEXT)`

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return e.Wrap("Не удалось создать таблицу", err)
	}

	return nil
}

// Save сохраняет присланную статью в хранилище.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url,user_name) VALUES(?,?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.Wrap("Не удалось сохранить страницу", err)
	}
	return nil
}

// PickRandom выбирает случайную статью из хранилища.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.Wrap("Не удалось выбрать статью", err)
	}
	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove удаляет выбранную статью из хранилища.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.Wrap("Не удалось удалить статью", err)
	}
	return nil
}

func (s *Storage) IsExist(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count); err != nil {
		return false, e.Wrap("Такой записи не существует", err)
	}

	return count > 0, nil
}
