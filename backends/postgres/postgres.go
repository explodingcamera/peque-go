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

func (s Service) Install() {
	sqlx.LoadFile(s.db, "")
}
