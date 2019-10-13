package details

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the app configuration (eg: ENV, Database, Network).
type Config struct {
	Database struct {
		Host     string `envconfig:"DATABASE_HOST" required:"true"`
		Port     string `envconfig:"DATABASE_PORT" required:"true"`
		Name     string `envconfig:"DATABASE_NAME" required:"true"`
		User     string `envconfig:"DATABASE_USERNAME" required:"true"`
		Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
	}
}

// NewConfig initialize the config.
func NewConfig() (*Config, error) {
	c := Config{}
	if err := envconfig.Process("", &c); err != nil {
		return &c, err
	}
	return &c, nil
}

// DatabaseDNS builds the DB data source name.
func (c *Config) DatabaseDNS() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}
