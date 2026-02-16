package database

type Config struct {
	Host     string `koanf:"db_host"`
	Port     string `koanf:"db_port"`
	User     string `koanf:"db_user"`
	Password string `koanf:"db_password"`
	Name     string `koanf:"db_name"`
	SSLMode  string `koanf:"db_ssl_mode"`
}
