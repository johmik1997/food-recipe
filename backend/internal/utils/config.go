package utils

import (
	"github.com/caarlos0/env/v8"
)

type Config struct {
    JWTSecret         string `env:"JWT_SECRET,required"`
    HasuraEndpoint    string `env:"HASURA_GRAPHQL_ENDPOINT,required"`
    HasuraAdminSecret string `env:"HASURA_GRAPHQL_ADMIN_SECRET,required"`
    HTTPPort          string `env:"HTTP_PORT" envDefault:"3001"`
    CloudinaryCloudName  string `env:"CLOUDINARY_CLOUD_NAME,required"`
    CloudinaryAPIKey        string `env:"CLOUDINARY_API_KEY,required"`
    CloudinaryAPISecret  string `env:"CLOUDINARY_API_SECRET,required"`
}

func LoadConfig() (*Config, error) {
    cfg := new(Config)
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}