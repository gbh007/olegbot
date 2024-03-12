package cmscontroller

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt *Controller) ffQuoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("quotes")
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		fileData, err := file.Open()
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		defer fileData.Close()

		req := make([]string, 0)

		err = json.NewDecoder(fileData).Decode(&req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = cnt.useCases.AddQuotes(c.Request().Context(), req)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
