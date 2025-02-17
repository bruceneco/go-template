package config

type EnvType string

const (
	EnvTypeDevelopment EnvType = "development"
	EnvTypeStaging     EnvType = "staging"
	EnvTypeProduction  EnvType = "production"
	EnvTypeTest        EnvType = "test"
)

func (e EnvType) String() string {
	return string(e)
}
