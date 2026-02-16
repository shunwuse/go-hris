package cache

type Config struct {
	RedisAddr     string `koanf:"redis_addr"`
	RedisPassword string `koanf:"redis_password"`
	RedisDB       int    `koanf:"redis_db"`
	UseMiniredis  bool   `koanf:"use_miniredis"`
}
