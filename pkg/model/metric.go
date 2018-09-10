package model

import (
	"fmt"
	"time"
)

// Metric DTO
type Metric struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Error     string        `json:"error,omitempty"`
}

func (m Metric) String() string {
	if m.Error != "" {
		return fmt.Sprintf(
			"{status: \"%s\", error: \"%s\", duration: %d, ts: %s}",
			m.Status,
			m.Error,
			m.Duration,
			m.Timestamp)
	}
	return fmt.Sprintf(
		"{status: \"%s\", duration: %d, ts: %s}",
		m.Status,
		m.Duration,
		m.Timestamp)
}
