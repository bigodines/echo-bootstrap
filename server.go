package main

import (
	"database/sql"
	"echo-bootstrap/checkin"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	CheckInService := &checkin.Service{DB: db}
	fmt.Printf("Setup... \n", CheckInService)
	// BestBy 12/2024 - this shouldn't be necessary and definitely not live in main()
	createTableStatement := `CREATE TABLE IF NOT EXISTS check_in (id INTEGER PRIMARY KEY, user_id INTEGER, timestamp INTEGER, date TEXT, payload TEXT)`
	_, err = db.Exec(createTableStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, createTableStatement)
		return
	}
	// load sample data into the database
	err = CheckInService.AddCheckIn(checkin.CheckInRequest{UserID: 1, Timestamp: "2021-01-01", Date: "2021-01-01", Payload: "payload"})
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.GET("/", MainPage)
	e.GET("/check-in", GetCheckins(*CheckInService))

	e.POST("/check-in", Inc(*CheckInService))

	e.Logger.Fatal(e.Start(":8001"))
}

func Inc(service checkin.Service) func(context2 echo.Context) error {
	return func(c echo.Context) error {
		// retrieve the value of variable "metric" from request body
		metric := c.FormValue("metric")

		return c.String(http.StatusOK, metric)
	}
}

func GetCheckins(service checkin.Service) func(context2 echo.Context) error {
	return func(c echo.Context) error {
		userID := c.QueryParam("user_id")

		r, err := service.GetCheckInCount(userID)
		if err != nil {
			fmt.Errorf("error: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, r)
	}
}

// BestBy 01/2025 - have a better default route.
func MainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
