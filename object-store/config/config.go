package config

import (
	"github.com/mcuadros/go-defaults"
	"os"
)

type ObjectStoreConfig struct {
	DBConfig *DBStoreConfig
}

func (obsConf *ObjectStoreConfig) Load() error {
	obsConf.DBConfig = &DBStoreConfig{}
	if err := obsConf.DBConfig.Load(); err != nil {
		return err
	}

	return nil
}

type DBStoreConfig struct {
	Type     string `default:"redis"`     // STORE_TYPE
	Host     string `default:"localhost"` // DB_Host
	Port     string `default:"6379"`      // DB_PORT
	User     string `default:""`     // DB_USER
	Password string `default:""`     // DB_PASSWORD
}

func (dbConf *DBStoreConfig) Load() error {
	defaults.SetDefaults(dbConf)
	if val, ok := os.LookupEnv("STORE_TYPE"); ok {
		dbConf.Type = val
	}

	if val, ok := os.LookupEnv("DB_Host"); ok {
		dbConf.Host = val
	}

	if val, ok := os.LookupEnv("DB_PORT"); ok {
		dbConf.Port = val
	}

	if val, ok := os.LookupEnv("DB_USER"); ok {
		dbConf.User = val
	}

	if val, ok := os.LookupEnv("DB_PASSWORD"); ok {
		dbConf.Password = val
	}

	return nil
}

func Load() (*ObjectStoreConfig, error) {
	obsConf := &ObjectStoreConfig{}
	if err := obsConf.Load(); err != nil {
		return nil, err
	}

	return obsConf, nil
}
