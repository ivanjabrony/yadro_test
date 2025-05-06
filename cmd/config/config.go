package config

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/ivanjabrony/yadro/internal/utils"
)

type Config struct {
	Laps        int           `json:"laps"`
	LapLen      int           `json:"lapLen"`
	PenaltyLen  int           `json:"penaltyLen"`
	FiringLines int           `json:"firingLines"`
	Start       time.Time     `json:"start"`
	StartDelta  time.Duration `json:"startDelta"`
}

type unformattedConfig struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

func New(filepath string) *Config {
	unformattedConfig := &unformattedConfig{}
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(file)
	jsonConfig, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(jsonConfig, unformattedConfig)
	if err != nil {
		log.Fatal(err)
	}

	cfg := Config{}
	cfg.Laps = unformattedConfig.Laps
	cfg.LapLen = unformattedConfig.LapLen
	cfg.PenaltyLen = unformattedConfig.PenaltyLen
	cfg.FiringLines = unformattedConfig.FiringLines
	cfg.Start, err = utils.UnmarshallTimeStamp(unformattedConfig.Start)
	if err != nil {
		log.Fatal(err)
	}
	cfg.StartDelta, err = utils.UnmarshallDuration(unformattedConfig.StartDelta)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
