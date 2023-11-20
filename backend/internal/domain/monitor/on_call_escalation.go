package monitor

import (
	"fmt"

	"github.com/friendsofgo/errors"

	"github.com/flowck/doberman/internal/domain"
)

type OnCallEscalation struct {
	notificationMethods     []NotificationMethod
	teamMembersToBeNotified []domain.ID
}

func NewOnCallEscalation(notificationMethods []NotificationMethod, teamMembersToBeNotified []domain.ID) (OnCallEscalation, error) {
	if len(notificationMethods) == 0 {
		return OnCallEscalation{}, errors.New("notificationMethods cannot be nil or empty")
	}

	if len(teamMembersToBeNotified) == 0 {
		return OnCallEscalation{}, errors.New("teamMembersToBeNotified cannot be nil or empty")
	}

	return OnCallEscalation{
		notificationMethods:     notificationMethods,
		teamMembersToBeNotified: teamMembersToBeNotified,
	}, nil
}

func (e OnCallEscalation) IsValid() bool {
	return len(e.notificationMethods) > 0 && len(e.teamMembersToBeNotified) > 0
}

var (
	NotificationMethodSMS   = NotificationMethod{"sms"}
	NotificationMethodEmail = NotificationMethod{"email"}
	NotificationMethodCall  = NotificationMethod{"call"}
)

type NotificationMethod struct {
	value string
}

func NewNotificationMethod(value string) (NotificationMethod, error) {
	switch value {
	case NotificationMethodEmail.value:
		return NotificationMethodEmail, nil
	case NotificationMethodSMS.value:
		return NotificationMethodSMS, nil
	case NotificationMethodCall.value:
		return NotificationMethodSMS, nil
	default:
		return NotificationMethod{}, fmt.Errorf("%s is an unknown notification method", value)
	}
}
