package config

import "github.com/spf13/viper"

type Config struct {
	StepOnePort   string `mapstructure:"STEP1_PORT"`
	StepTwoPort   string `mapstructure:"STEP2_PORT"`
	StepThreePort string `mapstructure:"STEP3_PORT"`
	StepFourPort  string `mapstructure:"STEP4_PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
