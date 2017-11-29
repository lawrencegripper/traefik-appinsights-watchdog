package types

import "time"

type StatsEvent struct {
	SourceTime      time.Time
	RequestDuration time.Duration
	IsSuccess       bool
	ErrorDetails    string
	Data            map[string]interface{}
}
