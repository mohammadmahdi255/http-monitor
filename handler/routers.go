package handler

import (
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mohammadmahdi255/http-monitor/common"
	"github.com/mohammadmahdi255/http-monitor/middleware"
)

// RegisterRoutes registers routes with their corresponding handler function
// functions are defined in handler package
func (h *Handler) RegisterRoutes(rg *echo.Group) {

	rg.Use(middleware.Logger())
	rg.Use(middleware.Recover())

	// adding white list
	mileware.AddToWhiteList("/api/users", "POST")
	mileware.AddToWhiteList("/api/users/login", "POST")

	rg.Use(echojwt.WithConfig(mileware.Config(common.JWTSecret)))

	userGroup := rg.Group("/users")
	userGroup.POST("", h.SignUp)
	userGroup.POST("/login", h.Login)

	alertGroup := rg.Group("/alerts")
	alertGroup.GET("", h.FetchAlerts)

	urlGroup := rg.Group("/urls")
	urlGroup.POST("", h.CreateURL)
	urlGroup.GET("", h.FetchURLs)
	urlGroup.GET("/:urlID", h.GetURLStats)

}
