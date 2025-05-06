package competitors

import (
	"time"

	"github.com/ivanjabrony/yadro/cmd/config"
)

type CompetitorStatus string

var (
	NotRegistered   CompetitorStatus = "not registered"
	Registered      CompetitorStatus = "registered"
	OnStartLine     CompetitorStatus = "on start line"
	NotStarted      CompetitorStatus = "not started"
	RunningMain     CompetitorStatus = "running main loop"
	RunningPenalty  CompetitorStatus = "running penalty loop"
	OnFiringRange   CompetitorStatus = "on firing range"
	LeftFiringRange CompetitorStatus = "left firing range"
	CantContinue    CompetitorStatus = "can't continue"
	NotFinished     CompetitorStatus = "not finished"
	Disqualified    CompetitorStatus = "discualified"
	Finished        CompetitorStatus = "finished"
	ErrorState      CompetitorStatus = "error state"
)

type Competitor struct {
	CurrentStatus        CompetitorStatus
	StartTime            time.Time
	ActualStartTime      time.Time
	LapsRemained         int
	FiringLinesRemained  int
	MainLapsStartTime    []time.Time
	MainLapsEndTime      []time.Time
	PenaltyLapsStartTime []time.Time
	PenaltyLapsEndTime   []time.Time
	NumberOfHits         int
}

func New(cfg *config.Config) *Competitor {
	return &Competitor{
		NotRegistered,
		time.Time{},
		time.Time{},
		cfg.Laps,
		cfg.FiringLines,
		make([]time.Time, cfg.Laps),
		make([]time.Time, cfg.Laps),
		make([]time.Time, 0),
		make([]time.Time, 0),
		0,
	}
}
