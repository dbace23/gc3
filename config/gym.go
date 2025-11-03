package config

type Config struct {
	AppPort     string `env:"APP_PORT,required"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	JWTSecret   string `env:"JWT_SECRET,required"`
	RapidKey    string `env:"RAPIDAPI_KEY,required"`
	RapidHost   string `env:"RAPIDAPI_HOST,required"`
	SwaggerOn   bool   `env:"SWAGGER_ENABLED" envDefault:"true"`
}
