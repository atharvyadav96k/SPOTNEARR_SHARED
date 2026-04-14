package app

import "github.com/atharvyadav96k/SPOTNEARR_SHARED/app/models/secrets"

func Init() App {
	return App{}
}

func (a *App) InitEnvironmentVariables() {
	a.Env = secrets.NewSecrets()
}
