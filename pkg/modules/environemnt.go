package modules

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Environment struct {
	viper.Viper
	name string
	log  *logrus.Logger
}

func NewEnvironment(name string) (*Environment, error) {
	env := &Environment{
		name: name,
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", name))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", name))
	viper.AddConfigPath(".")
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err := viper.ReadInConfig()
		if err != nil {
			return
		}
	})
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	log := logrus.New()
	logLevel := viper.Get("log.level")
	switch logLevel {
	case "debug":
		log.Level = logrus.DebugLevel
	case "info":
		log.Level = logrus.InfoLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	default:
		log.Level = logrus.InfoLevel
	}
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = viper.GetString("log.timestamp")
	formatter.FullTimestamp = true
	log.Formatter = formatter
	env.log = log

	return env, nil
}
