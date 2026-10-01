package health

import "time"

type Output struct {
	Status    string
	Timestamp time.Time
	Version   string
}

type Input struct {
	Version string
}
