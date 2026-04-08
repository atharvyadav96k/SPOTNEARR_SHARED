package app

import env "github.com/atharvyadav96k/SPOTNEARR_API/shared/app/Env"

func Init() *App {
	return &App{
		Env: env.Init(),
	}
}
