package checkin

import (
	"context"
	"database/sql"
	"log"
)

type (
	Service struct {
		DB *sql.DB
	}

	AddCheckInRequest struct {
		// TODO: BestBy 1/2025: double check all datatypes in this struct
		UserID    int    `json:"user_id"`
		Timestamp string `json:"timestamp"`
		Date      string `json:"date"`
		Payload   string `json:"payload"`
	}

	GetCheckInRequest struct {
		UserID    string
		DateFrom  string
		DateUntil string
	}

	GetCheckInCountResponse struct {
		Count int64
	}

	GetCheckInResponse struct {
		ID        int    `json:"id"`
		UserID    int    `json:"user_id"`
		Timestamp string `json:"timestamp"`
		Date      string `json:"date"`
		Payload   string `json:"payload"`
	}
)

func (s *Service) AddCheckIn(req AddCheckInRequest) error {
	_, err := s.DB.Exec("INSERT INTO check_in (user_id, timestamp, date, payload) VALUES (?, ?, ?, ?)", req.UserID, req.Timestamp, req.Date, req.Payload)
	return err
}

func (s *Service) GetCheckInCount(userID string) (GetCheckInCountResponse, error) {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) from check_in WHERE user_id = ? GROUP BY user_id", userID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return GetCheckInCountResponse{Count: int64(count)}, err
}

func (s *Service) GetCheckIns(req GetCheckInRequest) ([]GetCheckInResponse, error) {
	var checkIn []GetCheckInResponse
	// TODO: support date range
	rows, err := s.DB.QueryContext(context.Background(), "SELECT * FROM check_in WHERE user_id = ?", req.UserID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r GetCheckInResponse
		if err := rows.Scan(&r.ID, &r.UserID, &r.Timestamp, &r.Date, &r.Payload); err != nil {
			log.Fatal(err)
		}
		checkIn = append(checkIn, r)
	}

	return checkIn, err
}

func (s *Service) DeleteCheckIn(id string) error {
	_, err := s.DB.Exec("DELETE FROM check_in WHERE id = ?", id)
	return err
}
