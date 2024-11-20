package cmscontroller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

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

func (cnt *Controller) ffMediaHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var fileBody io.Reader

		actionType := c.FormValue("type")

		if actionType != "text" {
			file, err := c.FormFile("file-data")
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			fileData, err := file.Open()
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			fileBody = fileData

			defer fileData.Close()
		}

		chatID, err := strconv.ParseInt(c.FormValue("chat-id"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		switch actionType {
		case "audio":
			err = cnt.botController.SendAudio(c.Request().Context(), chatID, c.FormValue("filename"), fileBody)
		case "video":
			err = cnt.botController.SendVideo(c.Request().Context(), chatID, c.FormValue("filename"), fileBody)
		case "image":
			err = cnt.botController.SendImage(c.Request().Context(), chatID, c.FormValue("filename"), fileBody)
		case "text":
			err = cnt.botController.SendText(c.Request().Context(), chatID, c.FormValue("filename"))
		default:
			return c.String(http.StatusBadRequest, "unknown type")
		}

		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
