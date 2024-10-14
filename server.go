package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type (
	Service struct {
		db *sql.DB
	}

	CheckInRequest struct {
		// BestBy 1/2025: double check all datatypes in this struct
		UserID    int    `json:"user_id"`
		Timestamp string `json:"timestamp"`
		Date      string `json:"date"`
		Payload   string `json:"payload"`
	}

	GetCheckInRequest struct {
		UserID string
	}
)

func (s *Service) PostCheckIn(req CheckInRequest) error {
	_, err := s.db.Exec("INSERT INTO check_in (user_id, timestamp, date, payload) VALUES (?, ?, ?, ?)", req.UserID, req.Timestamp, req.Date, req.Payload)
	return err
}

func (s *Service) GetCheckIns() (GetCheckInRequest, error) {
	return nil, nil
}

func main() {
	db, err := sql.Open("sqlite3", "./maain.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	CheckInService := &Service{db: db}
	fmt.Printf("+%v\n", CheckInService)
	createTableStatement := `CREATE TABLE IF NOT EXISTS check_in (id INTEGER PRIMARY KEY, user_id INTEGER, timestamp INTEGER, date TEXT, payload TEXT)`
	_, err = db.Exec(createTableStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, createTableStatement)
		return
	}

	e := echo.New()

	e.GET("/", MainPage)

	e.Logger.Fatal(e.Start(":8001"))
}

func Inc(c echo.Context) error {
	// retrieve the value of variable "metric" from request body
	metric := c.FormValue("metric")

	return c.String(http.StatusOK, metric)
}

func MainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
