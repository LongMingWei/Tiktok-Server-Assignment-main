package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type SQLClient struct {
	db *sql.DB
}

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (c *SQLClient) InitClient(ctx context.Context, username, password, address, database string) error {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, address, database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	c.db = db

	_, err = c.db.ExecContext(ctx, `
		CREATE TABLE messages (
			id INT AUTO_INCREMENT PRIMARY KEY,
			room_id VARCHAR(255) NOT NULL,
			sender VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			timestamp INT NOT NULL
		)`)
	if err != nil {
		return err
	}

	return nil
}

func (c *SQLClient) SaveMessage(ctx context.Context, roomID string, message *Message) error {

	stmt, err := c.db.PrepareContext(ctx, "INSERT INTO messages (room_id, sender, message, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, roomID, message.Sender, message.Message, message.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (c *SQLClient) GetMessagesByRoomID(ctx context.Context, roomID string, start, end int64, reverse bool) ([]*Message, error) {
	var (
		getMessages []*Message
		err         error
	)

	var order string
	if reverse {
		order = "DESC"
	} else {
		order = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT sender, message, timestamp
		FROM messages
		WHERE room_id = ?
		ORDER BY timestamp %s
		LIMIT ?, ?
	`, order)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, roomID, start, end-start+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		temp := &Message{}
		err = rows.Scan(&temp.Sender, &temp.Message, &temp.Timestamp)
		if err != nil {
			return nil, err
		}
		getMessages = append(getMessages, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return getMessages, nil
}
