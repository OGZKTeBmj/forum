package app

import (
	"log/slog"
	"strconv"

	"github.com/OGZKTeBmj/forum/utils"
)

type Controller interface {
	Run(addr string) error
}

type App struct {
	controller Controller
	log        *slog.Logger
	port       int
}

func New(Controller Controller, log *slog.Logger, port int) *App {
	return &App{
		controller: Controller,
		log:        log,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "app.Run"

	log := a.log.With(
		"op", op,
		"port", a.port,
	)

	log.Info("app is running")

	addr := ":" + strconv.Itoa(a.port)

	if err := a.controller.Run(addr); err != nil {
		return utils.ErrWrap(op, err)
	}

	return nil
}
