package cmscontroller

import (
	"app/internal/controllers/cmscontroller/internal/binds"
	"app/internal/controllers/cmscontroller/internal/render"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt *Controller) listBotHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		bots, err := cnt.useCases.GetBots(c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.ConvertSliceWithAloc(bots, render.BotFromDomain))
	}
}

func (cnt *Controller) createBotHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.CreateBotRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.CreateBot(c.Request().Context(), render.BotToDomain(render.Bot(req)))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (cnt *Controller) updateBotHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.UpdateBotRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.UpdateBot(c.Request().Context(), render.BotToDomain(render.Bot(req)))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (cnt *Controller) deleteBotHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.DeleteBotRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.DeleteBot(c.Request().Context(), req.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (cnt *Controller) getBotHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.GetBotRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		raw, err := cnt.useCases.GetBot(c.Request().Context(), req.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.BotFromDomain(raw))
	}
}
