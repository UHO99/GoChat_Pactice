package config

import "github.com/spf13/viper"

type Config struct {
	DBSource      string `mapstructure:"DATABASE_URL"`
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	StepOnePort   string `mapstructure:"STEP1_PORT"`
	StepTwoPort   string `mapstructure:"STEP2_PORT"`
	StepThreePort string `mapstructure:"STEP3_PORT"`
	StepFourPort  string `mapstructure:"STEP4_PORT"`
	StepFivePort  string `mapstructure:"STEP5_PORT"`
	StepSixPort   string `mapstructure:"STEP6_PORT"`
	StepSevenPort string `mapstructure:"STEP7_PORT"`
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
