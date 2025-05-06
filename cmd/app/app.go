package app

import (
	"log"
	"os"

	"github.com/ivanjabrony/yadro/cmd/config"
	"github.com/ivanjabrony/yadro/internal/events"
	"github.com/ivanjabrony/yadro/internal/logger"
)

type App struct {
	cfg       *config.Config
	eo        *events.EventOperator
	eventPath string
}

func New(eventPath, configPath, outputPath string) App {
	cfg := config.New(configPath)
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("Output file couldn't be created")
	}
	eo := events.NewEventOperator(cfg, logger.MyLogger{}, file)
	return App{cfg, eo, eventPath}
}

func (app *App) Run() {
	app.eo.RunFromFile(app.eventPath)
}
