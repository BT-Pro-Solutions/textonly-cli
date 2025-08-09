package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	configDirName  = "textonly"
	configFileName = "config"
	configFileType = "yaml"
)

func Init() error {
	viper.SetEnvPrefix("TO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("api", "https://textonly.io/api")
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(Dir())
	_ = viper.ReadInConfig()
	return nil
}

func userConfigRoot() string {
	if v, err := os.UserConfigDir(); err == nil {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}

func Dir() string {
	return filepath.Join(userConfigRoot(), configDirName)
}

func ensureDir() error { return os.MkdirAll(Dir(), 0o755) }

func Path() string {
	_ = ensureDir()
	return filepath.Join(Dir(), configFileName+"."+configFileType)
}

func APIBaseURL() string {
	v := viper.GetString("api")
	if v == "" { v = "https://textonly.io/api" }
	return strings.TrimRight(v, "/")
}

func Get(key string) (string, bool) {
	if !viper.IsSet(key) {
		return "", false
	}
	return fmt.Sprint(viper.Get(key)), true
}

func Set(key, value string) error {
	if err := ensureDir(); err != nil { return err }
	viper.Set(key, value)
	return viper.WriteConfigAs(Path())
}

func Unset(key string) error {
	if !viper.IsSet(key) {
		return errors.New("key not set")
	}
	if err := ensureDir(); err != nil { return err }
	viper.Set(key, nil)
	return viper.WriteConfigAs(Path())
}
