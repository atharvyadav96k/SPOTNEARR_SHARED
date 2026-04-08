package app

import "github.com/atharvyadav96k/SPOTNEARR_SHARED/app/env"

func Init() *App {
	return &App{
		Env: *env.Init(),
	}
}
