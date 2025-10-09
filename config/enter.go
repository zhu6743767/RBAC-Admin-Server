package config

// Config 全局配置入口
type Config struct {
	System       SystemConfig     `yaml:"system"`
	DB           DBConfig         `yaml:"db"`
	Redis        RedisConfig      `yaml:"redis"`
	JWT          JWTConfig        `yaml:"jwt"`
	Log          LogConfig        `yaml:"log"`
	Security     SecurityConfig   `yaml:"security"`
	CORS         CORSConfig       `yaml:"cors"`
	Monitoring   MonitoringConfig `yaml:"monitoring"`
	Swagger      SwaggerConfig    `yaml:"swagger"`
	Performance  PerformanceConfig `yaml:"performance"`
	App          AppConfig        `yaml:"app"`
	Upload       UploadConfig     `yaml:"upload"`
	Email        Email            `yaml:"email"`
	Captcha      Captcha          `yaml:"captcha"`
}