package config

import (
	"os"
	"strings"
)

type Config struct {
	SortRulesOrder []string
}

func Load() *Config {
	sortRulesEnvStr := os.Getenv("SORT_RULES_ORDER")
	var sortRulesEnv []string
	if sortRulesEnvStr != "" {
		sortRulesEnv = strings.Split(sortRulesEnvStr, ",")
	}

	return &Config{
		SortRulesOrder: sortRulesEnv,
	}
}
