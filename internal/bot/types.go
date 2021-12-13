package bot

import "time"

// ExecutionStatus contains the status of an executed command or component.
type ExecutionStatus struct {
	Err error

	TimeStarted  time.Time
	TimeFinished time.Time
}

// Duration returns how long it took to execute.
func (s *ExecutionStatus) Duration() time.Duration {
	return s.TimeFinished.Sub(s.TimeStarted)
}
