package server

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	addr     string
	basePath string
}

func DefaultConfig() Config {
	home, err := homedir.Dir()
	fmt.Println(home)
	if err != nil {
		panic(err)
	}
	if !strings.HasSuffix(home, "/") {
		home += "/"
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".fileserver")
	viper.SetConfigType("yaml")
	viper.SetDefault("addr", ":11111")
	viper.SetDefault("basePath", home)
	err = viper.ReadInConfig()
	if err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}
	addr := viper.GetString("addr")
	basePath := viper.GetString("basePath")
	return Config{
		addr:     addr,
		basePath: basePath,
	}
}
