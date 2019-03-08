package modules

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/x-cray/logrus-prefixed-formatter"
	"strings"
	"time"
)

type Environment struct {
	viper.Viper
	config   string
	provider string
	uri      string
	name     string
	log      *logrus.Logger
	remote   bool
}

type EnvironmentBuilder struct {
	*Environment
	errors []error
}

func NewEnvironmentBuilder() *EnvironmentBuilder {
	return &EnvironmentBuilder{
		Environment: &Environment{},
		errors:      make([]error, 0),
	}
}

func (eb *EnvironmentBuilder) WithName(name string) *EnvironmentBuilder {
	eb.name = name
	return eb
}

func (eb *EnvironmentBuilder) WithRemoteConfiguration(provider, uri string) *EnvironmentBuilder {
	eb.remote = true
	eb.provider = provider
	eb.uri = uri
	return eb
}

func (eb *EnvironmentBuilder) WithConfigurationFile(config string) *EnvironmentBuilder {
	eb.remote = false
	eb.config = config
	return eb
}

func (eb *EnvironmentBuilder) verifyConfig() error {
	stack := make([]string, 0)
	if eb.name == "" {
		stack = append(stack, "name cannot be empty")
	}

	if eb.config == "" {
		stack = append(stack, "a configuration file must be provided")
	}

	if (eb.provider == "" || eb.uri == "") && eb.remote {
		stack = append(stack, "provider and uri cannot be empty")
	}

	if len(stack) > 0 {
		return errors.New(strings.Join(stack, "\n"))
	}
	return nil
}

func (eb *EnvironmentBuilder) Build() (*Environment, error) {
	if err := eb.verifyConfig(); err != nil {
		return nil, err
	}
	log := logrus.New()
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	log.Formatter = formatter

	v := viper.New()
	v.SetConfigName(eb.config)

	if eb.provider != "" && eb.uri != "" {
		log.Infof("Configuration %s -> %s", eb.provider, eb.uri)
		if err := v.AddRemoteProvider(eb.provider, eb.uri, eb.name); err != nil {
			return nil, err
		}
		if err := v.ReadRemoteConfig(); err != nil {
			return nil, err
		}
	} else {
		log.Infof("Configuration %s", eb.config)
		v.AddConfigPath(fmt.Sprintf("/etc/%s/", eb.name))
		v.AddConfigPath(fmt.Sprintf("$HOME/.%s", eb.name))
		v.AddConfigPath(".")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
		v.WatchConfig()
		v.OnConfigChange(func(in fsnotify.Event) {
			err := v.ReadInConfig()
			if err != nil {
				return
			}
		})
	}

	env := eb.Environment
	logLevel := env.GetStringOrDefault("log.level", "info")
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
	if timestamp := env.GetStringOrDefault("log.timestamp", ""); timestamp != "" {
		formatter.TimestampFormat = timestamp
	}
	env.log = log

	err := env.CheckErrors(eb.errors)
	if err != nil {
		return nil, err
	}
	return env, nil
}

func (e *Environment) GetIntOrDefault(path string, defaultValue int) int {
	if value := e.GetInt(path); value != 0 {
		return value
	}
	return defaultValue
}

func (e *Environment) GetBoolOrDefault(path string, defaultValue bool) bool {
	if value := e.GetBool(path); value {
		return value
	}
	return defaultValue
}

func (e *Environment) GetFloat64OrDefault(path string, defaultValue float64) float64 {
	if value := e.GetFloat64(path); value != 0.0 {
		return value
	}
	return defaultValue
}

func (e *Environment) GetDurationOrDefault(path string, defaultValue time.Duration) time.Duration {
	if value := e.GetDuration(path); value != 0.0 {
		return value
	}
	return defaultValue
}

func (e *Environment) GetStringOrDefault(path string, defaultValue string) string {
	if value := e.GetString(path); value != "" {
		return value
	}
	return defaultValue
}

func (e *Environment) CheckErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	msg := ""
	for _, err := range errs {
		msg += fmt.Sprintf("%s\n", err.Error())
	}

	return errors.New(msg)

}
