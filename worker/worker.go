package worker

import "app"

type Worker struct {
	App *app.App
}

func NewWorker(a *app.App) Worker {
	return Worker{App: a}
}
