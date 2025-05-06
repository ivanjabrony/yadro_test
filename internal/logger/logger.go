package logger

import "fmt"

type MyLogger struct{}

func (MyLogger) LogEvent(log string, params ...any) {
	fmt.Printf(log, params...)
}
