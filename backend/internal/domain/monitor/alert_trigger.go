package monitor

import (
	"fmt"
)

var (
	AlertTriggerUrlIsUnavailable          = AlertTrigger{"url_is_unavailable"}
	AlertTriggerResponseStatusIsDifferent = AlertTrigger{"response_status_is_different"}
)

type AlertTrigger struct {
	name string
}

func NewAlertTrigger(value string) (AlertTrigger, error) {
	switch value {
	case AlertTriggerUrlIsUnavailable.name:
		return AlertTriggerUrlIsUnavailable, nil
	case AlertTriggerResponseStatusIsDifferent.name:
		return AlertTriggerResponseStatusIsDifferent, nil
	default:
		return AlertTrigger{}, fmt.Errorf("%s is an unknown alertTrigger", value)
	}
}
