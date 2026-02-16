package config

type Config struct {
	EnableServer bool  `env:"ENABLE_SERVER" envDefault:"true"`
	Port         int32 `env:"PORT" envDefault:"8080"`
}
