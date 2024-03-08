package cmscontroller

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (cnt *Controller) newBaseAuth() echo.MiddlewareFunc {
	return middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/metrics" {
				return true
			}

			if cnt.cmsLogin == "" || cnt.cmsPass == "" {
				return true
			}

			return false
		},
		Realm: "olegbotcms",
		Validator: func(s1, s2 string, ctx echo.Context) (bool, error) {
			if cnt.cmsLogin == s1 && cnt.cmsPass == s2 {
				return true, nil
			}

			return false, nil
		},
	})
}
