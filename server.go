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

var checkingService *checkin.Service

func main() {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	checkingService = &checkin.Service{DB: db}
	fmt.Printf("Setup... \n", checkingService)
	// BestBy 12/2024 - this shouldn't be necessary and definitely not live in main()
	createTableStatement := `CREATE TABLE IF NOT EXISTS check_in (id INTEGER PRIMARY KEY, user_id INTEGER, timestamp INTEGER, date TEXT, payload TEXT)`
	_, err = db.Exec(createTableStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, createTableStatement)
		return
	}

	// temp sample data
	addSampleData()

	e := echo.New()

	e.GET("/", MainPage)

	e.GET("/check-in", GetCheckinHandler)
	e.POST("/check-in", PostCheckinHandler)
	e.DELETE("/check-in/:id", DeleteCheckinHandler)

	e.Logger.Fatal(e.Start(":8001"))
}

func addSampleData() error {
	// load sample data into the database
	err := checkingService.AddCheckIn(checkin.AddCheckInRequest{UserID: 1, Timestamp: "2021-01-01", Date: "2021-01-01", Payload: "payload"})
	if err != nil {
		return err
	}

	return nil
}

func GetCheckinHandler(c echo.Context) error {
	userID := c.QueryParam("user_id")

	r, err := checkingService.GetCheckInCount(userID)
	if err != nil {
		fmt.Errorf("error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, r)
}

func PostCheckinHandler(c echo.Context) error {
	var req checkin.AddCheckInRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err := checkingService.AddCheckIn(req)
	if err != nil {
		fmt.Errorf("error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func DeleteCheckinHandler(c echo.Context) error {
	id := c.Param("id")

	err := checkingService.DeleteCheckIn(id)
	if err != nil {
		fmt.Errorf("error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

// BestBy 01/2025 - have a better default route.
func MainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
