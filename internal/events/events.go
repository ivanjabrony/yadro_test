package events

import "time"

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  string
}

const (
	CompetitorRegisteredEventID         = 1
	StartTimeSetByDrawEventID           = 2
	CompetitorOnStartLineEventID        = 3
	CompetitorStartedEventID            = 4
	CompetitorOnFiringRangeEventID      = 5
	TargetBeenHitEventID                = 6
	CompetitorLeftFiringRangeEventID    = 7
	CompetitorEnteredPenaltyLapsEventID = 8
	CompetitorEndedPenaltyLapsEventID   = 9
	CompetitorEndedMainLapEventID       = 10
	CompetitorCantContinueEventID       = 11

	CompetitorDisqualifiedEventID = 32
	CompetitorHasFinishedEventID  = 33
)
