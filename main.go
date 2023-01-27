package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadmahdi255/http-monitor/database"
	"github.com/mohammadmahdi255/http-monitor/handler"
	"github.com/mohammadmahdi255/http-monitor/monitor"
	"log"
	"time"
)

func main() {
	db := database.New("http-monitor.db")

	h := handler.New(db)

	mnt := monitor.NewMonitor(h, nil, 10)
	sch, _ := monitor.NewScheduler(mnt)
	sch.DoWithIntervals(time.Second * 5)

	err := mnt.LoadFromDatabase()
	if err != nil {
		log.Println(err)
	}

	// init echo
	e := echo.New()
	rg := e.Group("/api")
	h.RegisterRoutes(rg)

	e.Logger.Fatal(e.Start("127.0.0.1:5000"))
}
