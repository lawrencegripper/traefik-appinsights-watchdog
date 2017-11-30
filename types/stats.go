package types

import "time"

// StatsEvent is used to track an event and its data
// these are produced by multiple different providers
type StatsEvent struct {
	Source          string
	SourceTime      time.Time
	RequestDuration time.Duration
	IsSuccess       bool
	ErrorDetails    string
	Data            map[string]interface{}
}
