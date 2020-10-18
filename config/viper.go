package config

import "github.com/spf13/viper"

type viperConfig struct {
	vp *viper.Viper
}

func (v *viperConfig) GetInt(key string) int {
	return v.vp.GetInt(key)
}

func (v *viperConfig) GetStr(key string) string {
	return v.vp.GetString(key)
}

