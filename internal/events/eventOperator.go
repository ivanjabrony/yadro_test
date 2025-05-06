package events

import (
	"bufio"
	"io"
	"os"

	"github.com/ivanjabrony/yadro/cmd/config"
	"github.com/ivanjabrony/yadro/internal/competitors"
	"github.com/ivanjabrony/yadro/internal/output"
)

type EventLogger interface {
	LogEvent(string, ...any)
}

type ProcessEvent func(
	EventLogger,
	*config.Config,
	map[int]*competitors.Competitor,
	*Event) error

type EventOperator struct {
	processors map[int]ProcessEvent
	personData map[int]*competitors.Competitor
	config     *config.Config
	logger     EventLogger
	output     io.Writer
}

func NewEventOperator(config *config.Config, logger EventLogger, output io.Writer) *EventOperator {
	processors := map[int]ProcessEvent{
		CompetitorRegisteredEventID:         RegistrationEP,
		StartTimeSetByDrawEventID:           StartTimeSetByDrawEP,
		CompetitorOnStartLineEventID:        OnStartLineEP,
		CompetitorStartedEventID:            CompetitorStartedEP,
		CompetitorOnFiringRangeEventID:      OnFiringRangeEP,
		TargetBeenHitEventID:                TargetBeenHitEP,
		CompetitorLeftFiringRangeEventID:    LeftFiringRangeEP,
		CompetitorEnteredPenaltyLapsEventID: EnteredPenaltyEP,
		CompetitorEndedPenaltyLapsEventID:   LeftPenaltyEP,
		CompetitorEndedMainLapEventID:       EndedMainLapEP,
		CompetitorCantContinueEventID:       CantContinueEP,
	}
	return &EventOperator{
		processors,
		make(map[int]*competitors.Competitor),
		config,
		logger,
		output,
	}
}

func (eo *EventOperator) RunFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		event, err := ParseEvent(scanner.Text())
		if err != nil {
			return err
		}

		processor, ok := eo.processors[event.EventID]
		if ok {
			err = processor(eo.logger, eo.config, eo.personData, event)
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	output.CalculateOutput(eo.config, eo.output, eo.personData)
	return nil
}
