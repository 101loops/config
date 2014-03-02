package config

type Env string

// IsDev returns true if the environment is 'development' (or empty).
func (self Env) IsDev() bool {
	env := string(self)
	return env == "" || env == "development"
}

// IsTest returns true if the environment is 'testing'.
func (self Env) IsTest() bool {
	return string(self) == "testing"
}

// IsStage returns true if the environment is 'staging'.
func (self Env) IsStage() bool {
	return string(self) == "staging"
}

// IsProd returns true if the environment is 'production'.
func (self Env) IsProd() bool {
	return string(self) == "production"
}
