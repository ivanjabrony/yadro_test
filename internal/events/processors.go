package events

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ivanjabrony/yadro/cmd/config"
	"github.com/ivanjabrony/yadro/internal/competitors"
	"github.com/ivanjabrony/yadro/internal/utils"
)

var ( // TODO: вынести в отдельный пакет
	ErrImpossibleState        = errors.New("impossible state of competitor")
	ErrInvalidTimeFormat      = errors.New("time of invalid format has been passed")
	ErrInvalidHitFormat       = errors.New("hit number of invalid format has been passed")
	ErrInvalidFireRangeFormat = errors.New("fire range number of invalid format has been passed")
)

var RegistrationEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorRegisteredEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorRegisteredEventID)
	}

	if _, ok := competitorsData[e.CompetitorID]; !ok {
		competitorsData[e.CompetitorID] = competitors.New(cfg)
	}
	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.NotRegistered:
		curCompetitor.CurrentStatus = competitors.Registered
		logger.LogEvent("[%v] The competitor(%v) registered\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var StartTimeSetByDrawEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != StartTimeSetByDrawEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, StartTimeSetByDrawEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.Registered:
		startTime, err := utils.UnmarshallTimeStamp(e.ExtraParams)
		if err != nil {
			return ErrInvalidTimeFormat
		}
		curCompetitor.StartTime = startTime
		logger.LogEvent("[%v] The start time for the competitor(%v) was set by a draw to %v\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID, e.ExtraParams)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var OnStartLineEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorOnStartLineEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorOnStartLineEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.Registered:
		curCompetitor.CurrentStatus = competitors.OnStartLine
		logger.LogEvent("[%v] The competitor(%v) is on the start line\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var CompetitorStartedEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorStartedEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorStartedEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.OnStartLine:
		curCompetitor.CurrentStatus = competitors.RunningMain
		curLap := cfg.Laps - curCompetitor.LapsRemained
		curCompetitor.MainLapsStartTime[curLap] = e.Time
		logger.LogEvent("[%v] The competitor(%v) has started\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var OnFiringRangeEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorOnFiringRangeEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorOnFiringRangeEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]
	firingRangeNumber, err := strconv.Atoi(e.ExtraParams)
	if err != nil || firingRangeNumber <= 0 {
		return err
	}
	logger.LogEvent("[%v] The competitor(%v) is on the firing range(%v)\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID, e.ExtraParams)

	switch curCompetitor.CurrentStatus {
	case competitors.RunningMain:
		curCompetitor.CurrentStatus = competitors.OnFiringRange

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var TargetBeenHitEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != TargetBeenHitEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, TargetBeenHitEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]
	targetNumber, err := strconv.Atoi(e.ExtraParams)
	if err != nil || targetNumber > 5 || targetNumber <= 0 {
		return ErrInvalidHitFormat
	}

	switch curCompetitor.CurrentStatus {
	case competitors.OnFiringRange:
		curCompetitor.NumberOfHits++
		logger.LogEvent("[%v] The target(%v) has been hit by competitor(%v)\n", utils.MarshallTimeStamp(e.Time), e.ExtraParams, e.ExtraParams)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var LeftFiringRangeEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorLeftFiringRangeEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorLeftFiringRangeEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.OnFiringRange:
		curCompetitor.CurrentStatus = competitors.LeftFiringRange
		curCompetitor.FiringLinesRemained--
		logger.LogEvent("[%v] The competitor(%v) left the firing range\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var EnteredPenaltyEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorEnteredPenaltyLapsEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorEnteredPenaltyLapsEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.LeftFiringRange:
		curCompetitor.CurrentStatus = competitors.RunningPenalty
		curCompetitor.PenaltyLapsStartTime = append(curCompetitor.PenaltyLapsStartTime, e.Time)
		logger.LogEvent("[%v] The competitor(%v) entered the penalty laps\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var LeftPenaltyEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorEndedPenaltyLapsEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorEndedPenaltyLapsEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.RunningPenalty:
		curCompetitor.CurrentStatus = competitors.RunningMain
		curCompetitor.PenaltyLapsEndTime = append(curCompetitor.PenaltyLapsEndTime, e.Time)
		logger.LogEvent("[%v] The competitor(%v) left the penalty laps\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var EndedMainLapEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorEndedMainLapEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorEndedMainLapEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.RunningMain, competitors.LeftFiringRange:
		curCompetitor.CurrentStatus = competitors.RunningMain
		curLap := cfg.Laps - curCompetitor.LapsRemained + 1
		curCompetitor.MainLapsEndTime[curLap-1] = e.Time
		curCompetitor.LapsRemained--

		if curCompetitor.LapsRemained == 0 {
			curCompetitor.CurrentStatus = competitors.Finished
		} else {
			curCompetitor.MainLapsStartTime[curLap] = e.Time
		}
		logger.LogEvent("[%v] The competitor(%v) ended the main lap\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID)

	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.ErrorState
		return fmt.Errorf("%w: Current eventID:(%v), CurrentCompetitorID(%v)", ErrImpossibleState, e.EventID, e.CompetitorID)
	}

	return nil
}

var CantContinueEP ProcessEvent = func(
	logger EventLogger,
	cfg *config.Config,
	competitorsData map[int]*competitors.Competitor,
	e *Event) error {
	if e.EventID != CompetitorCantContinueEventID {
		return fmt.Errorf("wrong event has been passed: %v, needed %v", e.EventID, CompetitorCantContinueEventID)
	}

	curCompetitor := competitorsData[e.CompetitorID]

	switch curCompetitor.CurrentStatus {
	case competitors.ErrorState:
		return nil

	default:
		curCompetitor.CurrentStatus = competitors.CantContinue
		logger.LogEvent("[%v] The competitor(%v) can`t continue: %v\n", utils.MarshallTimeStamp(e.Time), e.CompetitorID, e.ExtraParams)

		return nil
	}
}
