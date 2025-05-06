package events_test

import (
	"testing"
	"time"

	"github.com/ivanjabrony/yadro/internal/events"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		eventString   string
		expectedEvent *events.Event
		isErrExpected bool
	}{
		{
			name:        "valid event",
			eventString: "[09:05:59.867] 1 1",
			expectedEvent: &events.Event{
				Time:         time.Date(0, 1, 1, 9, 5, 59, 867000000, time.UTC),
				EventID:      1,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			isErrExpected: false,
		},
		{
			name:        "valid event with params",
			eventString: "[09:15:00.841] 2 1 09:30:00.000",
			expectedEvent: &events.Event{
				Time:         time.Date(0, 1, 1, 9, 15, 0, 841000000, time.UTC),
				EventID:      2,
				CompetitorID: 1,
				ExtraParams:  "09:30:00.000",
			},
			isErrExpected: false,
		},
		{
			name:          "invalid empty string",
			eventString:   "",
			expectedEvent: nil,
			isErrExpected: true,
		},
		{
			name:          "invalid not enough params",
			eventString:   "[09:05:59.867] 1",
			expectedEvent: nil,
			isErrExpected: true,
		},
		{
			name:          "invalid time format",
			eventString:   "[09:05:59] 1 1",
			expectedEvent: nil,
			isErrExpected: true,
		},
		{
			name:          "invalid EventID",
			eventString:   "[09:05:59.867] abc 1",
			expectedEvent: nil,
			isErrExpected: true,
		},
		{
			name:          "invalid competitorID",
			eventString:   "[09:05:59.867] 1 abc",
			expectedEvent: nil,
			isErrExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := events.ParseEvent(tt.eventString)

			if tt.isErrExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedEvent, got)
		})
	}
}
