package config

/*
service --- config | ---- info from file/net/...
config 侧重于固定结构的信息校验、读取、覆盖
进一步则需要考虑推送信息变更
 */

type Config interface {
	GetInt(string) int
	GetStr(string) string
}

func FromFile(confFile string) (Config, error) {
	return newVp(confFile)
}

func FromFiles(baseFile string, modeKey string) (Config, error) {
	return newVpWithMode(baseFile, modeKey)
}

