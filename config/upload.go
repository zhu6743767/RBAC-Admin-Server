package config

// UploadConfig 上传配置
type UploadConfig struct {
	MaxFileSize    int      `yaml:"max_file_size"`
	AllowedTypes   []string `yaml:"allowed_types"`
	SavePath       string   `yaml:"save_path"`
	UseHashName    bool     `yaml:"use_hash_name"`
	BackupEnabled  bool     `yaml:"backup_enabled"`
	BackupPath     string   `yaml:"backup_path"`
}