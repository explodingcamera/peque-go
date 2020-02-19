package postgres

import (
	"log"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

// Service ...
type Service struct {
	db *sqlx.DB
}

// Init initializes a new Service
func Init() *Service {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	service := Service{
		db: db,
	}

	return &service
}

// WriteMessage writes a new message to a stream
func (s Service) WriteMessage(data string, streamID string) error {
	_, err := s.db.NamedExec(
		`SELECT write_message(messageID, streamID, messageType, data) VALUES (:messageID, :streamID, :data)`,
		map[string]interface{}{
			"messageID": uuid.NewV4(),
			"streamID":  streamID,
			"data":      data,
		})

	return err
}

// Install messageDB on postgres server
func (s Service) Install() error {
	prefix := "messagedb/database/"
	loadFiles := []string{
		// Main install
		"roles/message-store.sql",
		"schema/message-store.sql",
		"extensions/pgcrypto.sql",
		"tables/messages.sql",

		// Functions
		"types/message.sql",
		"functions/message-store-version.sql",
		"functions/hash-64.sql",
		"functions/acquire-lock.sql",
		"functions/category.sql",
		"functions/is-category.sql",
		"functions/id.sql",
		"functions/cardinal-id.sql",
		"functions/stream-version.sql",
		"functions/write-message.sql",
		"functions/get-stream-messages.sql",
		"functions/get-category-messages.sql",
		"functions/get-last-stream-message.sql",

		// Indexes
		"indexes/messages-id.sql",
		"indexes/messages-stream.sql",
		"indexes/messages-category.sql",

		// Privileges
		"privileges/schema.sql",
		"privileges/table.sql",
		"privileges/sequence.sql",
		"privileges/functions.sql",
		"privileges/views.sql",

		// Views
		"views/stream-summary.sql",
		"views/type-summary.sql",
		"views/stream-type-summary.sql",
		"views/type-stream-summary.sql",
		"views/category-type-summary.sql",
		"views/type-category-summary.sql",
	}

	for _, file := range loadFiles {
		_, err := sqlx.LoadFile(s.db, prefix+file)
		if err != nil {
			return err
		}
	}

	return nil
}
