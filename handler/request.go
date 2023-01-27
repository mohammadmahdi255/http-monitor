package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/http-monitor/common"
	"github.com/mohammadmahdi255/http-monitor/model"
	"net/http"
)

type userAuthRequest struct {
	Username string `valid:"stringlength(4|32), alphanum" json:"username"`
	Password string `valid:"stringlength(4|32)" json:"password"`
}

// binding user auth request with model.User instance
func (u *userAuthRequest) bind(c echo.Context, user *model.User) error {
	if err := c.Bind(u); err != nil {
		return common.MakeRequestError("error binding user request", err, http.StatusBadRequest)
	}
	if _, err := govalidator.ValidateStruct(u); err != nil {
		e := common.NewValidationError(err, "Error validating sign-up request")
		return e
	}
	user.Username = u.Username
	user.Password = u.Password
	return nil
}
