package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ntfs32/gvm/src/gvm/file"
	"fmt"
	"bytes"
	"errors"
	"os"
)


type version struct{
	ReleaseVersion []string
	InstalledVersion []string
	UsedVersion string
}

type Config struct {
	Version              version
	GolangAddress        string
	GithubAddress        string
	ReleaseHistoryAddress string
}

var config Config

func init()  {
	if _, err := toml.DecodeFile(file.ConfigFile, &config); err != nil {
		fmt.Println(err)
	}
}
func Get() *Config  {
	return &config
}

func (config *Config) Save() bool {
	var secondBuffer bytes.Buffer
	e2 := toml.NewEncoder(&secondBuffer)
	err := e2.Encode(config)
	if err != nil {
		errors.New("parse config failed")
		return false
	}
	f, err := os.Create(file.ConfigFile)
	defer f.Close()
	_, err = secondBuffer.WriteTo(f)
	if err != nil {
		errors.New("write config file failed")
		return false
	}
	return true
}