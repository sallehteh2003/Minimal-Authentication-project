package Config

type Config struct {
	Database struct {
		Host     string `env:"DATABASE_HOST" env-default:"localhost"`
		Port     int    `env:"DATABASE_PORT" env-default:"5432"`
		Name     string `env:"DATABASE_NAME" env-default:"minimal-auth"`
		Username string `env:"DATABASE_USERNAME" env-default:"postgres"`
		Password string `env:"DATABASE_PASSWORD" env-default:"saleh"`
	}
	Server struct {
		Port int `json:"port" env-default:"8080"`
	}
}
