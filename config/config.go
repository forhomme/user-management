package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	SecretKey            string
	FilePath             string
	DatabaseUri          string
	Database             string
	CourseCollection     string
	CategoryCollection   string
	AuditTrailCollection string
	BucketName           string
	AuthExpire           time.Duration
	RefreshTokenExpire   time.Duration
}

// Load config file from given path
func LoadLocalConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
