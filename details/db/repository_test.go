package db

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"

	"github.com/gritt/maskada/core"
	"github.com/gritt/maskada/details"
	"github.com/gritt/maskada/test"
)

func TestNewRepository(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// arrange
	cfg, err := mockDBConfig()
	if err != nil {
		t.Fatalf("failed to mockDBConfig: %s", err)
	}

	// act
	gotRepo, gotErr := NewRepository(&cfg)

	// assert
	assert.NoError(t, gotErr)
	assert.IsType(t, &sqlx.DB{}, gotRepo.db)
}

func TestRepository_Create(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	amount := test.RandomNumber()
	name := test.RandomName()
	date := time.Now().UTC().Truncate(time.Hour * 24).Add(-time.Hour * 24)

	cfg, err := mockDBConfig()
	if err != nil {
		t.Fatalf("mockDBConfig failed: %s", err)
	}

	tests := map[string]func(*testing.T, *Repository){
		"when connection is down": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			teardown()

			// act
			_, gotErr := r.Create(core.Transaction{})

			// assert
			assert.EqualError(t, gotErr, "Repository.CreateCategory failed: sql: database is closed")
		},
		"when a date is given": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			defer teardown()

			given := core.Transaction{
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
				Date:     date,
			}

			// act
			got, gotErr := r.Create(given)

			want := core.Transaction{
				ID:       7,
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
				Date:     date,
			}

			// assert
			assert.NoError(t, gotErr)
			assert.Equal(t, want, got)
		},
		"when no date is given, use current time": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			defer teardown()

			given := core.Transaction{
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
			}

			// act
			got, gotErr := r.Create(given)

			want := core.Transaction{
				ID:       7,
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
				Date:     time.Now().UTC(),
			}

			// assert
			assert.NoError(t, gotErr)
			assert.Equal(t, want.ID, got.ID)
			assert.Equal(t, want.Amount, got.Amount)
			assert.Equal(t, want.Type, got.Type)
			assert.Equal(t, want.Category, got.Category)
			assert.Equal(t, want.Date.Format(time.RFC3339), got.Date.Format(time.RFC3339))
		},
		"when a name is given": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			defer teardown()

			given := core.Transaction{
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
				Date:     date,
				Name:     name,
			}

			// act
			got, gotErr := r.Create(given)

			want := core.Transaction{
				ID:       7,
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
				Date:     date,
				Name:     name,
			}

			// assert
			assert.NoError(t, gotErr)
			assert.Equal(t, want, got)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := NewRepository(&cfg)
			assert.NoError(t, err)

			run(t, r)
		})
	}
}

func TestRepository_CreateCategory(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	type row struct {
		Name string `db:"name"`
	}

	cfg, err := mockDBConfig()
	if err != nil {
		t.Fatalf("mockDBConfig failed: %s", err)
	}

	tests := map[string]func(*testing.T, *Repository){
		"when connection is down": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			teardown()

			// act
			_, gotErr := r.Create(core.Transaction{})

			// assert
			assert.EqualError(t, gotErr, "Repository.CreateCategory failed: sql: database is closed")
		},
		"when category does not exists": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			defer teardown()

			testCategory := test.RandomName()
			given := core.Category{Name: testCategory}

			// act
			gotErr := r.CreateCategory(given)

			// assert
			assert.NoError(t, gotErr)

			var rows []row
			if err := r.db.Select(&rows, `SELECT * FROM category WHERE name = (?)`, testCategory); err != nil {
				t.Fail()
			}
			assert.Equal(t, testCategory, rows[0].Name)
		},
		"when category exists": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			defer teardown()

			testCategory := "Entertainment"
			given := core.Category{Name: testCategory}

			// act
			gotErr := r.CreateCategory(given)

			// assert
			assert.NoError(t, gotErr)

			var rows []row
			if err := r.db.Select(&rows, `SELECT * FROM category WHERE name = (?)`, testCategory); err != nil {
				t.Fail()
			}
			assert.Len(t, rows, 1)
			assert.Equal(t, testCategory, rows[0].Name)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := NewRepository(&cfg)
			assert.NoError(t, err)

			run(t, r)
		})
	}
}

