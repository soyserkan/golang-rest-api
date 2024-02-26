package models

import (
	"time"

	"example.com/rest-api/db"
)

type Event struct {
	ID          int64
	Name        string `binding:"required"`
	Description string
	Location    string
	DateTime    time.Time
	UserId      int64
}

func (e *Event) Save() error {
	query := "INSERT INTO events (name, description, location, dateTime, user_id) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserId)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	e.ID = id

	return err
}

func GetEvent(eventId int64) (*Event, error) {
	sql := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(sql, eventId)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func GetAllEvents() ([]Event, error) {

	var events = []Event{}

	sql := "SELECT * FROM events"
	rows, err := db.DB.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)
		if err != nil {
			return nil, err
		}
		events = append(events, event)

	}

	return events, nil
}

func (event Event) Update() error {
	sql := "UPDATE events SET name = ?, description = ?, location = ?, dateTime = ?, user_id = ?"
	stmt, err := db.DB.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.UserId)
	return err
}

func (event Event) Delete() error {
	sql := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

func (event Event) Register(userId int64) error {
	sql := "INSERT INTO registrations(event_id, user_id) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID, userId)
	return err
}

func (event Event) CancelRegistration(userId int64) error {
	sql := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID, userId)
	return err
}
