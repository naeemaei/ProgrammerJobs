package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Host:     "localhost",
			Port:     1433,
			Username: "sa",
			Password: "123",
		},
	}
}
