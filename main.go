package main

import (
	"fmt"
	"github.com/mohammadmahdi255/http-monitor/database"
	"github.com/mohammadmahdi255/http-monitor/database/tables"
)

func main() {
	db := database.New("http-monitor.db")

	db.Create(&tables.User{Username: "mahdi", Password: "lae@110"})

	var p tables.User
	db.First(&p, "Username = ?", "mahdi")

	fmt.Println(p)
	//h := handler.NewHandler(db)

	// init echo
	//e := echo.New()
	//g := e.Group("/api")
	//h.RegisterRoutes(g)
	//
	//e.Logger.Fatal(e.Start("127.0.0.1:5000"))
}
