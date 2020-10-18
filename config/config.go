package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config interface {
	GetInt(string) int
	GetStr(string) string
}

func FromFile(confFile string) (Config, error) {
	vp := viper.New()
	vp.SetConfigFile(confFile)
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	return &viperConfig{vp: vp}, nil
}

func FromFiles(baseFile string, modeKey string) (Config, error) {
	base := viper.New()
	base.SetConfigFile(baseFile) // app.yml
	if err := base.ReadInConfig(); err != nil {
		return nil, err
	}
	vp := viper.New()
	for k, v := range base.AllSettings() {
		vp.SetDefault(k, v)
	}
	mode := base.GetString(modeKey)
	ss := strings.Split(baseFile, ".")
	suffix := "." + ss[len(ss)-1]
	vp.SetConfigFile(strings.Replace(baseFile, suffix, fmt.Sprintf("-%s%s", mode, suffix), 1)) // app-mode.yml
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	return &viperConfig{vp}, nil
}
