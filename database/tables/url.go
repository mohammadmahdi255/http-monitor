package tables

import (
	"net/http"
	"time"
)

type URL struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UserId      uint   `gorm:"uniqueIndex:index_addr_user"` // for preventing url duplication for a single user
	Address     string `gorm:"uniqueIndex:index_addr_user"`
	Threshold   int
	FailedTimes int
	Requests    []Request `gorm:"foreignkey:url_id"`
}

// SendRequest sends a HTTP GET request to the url
// returns a *Request with result status code
func (url *URL) SendRequest() (*Request, error) {
	resp, err := http.Get(url.Address)
	req := new(Request)
	req.UrlId = url.ID
	if err != nil {
		return req, err
	}
	req.Result = resp.StatusCode
	return req, nil
}
