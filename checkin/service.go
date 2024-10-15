package checkin

import (
	"database/sql"
	"log"
)

type (
	Service struct {
		DB *sql.DB
	}

	CheckInRequest struct {
		// TODO: BestBy 1/2025: double check all datatypes in this struct
		UserID    int    `json:"user_id"`
		Timestamp string `json:"timestamp"`
		Date      string `json:"date"`
		Payload   string `json:"payload"`
	}

	GetCheckInRequest struct {
		UserID string
	}

	GetCheckInResponse struct {
		Metric int64
	}

	GetCheckinsRow struct {
	}
)

func (s *Service) AddCheckIn(req CheckInRequest) error {
	_, err := s.DB.Exec("INSERT INTO check_in (user_id, timestamp, date, payload) VALUES (?, ?, ?, ?)", req.UserID, req.Timestamp, req.Date, req.Payload)
	return err
}

func (s *Service) GetCheckInCount(userID string) (GetCheckInResponse, error) {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) from check_in WHERE user_id = ? GROUP BY user_id", userID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return GetCheckInResponse{Metric: int64(count)}, err
}
