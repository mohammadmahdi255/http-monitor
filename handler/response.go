package handler

import (
	"github.com/mohammadmahdi255/http-monitor/common"
	"github.com/mohammadmahdi255/http-monitor/database/tables"
	"time"
)

type UserResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func NewUserResponse(user *tables.User) *UserResponse {
	token, _ := common.GenerateJWT(user.ID)
	ur := &UserResponse{Username: user.Username, Token: token}
	return ur
}

type urlResponse struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Threshold   int       `json:"threshold"`
	FailedTimes int       `json:"failed_times"`
}

func newURLResponse(url *tables.URL) *urlResponse {
	u := new(urlResponse)
	u.URL = url.Address
	u.UserID = url.UserId
	u.CreatedAt = url.CreatedAt
	u.Threshold = url.Threshold
	u.FailedTimes = url.FailedTimes
	return u
}

type urlListResponse struct {
	URLs     []*urlResponse `json:"urls"`
	UrlCount int            `json:"url_count"`
}

func newURLListResponse(list []tables.URL) *urlListResponse {
	resp := new(urlListResponse)
	resp.URLs = make([]*urlResponse, 0)
	for i := range list {
		resp.URLs = append(resp.URLs, newURLResponse(&list[i]))
	}
	resp.UrlCount = len(list)
	return resp
}

type requestResponse struct {
	ResultCode int       `json:"result_code"`
	CreatedAt  time.Time `json:"created_at"`
}

func newRequestResponse(req *tables.Request) *requestResponse {
	return &requestResponse{ResultCode: req.Result, CreatedAt: req.CreatedAt}
}

type requestListResponse struct {
	URL           string             `json:"url"`
	RequestsCount int                `json:"requests_count"`
	Requests      []*requestResponse `json:"requests"`
}

// TODO update request struct to have a field for url instance
func newRequestListResponse(reqs []tables.Request, url string) *requestListResponse {
	resp := new(requestListResponse)
	resp.Requests = make([]*requestResponse, len(reqs))
	for i := range reqs {
		resp.Requests[i] = newRequestResponse(&reqs[i])
	}
	resp.URL = url
	resp.RequestsCount = len(reqs)
	return resp
}
