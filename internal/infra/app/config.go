package app

type ServiceConfig struct {
	Environment              string `koanf:"env"`
	Port                     string `koanf:"port"`
	IdempotencyExpireMinutes int    `koanf:"idempotency_expire_minutes"`
	ProfilerToken            string `koanf:"profiler_token"`
}

type AuthConfig struct {
	JWTSecret               string `koanf:"jwt_secret"`
	JWTExpireMinutes        int    `koanf:"jwt_expire_minutes"`
	JWTRefreshExpireMinutes int    `koanf:"jwt_refresh_expire_minutes"`
}
