package config

type TracingConfig struct {
	ServiceName string
	JaegerURL   string
	Environment string
	Enabled     bool
}

func GetTracingConfig() *TracingConfig {
	return &TracingConfig{
		Enabled:     getEnv("TRACING_ENABLED", "false") == "true",
		ServiceName: getEnv("APP_NAME", "golang-starter-kit"),
		JaegerURL:   getEnv("JAEGER_URL", "http://localhost:14268/api/traces"),
		Environment: getEnv("APP_ENV", "development"),
	}
}
