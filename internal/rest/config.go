package rest

import (
	"os"
	"strconv"
)

const (
	defaultPort           = "8080"
	defaultStatsTimeoutMs = 3000
)

const (
	envVarPort           = "DMG_HTTP_PORT"
	envVarStatsTimeoutMs = "DMG_STATS_TIMEOUT_MS"
)

type config struct {
	port         string
	statsTimeout int
}

func (c *config) Read() {
	readSetting(envVarPort, defaultPort, &c.port)
	readNumericSetting(envVarStatsTimeoutMs, defaultStatsTimeoutMs, &c.statsTimeout)

	if c.statsTimeout <= 0 {
		c.statsTimeout = defaultStatsTimeoutMs
	}
}

func readSetting(setting, defaultValue string, result *string) {
	*result = os.Getenv(setting)
	if *result == "" {
		*result = defaultValue
	}
}

func readNumericSetting(setting string, defaultValue int, result *int) {
	val := os.Getenv(setting)

	if val != "" {
		valNum, err := strconv.Atoi(val)

		if err == nil {
			*result = valNum
			return
		}
	}

	*result = defaultValue
}
