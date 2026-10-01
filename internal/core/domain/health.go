package domain

import "time"

type Status string

const (
	StatusUP       Status = "UP"
	StatusDOWN     Status = "DOWN"
	StatusDEGRADED Status = "DEGRADED"
)

type Health struct {
	Status    Status
	Timestamp time.Time
	Version   string
}
