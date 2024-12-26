package cmscontroller

import (
	"app/internal/controllers/cmscontroller/internal/binds"
	"app/internal/controllers/cmscontroller/internal/render"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt *Controller) moderatorsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.ListModeratorRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		rawModerators, err := cnt.useCases.Moderators(c.Request().Context(), req.BotID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.ConvertSliceWithAlloc(rawModerators, render.ModeratorFromDomain))	}
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

		err = cnt.useCases.AddModerator(c.Request().Context(), req.BotID, req.UserID, req.Description)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (cnt *Controller) deleteModeratorHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.DeleteModeratorRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.DeleteModerator(c.Request().Context(), req.BotID, req.UserID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
