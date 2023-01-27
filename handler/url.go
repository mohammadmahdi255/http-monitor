package handler

import (
	"github.com/labstack/echo/v4"
)

// FetchURLs TODO: add pagination support
// FetchURLs is used to retrieve a user's urls
// accessible with GET /api/urls
func (h *Handler) FetchURLs(c echo.Context) error {

}

// CreateURL is used to add a url to monitor service
// urls are validated and if there isn't any error a response code 201 is returned
// json request format:
//
//	{
//		"address": "http://google.com",
//		"threshold": 20
//	}
func (h *Handler) CreateURL(c echo.Context) error {

}

// GetURLStats reports stats of a url
// returns error in case of invalid url_id or unauthenticated request
// param request format :
//
// /api/urls/:urlID
// you can also specify time intervals to get stats in
// just use unix timestamp with the syntax below (to_time is optional):
// /api/urls/:urlID?from_time=1579184689[&to_time]
func (h *Handler) GetURLStats(c echo.Context) error {

}
