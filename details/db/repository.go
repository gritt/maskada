package db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/gritt/maskada/core"
	"github.com/gritt/maskada/details"
)

type Repository struct {
	db *sqlx.DB
}

func (r *Repository) Create(t core.Transaction) (core.Transaction, error) {
	if t.Date.String() == "0001-01-01 00:00:00 +0000 UTC" {
		t.Date = time.Now().UTC()
	}

	query := "INSERT INTO `transaction` (`amount`, `type`, `category`, `description`, `date`) VALUES (?, ?, ?, ?, ?)"

	result, err := r.db.Exec(query, t.Amount, t.Type, t.Category.Name, t.Name, t.Date.UTC())
	if err != nil {
		return core.Transaction{}, errors.Wrap(err, "Repository.Create failed")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return core.Transaction{}, errors.Wrap(err, "Repository.Create failed")
	}
	t.ID = int(id)

	return t, nil
}

func NewRepository(cfg *details.Config) (repo *Repository, err error) {
	dns := cfg.DatabaseDNS()

	db, err := sqlx.Open("mysql", dns)
	if err != nil {
		return repo, errors.Wrap(err, "NewRepository failed")
	}

	return &Repository{db: db}, err
}
