package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/http-monitor/common"
	"github.com/mohammadmahdi255/http-monitor/database/tables"
	"net/http"
	"strconv"
	"time"
)

// CreateURL is used to add a url to monitor service
// urls are validated and if there isn't any error a response code 201 is returned
// json request format:
//
//	{
//		"address": "http://google.com",
//		"threshold": 20
//	}
func (h *Handler) CreateURL(c echo.Context) error {
	userID, err := extractID(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	req := &urlCreateRequest{}
	url := &tables.URL{}

	if err := req.bind(c, url); err != nil {
		return err
	}
	url.UserId = *userID

	// adding url to database
	if err := h.db.Create(url).Error; err != nil {
		// internal error
		return common.NewRequestError("error adding url to database", err, http.StatusInternalServerError)
	}

	// adding url to monitor scheduler
	//h.sch.Mnt.AddURL([]model.URL{*url})

	return c.JSON(http.StatusCreated, "URL created successfully")
}

// FetchURLs TODO: add pagination support
// FetchURLs is used to retrieve a user's urls
// accessible with GET /api/urls
func (h *Handler) FetchURLs(c echo.Context) error {
	userID, err := extractID(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	urls, err := h.GetURLsByUser(*userID)
	if err != nil {
		return common.NewRequestError("Error retrieving urls from database, maybe check your token again", err, http.StatusBadRequest)
	}
	resp := newURLListResponse(urls)
	return c.JSON(http.StatusOK, resp)
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
	userID, err := extractID(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	urlID, err := strconv.Atoi(c.Param("urlID"))
	if err != nil {
		return common.NewRequestError("Invalid path parameter", err, http.StatusBadRequest)
	}

	req := &urlStatusRequest{}
	url := new(tables.URL)
	if err := req.parse(c); err != nil {
		return err
	}

	from := time.Unix(time.Now().Unix()-86400, 0)
	to := time.Unix(time.Now().Unix(), 0)
	url, err = h.GetUserRequestsInPeriod(uint(urlID), from, to)
	if err != nil {
		e := common.NewRequestError("error retrieving url stats, invalid url id", err, http.StatusBadRequest)
		return c.JSON(e.Status, e)
	}
	if url.UserId != *userID {
		return common.NewRequestError("operation not permitted", errors.New("user is not the owner of url"), http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, newRequestListResponse(url.Requests, url.Address))
}
