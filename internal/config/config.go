package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type ProviderType string

const (
	ProviderAWS_S3  ProviderType = "AWS_S3"
	ProviderCF_R2   ProviderType = "CF_R2"
	ProviderGCP_GCS ProviderType = "GCP_GCS"
)

type ProviderConfig struct {
	Name          string       `yaml:"name"` // logical name, used in routes
	Type          ProviderType `yaml:"type"`
	Region        string       `yaml:"region,omitempty"`
	Endpoint      string       `yaml:"endpoint,omitempty"`
	AccessKey     string       `yaml:"access_key,omitempty"`
	SecretKey     string       `yaml:"secret_key,omitempty"`
	UseSSL        bool         `yaml:"use_ssl,omitempty"`
	AccountID     string       `yaml:"account_id,omitempty"`     // for R2
	ProjectID     string       `yaml:"project,omitempty"`        // for GCS
	CredentialRef string       `yaml:"credential_ref,omitempty"` // e.g. path to JSON
}

type RouteRule struct {
	Env            string `yaml:"env"`
	LogicalRegion  string `yaml:"logical_region"`
	Bucket         string `yaml:"bucket"`
	StorageClass   string `yaml:"storage_class"` // HOT/COLD/ARCHIVE
	ProviderName   string `yaml:"provider_name"` // reference to Providers[*].Name
	ProviderBucket string `yaml:"provider_bucket"`
}

type ObjectStorageConfig struct {
	DefaultStorageClass string           `yaml:"default_storage_class"`
	Routes              []RouteRule      `yaml:"routes"`
	Providers           []ProviderConfig `yaml:"providers"`
}

type RedisConfig struct {
	Addr         string        `yaml:"addr"`
	Password     string        `yaml:"password,omitempty"`
	DB           int           `yaml:"db"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"` // ":8080"
}

type DBConfig struct {
	Driver string `yaml:"driver"` // "pgx"
	DSN    string `yaml:"dsn"`    // connection string
}

type Config struct {
	Env           string              `yaml:"env"`
	HTTP          HTTPConfig          `yaml:"http"`
	Redis         RedisConfig         `yaml:"redis"`
	DB            DBConfig            `yaml:"db"`
	ObjectStorage ObjectStorageConfig `yaml:"object_storage"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.ObjectStorage.DefaultStorageClass == "" {
		cfg.ObjectStorage.DefaultStorageClass = "HOT"
	}
	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = ":8080"
	}
	return &cfg, nil
}
