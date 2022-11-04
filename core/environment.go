package core

import "os"

var EnvironmentVariables = []EnvironmentVariable{
	EEnvironmentVariable.AuthToken(),
}

type EnvironmentVariable struct {
	Name    string
	Default string
	Secret  bool
}

func (e EnvironmentVariable) Get() (val string, defaulted bool) {
	val = os.Getenv(e.Name)

	if val == "" {
		val = e.Default
		defaulted = true
	}

	return
}

type eEnvironmentVariable struct{}

var EEnvironmentVariable = &eEnvironmentVariable{}

func (*eEnvironmentVariable) AuthToken() EnvironmentVariable {
	return EnvironmentVariable{
		Name:   "AOCF_SESSION_COOKIE",
		Secret: true,
	}
}
