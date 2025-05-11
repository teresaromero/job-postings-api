package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	SortRulesOrder []string
	JwtSecret      []byte

	SeedData bool
	SeedFile string
}

func Load() *Config {

	seedDataStr := os.Getenv("SEED_DATA")
	seedData, err := strconv.ParseBool(seedDataStr)
	if err != nil {
		seedData = false
	}
	seedFile := os.Getenv("SEED_FILE")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	sortRulesEnvStr := os.Getenv("SORT_RULES_ORDER")
	var sortRulesEnv []string
	if sortRulesEnvStr != "" {
		sortRulesEnv = strings.Split(sortRulesEnvStr, ",")
	}

	ginMode := os.Getenv("GIN_MODE")
	// Only for testing purposes
	customNowDate := os.Getenv("NOW")
	if ginMode == "release" && customNowDate != "" {
		log.Fatalf("Custom NOW date is not allowed in release mode")
	}
	if customNowDate != "" {
		log.Default().Printf("Custom NOW date set to: %s", customNowDate)
		nowDate, err := time.Parse(time.RFC3339, customNowDate)
		if err != nil {
			log.Fatalf("Invalid NOW date format: %v", err)
		}
		Now = func() time.Time {
			return nowDate
		}
	}

	return &Config{
		SortRulesOrder: sortRulesEnv,
		JwtSecret:      []byte(jwtSecret),
		SeedData:       seedData || seedFile != "",
		SeedFile:       seedFile,
	}
}
