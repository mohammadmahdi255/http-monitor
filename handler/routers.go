package handler

import (
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mohammadmahdi255/http-monitor/common"
	Middleware "github.com/mohammadmahdi255/http-monitor/middleware"
)

// RegisterRoutes registers routes with their corresponding handler function
// functions are defined in handler package
func (h *Handler) RegisterRoutes(routerGroup *echo.Group) {

	routerGroup.Use(middleware.RemoveTrailingSlash())

	routerGroup.Use(middleware.Logger())
	routerGroup.Use(middleware.Recover())

	routerGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: common.JWTSecret,
	}))

	// adding white list
	Middleware.AddToWhiteList("/api/users", "POST")
	Middleware.AddToWhiteList("/api/users/login", "POST")

	userGroup := routerGroup.Group("/users")
	userGroup.POST("", h.SignUp)
	userGroup.POST("/login", h.Login)

	urlGroup := routerGroup.Group("/urls")
	urlGroup.GET("", h.FetchURLs)
	urlGroup.POST("", h.CreateURL)
	urlGroup.GET("/:urlID", h.GetURLStats)

	alertGroup := routerGroup.Group("/alerts")
	alertGroup.GET("", h.FetchAlerts)
}
