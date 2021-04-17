package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
)

/*
app.yml
app-dev.yml
app-prod.yml
*/
const fmtModeFile = "-%s%s"

func newVp(confFile string) (*vp, error) {
	v := viper.New()
	v.SetConfigFile(confFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.ReadInConfig(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("reload config file")
		}
	})
	v.WatchConfig()
	return &vp{Viper: v}, nil
}

func newVpWithMode(baseFile, modeKey string) (*vp, error) {
	b := viper.New()
	b.SetConfigFile(baseFile) // app.yml
	if err := b.ReadInConfig(); err != nil {
		return nil, err
	}
	c := viper.New()
	for k, v := range b.AllSettings() {
		c.SetDefault(k, v)
	}
	mode := b.GetString(modeKey)
	ss := strings.Split(baseFile, ".")
	suffix := "." + ss[len(ss)-1]
	c.SetConfigFile(strings.Replace(baseFile, suffix, fmt.Sprintf(fmtModeFile, mode, suffix), 1)) // app-{mode}.yml
	if err := c.ReadInConfig(); err != nil {
		return nil, err
	}
	return &vp{Viper: c}, nil
}

type vp struct {
	*viper.Viper
}

func (v *vp) GetInt(key string) int {
	return v.Viper.GetInt(key)
}

func (v *vp) GetStr(key string) string {
	return v.Viper.GetString(key)
}

