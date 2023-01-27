// Package middleware codes from https://github.com/labstack/echo/blob/master/middleware/jwt.go
package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RequestInfo struct {
	path   string
	method string
}

// authWhiteList specifies paths to be skipped by jwt authentication middleware
var authWhiteList = make([]RequestInfo, 0)

// AddToWhiteList is used to add a path to skipper white list
// provide path relative to api version like /api/your/path/here as skipper uses strings.Contains to find whether
// it is in context path or not
func AddToWhiteList(path string, method string) {
	authWhiteList = append(authWhiteList, RequestInfo{path, method})
}

func skipper(c echo.Context) bool {
	path := c.Path()
	method := c.Request().Method
	for _, v := range authWhiteList {
		if path == v.path && method == v.method {
			return true
		}
	}
	return false
}

func JWT(key interface{}) echo.MiddlewareFunc {
	c := middleware.DefaultJWTConfig
	c.SigningKey = key
	c.Skipper = skipper
	return middleware.JWTWithConfig(c)
}
