package data

import "os"

const (
	defaultHost     = "localhost"
	defaultPort     = "5432"
	defaultDatabase = "demography"
	defaultUser     = "postgres"
	defaultPassword = "postgres"
)

const (
	envVarHost     = "DMG_STORAGE_HOST"
	envVarPort     = "DMG_STORAGE_PORT"
	envVarDatabase = "DMG_STORAGE_DATABASE"
	envVarUser     = "DMG_STORAGE_USER"
	envVarPassword = "DMG_STORAGE_PASSWORD"
)

type config struct {
	host     string
	port     string
	database string
	user     string
	password string
}

func (c *config) Read() {
	readSetting(envVarHost, defaultHost, &c.host)
	readSetting(envVarPort, defaultPort, &c.port)
	readSetting(envVarDatabase, defaultDatabase, &c.database)
	readSetting(envVarUser, defaultUser, &c.user)
	readSetting(envVarPassword, defaultPassword, &c.password)
}

func readSetting(setting, defaultValue string, result *string) {
	*result = os.Getenv(setting)
	if *result == "" {
		*result = defaultValue
	}
}
