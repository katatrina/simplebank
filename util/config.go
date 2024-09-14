package util

import "github.com/spf13/viper"

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DriverName     string `mapstructure:"DRIVER_NAME"`
	DataSourceName string `mapstructure:"DATASOURCE_NAME"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.UnmarshalExact(&config)
	if err != nil {
		return
	}

	return
}
