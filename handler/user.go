package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/http-monitor/common"
	"github.com/mohammadmahdi255/http-monitor/database/tables"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	Username string `valid:"stringlength(4|32), alphanum" json:"username"`
	Password string `valid:"stringlength(4|32)" json:"password"`
}

func (h *Handler) SignUp(c echo.Context) error {
	dbUser := &tables.User{}
	user := &User{}
	var err error

	if err = c.Bind(user); err != nil {
		e := common.NewRequestError("error binding user request", err, http.StatusBadRequest)
		return c.JSON(e.Status, e)
	}

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		e := common.NewRequestError("could not Hash user password", err, http.StatusInternalServerError)
		return c.JSON(e.Status, e)
	}

	dbUser.Username = user.Username
	dbUser.Password = user.Password

	// saving user
	if err = h.db.Create(dbUser).Error; err != nil {
		e := common.NewRequestError("could not save user in database", err, http.StatusInternalServerError)
		return c.JSON(e.Status, e)
	}

	return c.JSON(http.StatusCreated, NewUserResponse(dbUser))
}

func (h *Handler) Login(c echo.Context) error {
	dbUser := &tables.User{}
	user := &User{}
	var err error

	if err = c.Bind(user); err != nil {
		e := common.NewRequestError("error binding user request", err, http.StatusBadRequest)
		return c.JSON(e.Status, e)
	}

	err = h.db.Model(&tables.User{}).First(&dbUser, "username = ?", user.Username).Error
	if err != nil || !user.ValidatePassword(dbUser.Password) {
		e := common.NewRequestError("Invalid username or password", err, http.StatusUnauthorized)
		return c.JSON(e.Status, e)
	}

	return c.JSON(http.StatusOK, NewUserResponse(dbUser))
}

// FetchAlerts retrieves all alerts for the user, returns a list of urls with alert
func (h *Handler) FetchAlerts(c echo.Context) error {
	userID, err := extractID(c)

	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	alerts, err := h.FetchAlertsDB(*userID)
	if err != nil {
		e := common.NewRequestError("could not get alerts from database", err, http.StatusBadRequest)
		return c.JSON(e.Status, e)
	}
	return c.JSON(http.StatusOK, alerts)
}

// HashPassword generates a hashed string from 'pass'
// returns error if 'pass' is empty
func HashPassword(pass string) (string, error) {

	if len(pass) == 0 {
		return "", errors.New("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	return string(hash), err
}

// ValidatePassword compares 'pass' with 'users' password
// returns true if their equivalent
func (user *User) ValidatePassword(hashPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(user.Password)) == nil
}
