package cmscontroller

import (
	"app/internal/controllers/cmscontroller/internal/binds"
	"app/internal/controllers/cmscontroller/internal/render"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (cnt *Controller) quotesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		rawQuotes, err := cnt.useCases.Quotes(c.Request().Context(), cnt.botID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.ConvertSliceWithAloc(rawQuotes, render.QuoteFromDomain))
	}
}

func (cnt *Controller) quoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		rawQuote, err := cnt.useCases.Quote(c.Request().Context(), id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.QuoteFromDomain(rawQuote))
	}
}

func (cnt *Controller) deleteQuoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.DeleteQuote(c.Request().Context(), id)
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

func (cnt *Controller) addQuoteHandler() echo.HandlerFunc {
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

		err = cnt.useCases.AddQuote(c.Request().Context(), cnt.botID, req.Text)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
