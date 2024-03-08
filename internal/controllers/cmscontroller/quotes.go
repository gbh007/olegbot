package cmscontroller

import (
	"app/internal/controllers/cmscontroller/internal/render"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt *Controller) quotesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		rawQuotes, err := cnt.useCases.Quotes(c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, render.ConvertSliceWithAloc(rawQuotes, render.QuoteFromDomain))
	}
}
