package competitors

import (
	"errors"
)

var ErrInvalidStatusData error = errors.New("competitor's status is invalid")
var ErrInvalidAmountOfMainLapsData error = errors.New("competitor's main laps data is invalid")

func ValidateCompetitor(competitor *Competitor) error {
	if competitor.CurrentStatus != Finished &&
		competitor.CurrentStatus != NotFinished &&
		competitor.CurrentStatus != Disqualified {
		return ErrInvalidStatusData
	}
	if len(competitor.MainLapsStartTime) != len(competitor.MainLapsEndTime) {
		return ErrInvalidAmountOfMainLapsData
	}
	return nil
}
