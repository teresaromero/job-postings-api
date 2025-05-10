package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	SortRulesOrder []string
	JwtSecret      []byte
}

func Load() *Config {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	sortRulesEnvStr := os.Getenv("SORT_RULES_ORDER")
	var sortRulesEnv []string
	if sortRulesEnvStr != "" {
		sortRulesEnv = strings.Split(sortRulesEnvStr, ",")
	}

	return &Config{
		SortRulesOrder: sortRulesEnv,
		JwtSecret:      []byte(jwtSecret),
	}
}
