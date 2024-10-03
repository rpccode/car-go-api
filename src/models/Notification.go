package models

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID      int       `json:"id"`
	UserID  int       `json:"user_id" binding:"required"`
	Message string    `json:"message" binding:"required"`
	SentAt  time.Time `json:"sent_at"`
}

// Enviar notificación
func (n *Notification) Send(db *sql.DB) error {
	query := `INSERT INTO notifications (user_id, message, sent_at) 
              VALUES ($1, $2, $3) RETURNING id`
	return db.QueryRow(query, n.UserID, n.Message, time.Now()).Scan(&n.ID)
}
func GetNotificationsByUserID(db *sql.DB, userID int) ([]Notification, error) {
	query := `SELECT id, user_id, message, created_at FROM notifications WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.SentAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
