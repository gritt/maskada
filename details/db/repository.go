package db

import (
	"database/sql"
	"time"

	// imports mysql db driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/gritt/maskada/core"
	"github.com/gritt/maskada/details"
)

// Repository is able to save and find a transaction(s).
type Repository struct {
	db *sqlx.DB
}

// NewRepository initialize the repository.
func NewRepository(cfg *details.Config) (repo *Repository, err error) {
	dns := cfg.DatabaseDNS()

	db, err := sqlx.Open("mysql", dns)
	if err != nil {
		return repo, errors.Wrap(err, "NewRepository failed")
	}

	return &Repository{db: db}, err
}

// Create persists a transaction in db.
func (r *Repository) Create(t core.Transaction) (core.Transaction, error) {
	if err := r.CreateCategory(t.Category); err != nil {
		return core.Transaction{}, err
	}

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

// CreateCategory persists a category in db.
func (r *Repository) CreateCategory(category core.Category) error {
	query := "INSERT IGNORE INTO `category` (`name`) VALUES (?)"

	_, err := r.db.Exec(query, category.Name)
	if err != nil {
		return errors.Wrap(err, "Repository.CreateCategory failed")
	}

	return nil
}

// Find transactions in db.
func (r *Repository) Find() ([]core.Transaction, error) {
	type row struct {
		ID       int            `db:"id"`
		Amount   int            `db:"amount"`
		Type     int            `db:"type"`
		Category string         `db:"category"`
		Date     time.Time      `db:"date"`
		Name     sql.NullString `db:"name"`
	}

	query := `SELECT 
				t.id "id", 
				t.amount "amount", 
				t.type "type", 
				t.category "category",
				t.date "date", 
				t.description "name"
				FROM transaction t
				ORDER by t.date`

	var rows []row
	if err := r.db.Select(&rows, query); err != nil {
		return []core.Transaction{}, errors.Wrap(err, "Repository.Find failed")
	}

	var trs []core.Transaction
	for _, row := range rows {
		trs = append(trs, core.Transaction{
			ID:       row.ID,
			Amount:   row.Amount,
			Type:     row.Type,
			Category: core.Category{Name: row.Category},
			Date:     row.Date,
			Name:     row.Name.String,
		})
	}

	return trs, nil
}
