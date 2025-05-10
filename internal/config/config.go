package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	SortRulesOrder []string
	JwtSecret      []byte

	SeedData bool
}

func Load() *Config {

	seedDataStr := os.Getenv("SEED_DATA")
	seedData, err := strconv.ParseBool(seedDataStr)
	if err != nil {
		log.Printf("SEED_DATA environment variable is not set, defaulting to false")
		seedData = false
	}

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
		SeedData:       seedData,
	}
}
