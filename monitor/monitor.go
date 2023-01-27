package monitor

import (
	"errors"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/mohammadmahdi255/http-monitor/database/tables"
	"github.com/mohammadmahdi255/http-monitor/handler"
	"net/http"
	"sync"
)

type Monitor struct {
	h          *handler.Handler
	URLs       []tables.URL
	wp         *workerpool.WorkerPool
	workerSize int
}

// NewMonitor creates a Monitor instance with 'store' and 'url'
// it also creates a worker pool of size 'workerSize'
// if 'urls' is set to nil it will be initialized with an empty slice
func NewMonitor(h *handler.Handler, urls []tables.URL, workerSize int) *Monitor {
	mnt := new(Monitor)
	if urls == nil {
		mnt.URLs = make([]tables.URL, 0)
	}
	mnt.URLs = urls
	mnt.h = h
	mnt.workerSize = workerSize
	// max number of workers
	mnt.wp = workerpool.New(workerSize)
	return mnt
}

// AddURL appends a slice of urls to the current list of urls
func (mnt *Monitor) AddURL(urls []tables.URL) {
	mnt.URLs = append(mnt.URLs, urls...)
}

// LoadFromDatabase loads all urls from database into monitor to start working on them
// this function will replace all of saved URLs with the ones from database
func (mnt *Monitor) LoadFromDatabase() error {
	urls, err := mnt.h.GetAllURLs()
	if err != nil {
		return err
	}
	mnt.URLs = urls
	return nil
}

// RemoveURL removes a URL from current list of monitor's urls
// returns error if the URL to be deleted was not found
func (mnt *Monitor) RemoveURL(url tables.URL) error {
	var index = -1
	for i := range mnt.URLs {
		if mnt.URLs[i].ID == url.ID {
			index = i
		}
	}
	if index == -1 {
		return errors.New("url to be deleted was not found in the slice")
	}
	// deleting from list efficiently
	mnt.URLs[index], mnt.URLs[len(mnt.URLs)-1] = mnt.URLs[len(mnt.URLs)-1], mnt.URLs[index]
	mnt.URLs = mnt.URLs[:len(mnt.URLs)-1]
	return nil
}

// Cancel stops all tasks of fetching urls
// it will wait for current running jobs to finish
// note that if you call this method, for reusing the monitor
// you need to instantiate it again.
func (mnt *Monitor) Cancel() error {
	mnt.wp.Stop()
	if !mnt.wp.Stopped() {
		return errors.New("could not stop monitor")
	}
	return nil
}

// DoURL checks a single URLs response and saves its request into database
func (mnt *Monitor) DoURL(url tables.URL) {
	var wg sync.WaitGroup
	wg.Add(1)
	mnt.wp.Submit(func() {
		defer wg.Done()
		mnt.monitorURL(url)
	})
	wg.Wait()
}

// Do ranges over URLs currently inside Monitor instance
// and save each one's request inside database
// this function does not block
func (mnt *Monitor) Do() {
	var wg sync.WaitGroup

	for urlIndex := range mnt.URLs {
		url := mnt.URLs[urlIndex]
		wg.Add(1)
		mnt.wp.Submit(func() {
			defer wg.Done()
			mnt.monitorURL(url)
		})
	}
	wg.Wait()
}

func (mnt *Monitor) monitorURL(url tables.URL) {
	// sending request
	req, err := url.SendRequest()
	if err != nil {
		fmt.Println(err, "could not make request")
		req = new(tables.Request)
		req.UrlId = url.ID
		req.Result = http.StatusBadRequest
	}
	// add request to database
	if err = mnt.h.AddRequest(req); err != nil {
		fmt.Println(err, "could not save request to database")
	}
	// status code was other than 2XX
	if req.Result/100 != 2 {
		if err = mnt.h.IncrementFailed(&url); err != nil {
			fmt.Println(err, "could not increment failed times for url")
		}
	}
}
