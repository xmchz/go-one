package config

type Config interface {
	GetInt(string) int
	GetStr(string) string
}

func FromFile(confFile string) (Config, error) {
	return newVp(confFile)
}

/*
app.yml
app-dev.yml
app-prod.yml
*/
func FromFiles(baseFile string, modeKey string) (Config, error) {
	return newVpWithMode(baseFile, modeKey)
}
