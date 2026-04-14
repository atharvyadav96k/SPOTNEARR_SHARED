package secrets

import "os"

func NewSecrets() Env {
	env := Env{}
	env.LoadSecrets()
	return env
}

func (e *Env) LoadSecrets() {
	e.GCP_PROJECT_ID = os.Getenv("GCP_PROJECT_ID")
}
