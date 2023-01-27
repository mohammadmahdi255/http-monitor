package handler

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/http-monitor/database/tables"
	"gorm.io/gorm"
	"time"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}

// FetchAlertsDB FetchAlerts retrieves urls which "failed_times" is greater than it's "threshold" for given userID
// TODO: write tests for this function
func (h *Handler) FetchAlertsDB(userID uint) ([]tables.URL, error) {
	var urls []tables.URL
	err := h.db.Model(&tables.URL{}).Where("id == ? and failed_times >= threshold", userID).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (h *Handler) GetURLsByUser(userID uint) ([]tables.URL, error) {
	var urls []tables.URL
	if err := h.db.Model(&tables.URL{}).Where("user_id == ?", userID).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}

func extractID(c echo.Context) (*uint, error) {
	token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return nil, errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return nil, errors.New("failed to cast claims as jwt.MapClaims")
	}
	id := uint(claims["id"].(float64))
	return &id, nil
}

// GetUserRequestsInPeriod retrieves requests between 2 time intervals
func (h *Handler) GetUserRequestsInPeriod(urlID uint, from, to time.Time) (*tables.URL, error) {
	url := &tables.URL{}
	url.ID = urlID
	if err := h.db.Model(url).Preload("Requests", "created_at >= ? and created_at <= ?", from, to).First(url).Error; err != nil {
		return nil, err
	}
	return url, nil
}

// IncrementFailed increments failed_times of a URL
func (h *Handler) IncrementFailed(url *tables.URL) error {
	url.FailedTimes += 1
	return h.db.Model(url).Updates(url).Error
}

// AddRequest adds a request to database
func (h *Handler) AddRequest(req *tables.Request) error {
	return h.db.Create(req).Error
}

func (h *Handler) GetAllURLs() ([]tables.URL, error) {
	var urls []tables.URL
	if err := h.db.Model(&tables.URL{}).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}
