package main

import (
	"log"
	"os"

	"github.com/ivanjabrony/yadro/cmd/app"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Not enough arguments: pls pass the filepaths in the exact order: events file, config file, destination file")
	}
	eventPath := os.Args[1]
	cfgPath := os.Args[2]
	outputPath := os.Args[3]
	app := app.New(eventPath, cfgPath, outputPath)
	app.Run()
}
