package details

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gritt/maskada/test"
)

func TestNewConfig(t *testing.T) {
	variables := getEnvironmentVariables()

	os.Clearenv()
	for wantVariable, wantValue := range variables {
		if err := os.Setenv(wantVariable, wantValue); err != nil {
			t.Fatalf("failed to: Setenv %s with value %s", wantVariable, wantValue)
		}
	}

	gotCfg, gorErr := NewConfig()
	assert.NoError(t, gorErr)

	assert.Equal(t, variables["DATABASE_HOST"], gotCfg.Database.Host)
	assert.Equal(t, variables["DATABASE_PORT"], gotCfg.Database.Port)
	assert.Equal(t, variables["DATABASE_NAME"], gotCfg.Database.Name)
	assert.Equal(t, variables["DATABASE_USERNAME"], gotCfg.Database.User)
	assert.Equal(t, variables["DATABASE_PASSWORD"], gotCfg.Database.Password)
}

func TestNewConfig_with_missing_environment_variables(t *testing.T) {
	// arrange
	variables := getEnvironmentVariables()

	for wantEnv, _ := range variables {
		os.Clearenv()

		for env, value := range variables {
			if env == wantEnv {
				continue
			}

			if err := os.Setenv(env, value); err != nil {
				t.Fatalf("failed to: Setenv %s with value %s", env, value)
			}
		}

		// act
		_, gotErr := NewConfig()

		// assert
		assert.EqualError(t, gotErr, fmt.Sprintf("required key %s missing value", wantEnv))
	}
}

func TestConfig_DatabaseDNS(t *testing.T) {
	// arrange
	variables := getEnvironmentVariables()

	os.Clearenv()
	for wantVariable, wantValue := range variables {
		if err := os.Setenv(wantVariable, wantValue); err != nil {
			t.Fatalf("failed to: Setenv %s with value %s", wantVariable, wantValue)
		}
	}

	gotCfg, _ := NewConfig()

	// act
	gotDNS := gotCfg.DatabaseDNS()

	wantDNS := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		variables["DATABASE_USERNAME"],
		variables["DATABASE_PASSWORD"],
		variables["DATABASE_HOST"],
		variables["DATABASE_PORT"],
		variables["DATABASE_NAME"],
	)

	// assert
	assert.Equal(t, wantDNS, gotDNS)
}

func getEnvironmentVariables() map[string]string {
	return map[string]string{
		"DATABASE_HOST":     test.RandomDomain(),
		"DATABASE_PORT":     strconv.Itoa(test.RandomNumber()),
		"DATABASE_NAME":     test.RandomUsername(),
		"DATABASE_USERNAME": test.RandomUsername(),
		"DATABASE_PASSWORD": test.RandomPassword(),
	}
}
