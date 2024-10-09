package conf

type AppConfig struct {
	AppPort     int    `json:"APP_PORT"`
	AppMode     string `json:"APP_MODE"`
	DatabaseDSN string `json:"DATABASE_DSN"`
	ServiceName string `json:"SERVICE_NAME"`
}
