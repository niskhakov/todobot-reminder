package config

import "github.com/spf13/viper"

type Config struct {
	TelegramToken       string
	TodoistClientID     string
	TodoistClientSecret string
	AuthServerURL       string
	TelegramBotURL      string `mapstructure:"bot_url"`
	DBPath              string `mapstructure:"db_file"`

	Messages Messages
}

type Messages struct {
	Errors    Errors
	Responses Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	NoMessage         string `mapstructure:"no_message"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	if err := viper.BindEnv("client_id"); err != nil {
		return err
	}
	if err := viper.BindEnv("client_secret"); err != nil {
		return err
	}
	if err := viper.BindEnv("auth_server_url"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("token")
	cfg.TodoistClientID = viper.GetString("client_id")
	cfg.TodoistClientSecret = viper.GetString("client_secret")
	cfg.AuthServerURL = viper.GetString("auth_server_url")

	return nil
}
