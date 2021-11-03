package utils

import (
	"MyEnvelope/algo"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func Run(cfg string) error {
	c := Config{
		Name: cfg,
	}

	if err := c.init(); err != nil {
		return err
	}

	c.watchConfig()

	return nil
}
func (c *Config) init() error {
	v := viper.New()
	v.AddConfigPath("./configs")
	v.SetConfigType("yaml")
	v.SetConfigName("default")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	cfgs := v.AllSettings()

	for k, v := range cfgs {
		viper.SetDefault(k, v)
	}
	viper.SetConfigFile("./configs/default.yaml")

	if c.Name != "" {
		viper.AddConfigPath("./configs")
		viper.SetConfigName(c.Name)
		viper.SetConfigType("yaml")
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	env := os.Getenv("GO_ENV")
	if env != "" {
		fmt.Println(env)
		viper.SetConfigName(env)
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yaml")
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	return nil
}

func (c *Config) watchConfig() {
	fmt.Println("start watch config")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		algo.InitConfig()
	})
}
