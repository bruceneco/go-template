package config

type EnvType string

const (
	EnvTypeDevelopment EnvType = "development"
	EnvTypeStaging     EnvType = "staging"
	EnvTypeProduction  EnvType = "production"
)
