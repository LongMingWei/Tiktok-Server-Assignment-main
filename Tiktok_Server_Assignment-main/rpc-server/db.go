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

	return nil
}

func (c *SQLClient) SaveMessage(ctx context.Context, roomID string, message *Message) error {

	stmt, err := c.db.PrepareContext(ctx, "INSERT INTO messages (room_id, sender, message, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

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
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	rows, err := stmt.QueryContext(ctx, roomID, start, end-start+1)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

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
