package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Security SecurityConfig `mapstructure:"security"`
	Modules  ModuleConfig   `mapstructure:"modules"`
}

type ServerConfig struct {
	Addr          string `mapstructure:"addr"`
	PublicBaseURL string `mapstructure:"public_base_url"`
	MaxBodyBytes  int64  `mapstructure:"max_body_bytes"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type StorageConfig struct {
	Driver   string `mapstructure:"driver"`
	LocalDir string `mapstructure:"local_dir"`
}

type SecurityConfig struct {
	AdminToken           string `mapstructure:"admin_token"`
	ContentEncryptionKey string `mapstructure:"content_encryption_key"`
}

type ModuleConfig struct {
	ShortLink    TTLConfig     `mapstructure:"short_link"`
	ImageHosting AssetConfig   `mapstructure:"image_hosting"`
	Clipboard    ClipboardConf `mapstructure:"clipboard"`
	FileStash    AssetConfig   `mapstructure:"file_stash"`
}

type TTLConfig struct {
	DefaultTTL      time.Duration     `mapstructure:"default_ttl"`
	AllowCustomSlug bool              `mapstructure:"allow_custom_slug"`
	DomainMappings  map[string]string `mapstructure:"domain_mappings"`
}

type AssetConfig struct {
	DefaultTTL time.Duration `mapstructure:"default_ttl"`
	MaxBytes   int64         `mapstructure:"max_bytes"`
}

type ClipboardConf struct {
	DefaultTTL time.Duration `mapstructure:"default_ttl"`
	MaxVisits  int           `mapstructure:"max_visits"`
}

func Load(path string) (Config, error) {
	cfg := defaults()
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType("toml")
	v.SetEnvPrefix("COMICAL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return cfg, err
		}
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func defaults() Config {
	return Config{
		Server: ServerConfig{
			Addr:          "127.0.0.1:8080",
			PublicBaseURL: "http://localhost:8080",
			MaxBodyBytes:  104857600,
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    "file:comical.db?_foreign_keys=on",
		},
		Storage: StorageConfig{Driver: "local", LocalDir: "./data/objects"},
		Security: SecurityConfig{
			AdminToken:           "change-me",
			ContentEncryptionKey: "change-me-32-bytes",
		},
		Modules: ModuleConfig{
			ShortLink:    TTLConfig{DefaultTTL: 168 * time.Hour, AllowCustomSlug: true, DomainMappings: map[string]string{}},
			ImageHosting: AssetConfig{DefaultTTL: 720 * time.Hour, MaxBytes: 10485760},
			Clipboard:    ClipboardConf{DefaultTTL: time.Hour, MaxVisits: 5},
			FileStash:    AssetConfig{DefaultTTL: 168 * time.Hour, MaxBytes: 104857600},
		},
	}
}
