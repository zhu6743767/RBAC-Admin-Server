package config

// SwaggerConfig Swagger文档配置
type SwaggerConfig struct {
	Enabled          bool   `yaml:"enabled"`
	Path             string `yaml:"path"`
	Title            string `yaml:"title"`
	Description      string `yaml:"description"`
	Version          string `yaml:"version"`
	TermsOfService   string `yaml:"terms_of_service"`
	ContactName      string `yaml:"contact_name"`
	ContactURL       string `yaml:"contact_url"`
	ContactEmail     string `yaml:"contact_email"`
	LicenseName      string `yaml:"license_name"`
	LicenseURL       string `yaml:"license_url"`
}