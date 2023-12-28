package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/domain"
)

func (h handlers) GetIncidentByID(c echo.Context, incidentID string) error {
	_, err := retrieveUserFromCtx(c)
	if err != nil {
		return NewUnableToRetrieveUserFromCtx(err)
	}

	iID, err := domain.NewIdFromString(incidentID)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "invalid-incident-id", http.StatusBadRequest)
	}

	foundIncident, err := h.application.Queries.IncidentByID.Execute(c.Request().Context(), query.IncidentByID{
		ID: iID,
	})
	if err != nil {
		return NewHandlerError(err, "unable-to-get-incident")
	}

	return c.JSON(http.StatusOK, GetIncidentByByIdPayload{
		Data: mapIncidentToFullIncidentResponse(foundIncident),
	})
}
