package monitors

import (
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/flowck/doberman/internal/adapters/models"
	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/internal/domain/monitor"
)

func mapFromMonitorToModel(mo *monitor.Monitor) *models.Monitor {
	return &models.Monitor{
		ID:                          mo.ID().String(),
		AccountID:                   mo.AccountID().String(),
		TeamID:                      mo.TeamID().String(),
		Endpoint:                    mo.Endpoint().String(),
		RecoveredOnlyAfter:          int16(mo.AdvancedSettings().RecoveredOnlyAfter().Minutes()),
		StartAnIncidentAfter:        int16(mo.AdvancedSettings().StartAnIncidentAfter().Minutes()),
		CheckInterval:               int16(mo.AdvancedSettings().CheckInterval().Minutes()),
		AlertDomainExpirationWithin: int16(mo.AdvancedSettings().AlertDomainExpirationWithin().Minutes()),
		SSLVerificationEnabled:      mo.SslVerification().IsEnabled(),
		VerifySSLExpirationWithin:   null.Int16From(int16(mo.SslVerification().VerifyExpirationIn().Minutes())),
		RequestMethod:               mo.RequestParameters().Method().String(),
		RequestTimeout:              mo.RequestParameters().Timeout(),
		RequestBody:                 null.StringFrom(mo.RequestParameters().Body()),
		FollowRedirects:             mo.RequestParameters().FollowRedirects(),
		KeepCookiesWhileRedirecting: mo.RequestParameters().KeepCookiesWhileRedirecting(),
		ExpectedResponseStatus:      int(mo.ExpectedResponseStatus()),
		BasicAuthUsername:           null.StringFrom(mo.HttpAuth().Username()),
		BasicAuthPassword:           null.StringFrom(mo.HttpAuth().Password()),
		MaintenanceFrom:             null.TimeFrom(mo.Maintenance().From()),
		MaintenanceTo:               null.TimeFrom(mo.Maintenance().To()),
		LastCheckedAt:               mo.LastCheckedAt(),
		CheckStatus:                 mo.CheckStatus().String(),
		IsPaused:                    mo.IsPaused(),
		IsUp:                        mo.IsUp(),
	}
}

func mapFromMonitorsToModels(mos []*monitor.Monitor) []*models.Monitor {
	result := make([]*models.Monitor, len(mos))

	for i, mo := range mos {
		result[i] = mapFromMonitorToModel(mo)
	}

	return result
}

func mapFromModelToMonitor(model *models.Monitor) (*monitor.Monitor, error) {
	ID, err := domain.NewIdFromString(model.ID)
	if err != nil {
		return nil, err
	}

	accountID, err := domain.NewIdFromString(model.AccountID)
	if err != nil {
		return nil, err
	}

	teamID, err := domain.NewIdFromString(model.TeamID)
	if err != nil {
		return nil, err
	}

	endpoint, err := monitor.NewEndpoint(model.Endpoint)
	if err != nil {
		return nil, err
	}

	settings, err := monitor.NewSettings(
		time.Duration(model.RecoveredOnlyAfter),
		time.Duration(model.StartAnIncidentAfter),
		time.Duration(model.CheckInterval),
		time.Duration(model.AlertDomainExpirationWithin),
	)
	if err != nil {
		return nil, err
	}

	requestParameters, err := monitor.NewRequestParameters(model.RequestMethod, model.RequestBody.String, model.RequestTimeout, model.FollowRedirects, model.KeepCookiesWhileRedirecting)
	if err != nil {
		return nil, err
	}

	sslVerification := &monitor.SslVerification{}
	requestHeaders := monitor.RequestHeaders{}
	httpAuth := monitor.HttpAuth{}
	maintenance := monitor.Maintenance{}
	regions := []monitor.Region{}
	alertTriggers := []monitor.AlertTrigger{}
	onCallEscalation := monitor.OnCallEscalation{}
	checkStatus := monitor.CheckStatus{}

	return monitor.NewFrom(
		ID,
		accountID,
		teamID,
		endpoint,
		settings,
		sslVerification,
		requestParameters,
		requestHeaders,
		httpAuth,
		monitor.ResponseStatus(model.ExpectedResponseStatus),
		maintenance,
		regions,
		alertTriggers,
		onCallEscalation,
		model.IsPaused,
		model.IsUp,
		checkStatus,
		model.LastCheckedAt,
	)
}

func mapFromModelsToMonitors(m []*models.Monitor) ([]*monitor.Monitor, error) {
	result := make([]*monitor.Monitor, len(m))

	for i, model := range m {
		mo, err := mapFromModelToMonitor(model)
		if err != nil {
			return nil, err
		}

		result[i] = mo
	}

	return result, nil
}
