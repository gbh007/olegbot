package cmscontroller

import (
	"app/internal/controllers/cmscontroller/internal/binds"
	"app/internal/controllers/cmscontroller/internal/render"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt *Controller) quoteListHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.QuoteListRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		rawQuotes, err := cnt.useCases.Quotes(c.Request().Context(), req.BotID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.ConvertSliceWithAlloc(rawQuotes, render.QuoteFromDomain))
	}
}

func (cnt *Controller) quoteGetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.GetQuoteRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		rawQuote, err := cnt.useCases.Quote(c.Request().Context(), req.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.QuoteFromDomain(rawQuote))
	}
}

func (cnt *Controller) deleteQuoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.DeleteQuoteRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.DeleteQuote(c.Request().Context(), req.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (cnt *Controller) updateQuoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.UpdateQuoteRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.UpdateQuoteText(c.Request().Context(), req.ID, req.Text)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (cnt *Controller) createQuoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := binds.AddQuoteRequest{}

		err := c.Bind(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.AddQuote(c.Request().Context(), req.BotID, req.Text)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