func TestRepository_Find(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	transactions := []core.Transaction{
		{ID: 1, Amount: 99, Type: core.Credit, Category: core.Category{Name: "Entertainment"}, Date: time.Now().UTC()},
		{ID: 2, Amount: 11, Type: core.Credit, Category: core.Category{Name: "Food"}, Date: time.Now().UTC()},
		{ID: 3, Amount: 32, Type: core.Credit, Category: core.Category{Name: "Food"}, Date: time.Now().UTC()},
		{ID: 4, Amount: 5300, Type: core.Income, Category: core.Category{Name: "Work"}, Date: time.Now().UTC()},
		{ID: 5, Amount: 129, Type: core.Debit, Category: core.Category{Name: "Home"}, Date: time.Now().UTC(), Name: "Internet"},
		{ID: 6, Amount: 129, Type: core.Debit, Category: core.Category{Name: "Home"}, Date: time.Now().UTC(), Name: "Electricity"},
	}

	cfg, err := mockDBConfig()
	if err != nil {
		t.Fatalf("mockDBConfig failed: %s", err)
	}

	tests := map[string]func(t *testing.T, r *Repository){
		"when connection is down": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			teardown()

			// act
			_, gotErr := r.Find()

			// assert
			assert.EqualError(t, gotErr, "Repository.Find failed: sql: database is closed")
		},
		"when transactions are found": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			defer teardown()

			// act
			got, gotErr := r.Find()

			// assert
			assert.NoError(t, gotErr)
			for want, gotTrs := range got {
				assert.Equal(t, transactions[want].ID, gotTrs.ID)
				assert.Equal(t, transactions[want].Amount, gotTrs.Amount)
				assert.Equal(t, transactions[want].Type, gotTrs.Type)
				assert.Equal(t, transactions[want].Category, gotTrs.Category)
				assert.Equal(t, transactions[want].Date.Format(time.RFC822), gotTrs.Date.Format(time.RFC822))
				assert.Equal(t, transactions[want].Name, gotTrs.Name)
			}
		},
		"when no transactions are found": func(t *testing.T, r *Repository) {
			// arrange
			teardown := setupDBData(t, r.db)
			defer teardown()

			_, err := r.db.Exec(`DELETE from transaction`)
			if err != nil {
				t.Fatalf("when no transactions are found failed: %s", err)
			}

			// act
			got, gotErr := r.Find()

			// assert
			assert.Empty(t, got)
			assert.NoError(t, gotErr)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := NewRepository(&cfg)
			assert.NoError(t, err)

			run(t, r)
		})
	}
}

func mockDBConfig() (details.Config, error) {
	type MockConfig struct {
		Host     string `envconfig:"DATABASE_HOST" required:"true"`
		Port     string `envconfig:"DATABASE_PORT" required:"true"`
		Name     string `envconfig:"DATABASE_NAME" required:"true"`
		User     string `envconfig:"DATABASE_USERNAME" required:"true"`
		Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
	}

	var mockConfig MockConfig

	if err := envconfig.Process("", &mockConfig); err != nil {
		return details.Config{}, err
	}

	cfg := details.Config{}
	cfg.Database.Host = mockConfig.Host
	cfg.Database.Port = mockConfig.Port
	cfg.Database.Name = mockConfig.Name
	cfg.Database.User = mockConfig.User
	cfg.Database.Password = mockConfig.Password

	return cfg, nil
}

func setupDBData(t *testing.T, db *sqlx.DB) func() {
	// create db schema
	script, err := ioutil.ReadFile("../../details/db/migrations/schema.sql")
	if err != nil {
		t.Fatalf("setupDBData failed: %s", err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// setup data
	script, err = ioutil.ReadFile("../../details/db/test/setup.sql")
	if err != nil {
		t.Fatalf("setupDBData failed: %s", err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// teardown data
	return func() {
		script, err := ioutil.ReadFile("../../details/db/test/teardown.sql")
		if err != nil {
			t.Fatalf("setupDBData failed: %s", err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatalf("setupDBData failed: %s", err)
		}

		db.Close()
	}
}
