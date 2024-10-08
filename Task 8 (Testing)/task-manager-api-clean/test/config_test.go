package tests

import (
    "os"
    "testing"
    "task-manager-api-clean/config"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"

)

type ConfigTestSuite struct {
    suite.Suite
}

func (suite *ConfigTestSuite) SetupSuite() {
    os.Setenv("DATABASE_URL", "mongodb://localhost:27017")
    os.Setenv("JWT_SECRET", "mysecret")
    os.Setenv("JWT_EXPIRATION", "3600")
    os.Setenv("PORT", "8080")
    os.Setenv("TIMEOUT", "30s")
    os.Setenv("DATABASE_NAME", "testdb")

}

func (suite *ConfigTestSuite) TestLoad_Success() {
    env, err := config.Load()
    require.NoError(suite.T(), err)
    assert.Equal(suite.T(), "mongodb://localhost:27017", env.DatabaseURL)
	assert.Equal(suite.T(), "mysecret", env.JwtSecret)
    assert.Equal(suite.T(), 3600, env.JwtExpiration)
    assert.Equal(suite.T(), "8080", env.Port)
    assert.Equal(suite.T(), "30s", env.TimeOut)
    assert.Equal(suite.T(), "testdb", env.DatabaseName)
}

func TestConfigTestSuite(t *testing.T) {
    suite.Run(t, new(ConfigTestSuite))
}