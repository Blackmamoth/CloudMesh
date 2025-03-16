package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type AppConfiguration struct {
	ENVIRONMENT   string `envconfig:"ENVIRONMENT"   required:"true"`
	HOST          string `envconfig:"HOST"          required:"true"`
	PORT          string `envconfig:"PORT"          required:"true"`
	FRONTEND_HOST string `envconfig:"FRONTEND_HOST" required:"true"`
}

type PostgreConfiguration struct {
	POSTGRES_HOST    string `envconfig:"POSTGRES_HOST"     required:"true"`
	POSTGRES_PORT    string `envconfig:"POSTGRES_PORT"     required:"true"`
	POSTGRES_USER    string `envconfig:"POSTGRES_USER"     required:"true"`
	POSTGRES_PASS    string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	POSTGRES_DBNAME  string `envconfig:"POSTGRES_DB"       required:"true"`
	POSTGRES_SSLMODE string `envconfig:"POSTGRES_SSLMODE"  required:"true"`
}

type OAuthConfiguration struct {
	SUPPORTED_PROVIDERS []string
	GOOGLE              OAuthProviderConfig `yaml:"google"`
	DROPBOX             OAuthProviderConfig `yaml:"dropbox"`
}

type OAuthProviderConfig struct {
	CLIENT_ID     string   `yaml:"client_id"`
	CLIENT_SECRET string   `yaml:"client_secret"`
	CALLBACK_URL  string   `yaml:"callback_url"`
	SCOPES        []string `yaml:"scopes"`
}

type JWTConfiguration struct {
	ACCESS_TOKEN_SECRET              string `envconfig:"ACCESS_TOKEN_SECRET"              required:"true"`
	REFRESH_TOKEN_SECRET             string `envconfig:"REFRESH_TOKEN_SECRET"             required:"true"`
	ACCESS_TOKEN_EXPIRATION_IN_HOURS int    `envconfig:"ACCESS_TOKEN_EXPIRATION_IN_HOURS" required:"true"`
	REFRESH_TOKEN_EXPIRATION_IN_DAYS int    `envconfig:"REFRESH_TOKEN_EXPIRATION_IN_DAYS" required:"true"`
}

var (
	AppConfig      AppConfiguration
	PostgresConfig PostgreConfiguration
	JWTConfig      JWTConfiguration
	OAuthConfig    OAuthConfiguration
)

func init() {
	loadEnv()
	loadConfig()
}

func loadEnv() {
	godotenv.Load()

	if err := envconfig.Process("", &AppConfig); err != nil {
		log.Fatalf("An error occured while loading environment variables: %v", err)
	}

	if err := envconfig.Process("", &PostgresConfig); err != nil {
		log.Fatalf("An error occured while loading environment variables: %v", err)
	}

	if err := envconfig.Process("", &JWTConfig); err != nil {
		log.Fatalf("An error occured while loading environment variables: %v", err)
	}
}

func loadConfig() {
	data, err := os.ReadFile("oauth.config.yml")
	if err != nil {
		log.Fatalf("An error occured while reading `oauth.config.yaml`: %v", err)
	}

	OAuthConfig = OAuthConfiguration{
		SUPPORTED_PROVIDERS: []string{"google", "dropbox"},
	}

	if err := yaml.Unmarshal(data, &OAuthConfig); err != nil {
		log.Fatalf("An error occured while parsing `oauth.config.yaml`: %v", err)
	}
}
