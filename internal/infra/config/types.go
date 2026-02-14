package config

type Config struct {
	Service  ServiceConfig  `koanf:",squash"`
	Log      LogConfig      `koanf:",squash"`
	Database DatabaseConfig `koanf:",squash"`
	Auth     AuthConfig     `koanf:",squash"`
	Cache    CacheConfig    `koanf:",squash"`
}

type ServiceConfig struct {
	Environment              string `koanf:"env"`
	Port                     string `koanf:"port"`
	IdempotencyExpireMinutes int    `koanf:"idempotency_expire_minutes"`
	ProfilerToken            string `koanf:"profiler_token"` // token for pprof access
}

type LogConfig struct {
	Level      string `koanf:"log_level"`       // debug, info, warn, error
	FilePath   string `koanf:"log_file_path"`   // save destination
	MaxSize    int    `koanf:"log_max_size"`    // megabytes
	MaxBackups int    `koanf:"log_max_backups"` // number of backups
	MaxAge     int    `koanf:"log_max_age"`     // days
	Compress   bool   `koanf:"log_compress"`    // compress old files
}

type DatabaseConfig struct {
	Host     string `koanf:"db_host"`
	Port     string `koanf:"db_port"`
	User     string `koanf:"db_user"`
	Password string `koanf:"db_password"`
	Name     string `koanf:"db_name"`
	SSLMode  string `koanf:"db_ssl_mode"`
}

type AuthConfig struct {
	JWTSecret               string `koanf:"jwt_secret"`
	JWTExpireMinutes        int    `koanf:"jwt_expire_minutes"`
	JWTRefreshExpireMinutes int    `koanf:"jwt_refresh_expire_minutes"`
}

type CacheConfig struct {
	RedisAddr     string `koanf:"redis_addr"`
	RedisPassword string `koanf:"redis_password"`
	RedisDB       int    `koanf:"redis_db"`

	UseMiniredis bool `koanf:"use_miniredis"`
}
