package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/http-monitor/common"
	"github.com/mohammadmahdi255/http-monitor/database/tables"
	"net/http"
	"time"
)

type urlCreateRequest struct {
	Address   string `json:"address" valid:"url"`
	Threshold int    `json:"threshold" valid:"int"`
}

func (r *urlCreateRequest) bind(c echo.Context, url *tables.URL) error {
	if err := c.Bind(r); err != nil {
		return common.NewRequestError("error binding url create request, check json structure and try again", err, http.StatusBadRequest)
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "Error validating create url request")
		return e
	}
	url.Address = r.Address
	url.Threshold = r.Threshold
	url.FailedTimes = 0
	return nil
}

type urlStatusRequest struct {
	FromTime int64 `valid:"optional, time~Provide time as unix timestamp before current time" json:"from_time, omitempty" query:"from_time"`
	ToTime   int64 `valid:"optional, time~Provide time as unix timestamp before current time" json:"to_time, omitempty" query:"to_time"`
}

func (r *urlStatusRequest) parse(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return common.NewRequestError("error parsing url status request, if you want to specify time, use unix timestamp", err, http.StatusBadRequest)
	}
	govalidator.CustomTypeTagMap.Set("time", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
		if _, ok := context.(urlStatusRequest); ok {
			if t, ok := i.(int64); ok {
				if time.Now().Unix() > t {
					return true
				}
			}
		}
		return false
	}))
	if _, err := govalidator.ValidateStruct(r); err != nil {
		e := common.NewValidationError(err, "error validating url status request")
		return e
	}
	if r.FromTime > r.ToTime && r.ToTime != 0 {
		return common.NewRequestError("end of time interval must be later than it's start", nil, http.StatusBadRequest)
	}
	return nil
}
