package cfg

import (
	"OrcAI/internal/config"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func MustConfig(logger zerolog.Logger) config.Config {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	var c config.Config
	err = env.Parse(&c)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	return c
}
