package config

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled            bool   `yaml:"enabled"`
	PrometheusPort     int    `yaml:"prometheus_port"`
	HealthCheckPath    string `yaml:"health_check_path"`
	MetricsPath        string `yaml:"metrics_path"`
	TraceSamplingRate  float64 `yaml:"trace_sampling_rate"`
}