package output

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/ivanjabrony/yadro/cmd/config"
	"github.com/ivanjabrony/yadro/internal/competitors"
	"github.com/ivanjabrony/yadro/internal/utils"
)

type outputData struct {
	competitorID int
	totalTime    string
	mainLaps     []*LapData
	penaltyLaps  *LapData
	hits         int
	shots        int
}

type LapData struct {
	avgSpeed  float64
	timeTaken string
}

func CalculateOutput(cfg *config.Config, output io.Writer, personData map[int]*competitors.Competitor) {
	outputSlice := make([]*outputData, len(personData))
	pos := 0
	for id, competitor := range personData {
		err := competitors.ValidateCompetitor(competitor)
		if err != nil {
			_, _ = fmt.Fprintf(output, "Competitor(%v) has invalid data", id)
			continue
		}

		if competitor.CurrentStatus == competitors.NotStarted {
			outputSlice[pos] = &outputData{competitorID: id, totalTime: "NotStarted"}
			pos++
			continue
		}

		competitorData := outputData{competitorID: id, penaltyLaps: &LapData{}}
		totalTime := competitor.ActualStartTime.Sub(competitor.StartTime)

		for i := range len(competitor.MainLapsEndTime) {
			lapTime := competitor.MainLapsEndTime[i].Sub(competitor.MainLapsStartTime[i])
			totalTime += lapTime
			avgSpeed := float64(cfg.LapLen) / lapTime.Seconds()

			competitorData.mainLaps = append(competitorData.mainLaps, &LapData{avgSpeed, utils.MarshallDurationToTimestamp(lapTime)})
		}

		var totalPenaltyTime time.Duration
		totaPenaltyLen := (cfg.FiringLines*5 - competitor.NumberOfHits) * cfg.PenaltyLen
		for i := range competitor.PenaltyLapsEndTime {
			totalTime += competitor.PenaltyLapsEndTime[i].Sub(competitor.PenaltyLapsStartTime[i])
			totalPenaltyTime += competitor.PenaltyLapsEndTime[i].Sub(competitor.PenaltyLapsStartTime[i])
		}
		competitorData.totalTime = utils.MarshallDurationToTimestamp(totalTime)
		competitorData.penaltyLaps.timeTaken = utils.MarshallDuration(totalPenaltyTime)
		competitorData.penaltyLaps.avgSpeed = float64(totaPenaltyLen) / totalPenaltyTime.Seconds()
		competitorData.shots = cfg.FiringLines * 5
		competitorData.hits = competitor.NumberOfHits

		switch competitor.CurrentStatus {
		case competitors.Disqualified, competitors.CantContinue, competitors.NotFinished:
			competitorData.totalTime = "NotFinished"
		}
		outputSlice[pos] = &competitorData
		pos++
	}

	printCompetitorsData(output, outputSlice)
}

func printCompetitorsData(output io.Writer, competitors []*outputData) {
	outputFormat := "[%v] %v [%v] {%v, %v} %v/%v\n"
	mainLapFormat := "{%v, %v}"
	dataToOrder := make([]string, len(competitors))

	for i, competitor := range competitors {
		var mainLapString strings.Builder
		penaltySpeed := ""
		for i, lap := range competitor.mainLaps {
			if math.IsNaN(lap.avgSpeed) {
				_, _ = mainLapString.WriteString(fmt.Sprintf(mainLapFormat, lap.timeTaken, ""))
			} else {
				_, _ = mainLapString.WriteString(fmt.Sprintf(mainLapFormat, lap.timeTaken, fmt.Sprintf("%.3f", lap.avgSpeed)))
			}
			if i != len(competitor.mainLaps)-1 {
				_, _ = mainLapString.WriteString(", ")
			}
		}

		if !math.IsNaN(competitor.penaltyLaps.avgSpeed) {
			penaltySpeed = fmt.Sprintf("%.3f", competitor.penaltyLaps.avgSpeed)
		}

		dataToOrder[i] = fmt.Sprintf(outputFormat,
			competitor.totalTime,
			competitor.competitorID,
			mainLapString.String(),
			competitor.penaltyLaps.timeTaken,
			penaltySpeed,
			competitor.hits,
			competitor.shots)

	}

	sort.Slice(dataToOrder, func(i, j int) bool {
		return dataToOrder[i] < dataToOrder[j]
	})
	for _, s := range dataToOrder {
		_, _ = output.Write([]byte(s))
	}
}
