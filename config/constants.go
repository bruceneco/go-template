package config

type EnvType string

const (
	EnvTypeDevelopment EnvType = "development"
	EnvTypeStaging     EnvType = "staging"
	EnvTypeProduction  EnvType = "production"
)

func (e *EnvType) String() string {
	return string(*e)
}
