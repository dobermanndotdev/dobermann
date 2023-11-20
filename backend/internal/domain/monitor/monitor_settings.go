package monitor

import "time"

const (
	startAnIncidentImmediately = time.Duration(0)
	dontAlertWhenDomainExpires = -1
)

type Settings struct {
	recoveredOnlyAfter          time.Duration
	startAnIncidentAfter        time.Duration
	checkInterval               time.Duration
	alertDomainExpirationWithin time.Duration
}

func NewSettings(
	recoveredOnlyAfter time.Duration,
	startAnIncidentAfter time.Duration,
	checkInterval time.Duration,
	alertDomainExpirationWithin time.Duration) (*Settings, error) {
	return &Settings{
		recoveredOnlyAfter:          time.Minute * 3,
		checkInterval:               time.Minute * 3,
		startAnIncidentAfter:        startAnIncidentImmediately,
		alertDomainExpirationWithin: dontAlertWhenDomainExpires,
	}, nil
}

func newDefaultSettings() *Settings {
	return &Settings{
		recoveredOnlyAfter:          time.Minute * 3,
		checkInterval:               time.Minute * 3,
		startAnIncidentAfter:        startAnIncidentImmediately,
		alertDomainExpirationWithin: dontAlertWhenDomainExpires,
	}
}

func (s *Settings) RecoveredOnlyAfter() time.Duration {
	return s.recoveredOnlyAfter
}

func (s *Settings) StartAnIncidentAfter() time.Duration {
	return s.startAnIncidentAfter
}

func (s *Settings) CheckInterval() time.Duration {
	return s.checkInterval
}

func (s *Settings) AlertDomainExpirationWithin() time.Duration {
	return s.alertDomainExpirationWithin
}
