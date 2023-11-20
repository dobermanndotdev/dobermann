package monitor

import (
	"net/http"
	"time"

	"github.com/friendsofgo/errors"

	"github.com/flowck/doberman/internal/domain"
)

type Monitor struct {
	id                     domain.ID
	accountID              domain.ID
	teamID                 domain.ID
	endpoint               Endpoint
	advancedSettings       *Settings
	sslVerification        *SslVerification
	requestParameters      RequestParameters
	requestHeaders         RequestHeaders
	httpAuth               HttpAuth
	expectedResponseStatus ResponseStatus
	maintenance            Maintenance
	regions                []Region
	alertTriggers          []AlertTrigger
	onCallEscalation       OnCallEscalation
	isPaused               bool
	isUp                   bool
	checkStatus            CheckStatus
	lastCheckedAt          time.Time
}

func New(
	id,
	accountID,
	teamID domain.ID,
	alertTriggers []AlertTrigger,
	endpoint Endpoint,
	onCallEscalation OnCallEscalation,
) (*Monitor, error) {
	if id.IsEmpty() {
		return nil, errors.New("id cannot be invalid")
	}

	if accountID.IsEmpty() {
		return nil, errors.New("accountID cannot be invalid")
	}

	if teamID.IsEmpty() {
		return nil, errors.New("teamID cannot be invalid")
	}

	if len(alertTriggers) == 0 {
		return nil, errors.New("alertTriggers cannot be nil or empty")
	}

	if !endpoint.IsValid() {
		return nil, errors.New("endpoint cannot be invalid")
	}

	if !onCallEscalation.IsValid() {
		return nil, errors.New("onCallEscalation cannot be invalid")
	}

	return &Monitor{
		id:               id,
		teamID:           teamID,
		accountID:        accountID,
		alertTriggers:    alertTriggers,
		onCallEscalation: onCallEscalation,
		endpoint:         endpoint,
		checkStatus:      CheckStatusPending,
		advancedSettings: newDefaultSettings(),
		sslVerification: &SslVerification{
			isEnabled:          false,
			verifyExpirationIn: 0,
		},
		requestParameters:      newDefaultRequestParameters(),
		requestHeaders:         nil,
		httpAuth:               HttpAuth{},
		expectedResponseStatus: http.StatusOK,
		maintenance:            Maintenance{},
		regions:                NewDefaultRegions(),
		lastCheckedAt:          time.Now(),
	}, nil
}

func NewFrom(
	id domain.ID,
	accountID domain.ID,
	teamID domain.ID,
	endpoint Endpoint,
	advancedSettings *Settings,
	sslVerification *SslVerification,
	requestParameters RequestParameters,
	requestHeaders RequestHeaders,
	httpAuth HttpAuth,
	expectedResponseStatus ResponseStatus,
	maintenance Maintenance,
	regions []Region,
	alertTriggers []AlertTrigger,
	onCallEscalation OnCallEscalation,
	isPaused bool,
	isUp bool,
	checkStatus CheckStatus,
	lastCheckedAt time.Time,
) (*Monitor, error) {
	return &Monitor{
		id:                     id,
		teamID:                 teamID,
		accountID:              accountID,
		alertTriggers:          alertTriggers,
		onCallEscalation:       onCallEscalation,
		endpoint:               endpoint,
		isUp:                   isUp,
		isPaused:               isPaused,
		checkStatus:            checkStatus,
		advancedSettings:       advancedSettings,
		sslVerification:        sslVerification,
		requestParameters:      requestParameters,
		requestHeaders:         requestHeaders,
		httpAuth:               httpAuth,
		expectedResponseStatus: expectedResponseStatus,
		maintenance:            maintenance,
		regions:                regions,
		lastCheckedAt:          lastCheckedAt,
	}, nil
}

func (m *Monitor) Enqueue() {
	m.checkStatus = CheckStatusEnqueued
}

func (m *Monitor) Pause() {
	m.isPaused = true
}

func (m *Monitor) UnPause() {
	m.isPaused = false
}

func (m *Monitor) SetStatusAsDown() {
	m.isUp = false
}

func (m *Monitor) CheckStatus() CheckStatus {
	return m.checkStatus
}

func (m *Monitor) AccountID() domain.ID {
	return m.accountID
}

func (m *Monitor) ID() domain.ID {
	return m.id
}

func (m *Monitor) TeamID() domain.ID {
	return m.teamID
}

func (m *Monitor) AdvancedSettings() *Settings {
	return m.advancedSettings
}

func (m *Monitor) SslVerification() *SslVerification {
	return m.sslVerification
}

func (m *Monitor) RequestParameters() RequestParameters {
	return m.requestParameters
}

func (m *Monitor) RequestHeaders() RequestHeaders {
	return m.requestHeaders
}

func (m *Monitor) HttpAuth() HttpAuth {
	return m.httpAuth
}

func (m *Monitor) ExpectedResponseStatus() ResponseStatus {
	return m.expectedResponseStatus
}

func (m *Monitor) Maintenance() Maintenance {
	return m.maintenance
}

func (m *Monitor) Regions() []Region {
	return m.regions
}

func (m *Monitor) AlertTriggers() []AlertTrigger {
	return m.alertTriggers
}

func (m *Monitor) OnCallEscalation() OnCallEscalation {
	return m.onCallEscalation
}

func (m *Monitor) IsPaused() bool {
	return m.isPaused
}

func (m *Monitor) IsUp() bool {
	return m.isUp
}

func (m *Monitor) Endpoint() Endpoint {
	return m.endpoint
}

func (m *Monitor) LastCheckedAt() time.Time {
	return m.lastCheckedAt
}

type IncidentCount int

func (i IncidentCount) IsValid() bool {
	return i >= 0
}

type ResponseStatus int

func (r ResponseStatus) IsValid() bool {
	return r >= 100 && r <= 599
}

func (r ResponseStatus) Int() int {
	return int(r)
}

type Integration struct{}
