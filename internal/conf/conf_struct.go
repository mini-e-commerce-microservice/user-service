package conf

type Config struct {
	AppPort       int                 `mapstructure:"APP_PORT"`
	AppMode       string              `mapstructure:"APP_MODE"`
	OpenTelemetry ConfigOpenTelemetry `mapstructure:"OPEN_TELEMETRY"`
	Minio         ConfigMinio         `mapstructure:"MINIO"`
	DatabaseDSN   string              `mapstructure:"DATABASE_DSN"`
	RabbitMQ      ConfigRabbitMQ      `mapstructure:"RABBIT_MQ"`
	Jwt           ConfigJWT           `mapstructure:"JWT"`
}

type ConfigRabbitMQ struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Url      string `mapstructure:"URL"`
}

type ConfigOpenTelemetry struct {
	Password   string `mapstructure:"PASSWORD"`
	Username   string `mapstructure:"USERNAME"`
	Endpoint   string `mapstructure:"ENDPOINT"`
	TracerName string `mapstructure:"TRACER_NAME"`
}

type ConfigMinio struct {
	Endpoint        string `mapstructure:"ENDPOINT"`
	AccessID        string `mapstructure:"ACCESS_ID"`
	SecretAccessKey string `mapstructure:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `mapstructure:"USE_SSL"`
	PrivateBucket   string `mapstructure:"PRIVATE_BUCKET"`
}

type ConfigJWT struct {
	OtpUsecase string `mapstructure:"OTP_USECASE"`
}
