package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	confFile = "./viper_test.yml"
	yml1     = "k1: v1"
	yml2     = "k1: v1-modified"
)

func TestConfig(t *testing.T) {
	f, err := os.Create(confFile)
	require.Nil(t, err)
	_, err = f.Write([]byte(yml1))
	require.Nil(t, err)
	require.Nil(t, f.Close())
	//err := ioutil.WriteFile(confFile, []byte(yml1), os.ModePerm)
	require.Nil(t, err)
	var conf Config
	conf, err = newVp(confFile)
	require.Nil(t, err)
	require.Equal(t, "v1", conf.GetStr("k1"))
	//err = ioutil.WriteFile(confFile, []byte(yml2), os.ModePerm)
	f, err = os.Create(confFile)
	require.Nil(t, err)
	_, err = f.Write([]byte(yml2))
	require.Nil(t, err)
	require.Nil(t, f.Sync())            // different with flush
	require.Nil(t, f.Close())
	require.Equal(t, "v1-modified", conf.GetStr("k1")) // not sync yet
}
