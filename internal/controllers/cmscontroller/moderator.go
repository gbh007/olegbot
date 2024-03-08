package cmscontroller

import (
	"app/internal/controllers/cmscontroller/internal/binds"
	"app/internal/controllers/cmscontroller/internal/render"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (cnt *Controller) moderatorsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		rawModerators, err := cnt.useCases.Moderators(c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.ConvertSliceWithAloc(rawModerators, render.ModeratorFromDomain))
	}
}

func (cnt *Controller) addModeratorHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.AddModeratorRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.AddModerator(c.Request().Context(), req.UserID, req.Description)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (cnt *Controller) deleteModeratorHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.DeleteModerator(c.Request().Context(), id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
