package config

type Config struct {
	NATS NATS
}

type NATS struct {
	URL string `env:"NATS_URL,required"`
}
