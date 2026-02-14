package clock

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now().UTC()
}

type MockClock struct {
	Current time.Time
}

func (m MockClock) Now() time.Time {
	return m.Current
}
