package monitor

import (
	"time"

	"github.com/flowck/doberman/internal/domain"
)

var (
	EscalationPolicyRepetitionRuleNone EscalationPolicyRepetitionRule = 0
	EscalationUrgencyLowUrgency                                       = EscalationUrgency{"low_urgency"}
	EscalationUrgencyHighUrgency                                      = EscalationUrgency{"high_urgency"}
)

type EscalationPolicy struct {
	name                           string
	waitPeriodBeforeEscalation     time.Duration
	peopleToBeNotified             []domain.ID
	servicesToBeNotified           []Integration
	urgency                        EscalationUrgency
	repeatPolicyOnNoAcknowledgment EscalationPolicyRepetitionRule
}

func (e EscalationPolicy) Name() string {
	return e.name
}

func (e EscalationPolicy) WaitPeriodBeforeEscalation() time.Duration {
	return e.waitPeriodBeforeEscalation
}

func (e EscalationPolicy) PeopleToBeNotified() []domain.ID {
	return e.peopleToBeNotified
}

func (e EscalationPolicy) ServicesToBeNotified() []Integration {
	return e.servicesToBeNotified
}

func (e EscalationPolicy) Urgency() EscalationUrgency {
	return e.urgency
}

func (e EscalationPolicy) RepeatPolicyOnNoAcknowledgment() EscalationPolicyRepetitionRule {
	return e.repeatPolicyOnNoAcknowledgment
}

type EscalationPolicyRepetitionRule int

type EscalationUrgency struct {
	value string
}
