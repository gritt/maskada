package db

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"

	"github.com/gritt/maskada/app"
	"github.com/gritt/maskada/core"
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

	tests := map[string]func(*testing.T, *app.Config){
		"should create with success with given date": func(t *testing.T, cfg *app.Config) {
			// arrange
			repo, err := NewRepository(cfg)
			assert.NoError(t, err)
			teardown := setupDBData(t, repo.db)
			defer teardown()

			given := core.Transaction{
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
				Date:     date,
			}

			// act
			got, gotErr := repo.Create(given)

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
		"should create with success with current date": func(t *testing.T, cfg *app.Config) {
			// arrange
			repo, err := NewRepository(cfg)
			assert.NoError(t, err)
			teardown := setupDBData(t, repo.db)
			defer teardown()

			given := core.Transaction{
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
			}

			// act
			got, gotErr := repo.Create(given)

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
			assert.Equal(t, want.Date.Day(), got.Date.Day())
			assert.Equal(t, want.Date.Month(), got.Date.Month())
			assert.Equal(t, want.Date.Year(), got.Date.Year())
			assert.Equal(t, want.Date.Hour(), got.Date.Hour())
			assert.Equal(t, want.Date.Minute(), got.Date.Minute())
		},
		"should create with success with name": func(t *testing.T, cfg *app.Config) {
			// arrange
			repo, err := NewRepository(cfg)
			assert.NoError(t, err)
			teardown := setupDBData(t, repo.db)
			defer teardown()

			given := core.Transaction{
				Amount:   amount,
				Type:     core.Credit,
				Category: core.Category{Name: "Food"},
				Date:     date,
				Name:     name,
			}

			// act
			got, gotErr := repo.Create(given)

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
			run(t, &cfg)
		})
	}
}

func mockDBConfig() (app.Config, error) {
	type MockConfig struct {
		Host     string `envconfig:"DATABASE_HOST" required:"true"`
		Port     string `envconfig:"DATABASE_PORT" required:"true"`
		Name     string `envconfig:"DATABASE_NAME" required:"true"`
		User     string `envconfig:"DATABASE_USERNAME" required:"true"`
		Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
	}

	var mockConfig MockConfig

	if err := envconfig.Process("", &mockConfig); err != nil {
		return app.Config{}, err
	}

	cfg := app.Config{}
	cfg.Database.Host = mockConfig.Host
	cfg.Database.Port = mockConfig.Port
	cfg.Database.Name = mockConfig.Name
	cfg.Database.User = mockConfig.User
	cfg.Database.Password = mockConfig.Password

	return cfg, nil
}

func setupDBData(t *testing.T, db *sqlx.DB) func() {
	// create db schema
	script, err := ioutil.ReadFile("../../app/db/migrations/schema.sql")
	if err != nil {
		t.Fatalf("setupDBData failed: %s", err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// setup data
	script, err = ioutil.ReadFile("../../app/db/test/setup.sql")
	if err != nil {
		t.Fatalf("setupDBData failed: %s", err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// teardown data
	return func() {
		script, err := ioutil.ReadFile("../../app/db/test/teardown.sql")
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
