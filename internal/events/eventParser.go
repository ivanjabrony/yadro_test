package events

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ivanjabrony/yadro/internal/utils"
)

func ParseEvent(event string) (*Event, error) {
	splittedEvent := strings.Split(event, " ")
	if len(splittedEvent) < 3 || len(splittedEvent) > 4 {
		return nil, errors.New("invalid Event format has been passed")
	}

	time, err := utils.UnmarshallTimeStamp(splittedEvent[0][1 : len(splittedEvent[0])-1])
	if err != nil {
		return nil, err
	}
	eventID, err := strconv.Atoi(splittedEvent[1])
	if err != nil || eventID < 0 {
		return nil, errors.New("invalid EventID format has been passed")
	}
	competitorID, err := strconv.Atoi(splittedEvent[2])
	if err != nil || competitorID < 0 {
		return nil, errors.New("invalid CompetitorID format has been passed")
	}

	parsedEvent := Event{time, eventID, competitorID, ""}

	if len(splittedEvent) == 4 {
		parsedEvent.ExtraParams = splittedEvent[3]
	}

	return &parsedEvent, nil
}
