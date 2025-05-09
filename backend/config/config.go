package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	DBConnection string
	JWTSecret    string
	Port         string
	IsDev        bool
	TokenExpiry  time.Duration
}

func Load() *Config {
	return &Config{
		DBConnection: getEnv("DB_CONNECTION", "postgres://postgres:postgrespassword@postgres:5432/postgres"),
		JWTSecret:    getEnv("JWT_SECRET", "9a4049355a147359591e126f535fe1146eb3de0bb00c21532b79032efdc39b6c84d260618b723126afc587d8d5087ba0a9b91817b48fb72f45effbbd86a2f2186dbd6948c06be8b579e440a8ed8f19e59b48c4f785f49d2f6ff86f165d052fea0ea2c416cba31da1dec6c1ca547ec2606378bbb666dc5d82e921d458af823acfdb36c3fecbb7ba14add576ad88eefba0b31f3606f8fc92fd52cf144a3f4891b5951822e49a0f966799fadab1bed093f0067d3c33f2660b42fcb41cb6b3e5e2da71ef613e8f117454def63f20e2e8e975677fa10bc75cd76626696ff6614ea1d8c7c866592e4118d439b3180bc1f0546aec310efd5e8ed2d18d994502d9d5c7ed"),
		Port:         getEnv("PORT", "3001"),
		IsDev:        os.Getenv("HASURA_GRAPHQL_DEV_MODE") == "true",
		TokenExpiry:  24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if fallback == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return fallback
}