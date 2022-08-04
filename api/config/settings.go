package config

import (
	"os"
	"strings"
)

type Settings struct {
	CurrentEnv         string
	DatadogEnvironment string
	DatadogHost        string
	DatadogService     string
	GitSHA             string
	Port               string
	PostgresURL        string
}

var GitSHA string

func Init() *Settings {
	env := strings.ToLower(os.Getenv("AL_API_ENV"))

	settings := settingsLocal
	if env == "production" {
		settings = settingsProd
	}

	settings.CurrentEnv = env

	settings.DatadogEnvironment = os.Getenv("DD_ENV")
	settings.DatadogHost = os.Getenv("DD_AGENT_HOST")
	settings.DatadogService = os.Getenv("DD_SERVICE")

	settings.GitSHA = GitSHA

	settings.Port = os.Getenv("PORT")
	if settings.Port == "" {
		settings.Port = "8080"
	}

	settings.PostgresURL = os.Getenv("AL_DATABASE_URL")

	return &settings
}
