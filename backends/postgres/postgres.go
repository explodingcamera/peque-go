package postgres

import (
	"fmt"
	"log"
	"strconv"

	// load postgrs driver
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

// Service ...
type Service struct {
	db      *sqlx.DB
	Options Options
}

// Options ...
type Options struct {
	DatabaseName           string // default peque
	DatabaseUser           string // default peque
	DatabasePassword       string // default peque
	DatabaseHost           string // default localhost
	DatabasePort           int64  // default 5432
	DatabaseConnectTimeout int64  // default 0

	DatabaseSSLMode     string
	DatabaseSSLCert     string
	DatabaseSSLKey      string
	DatabaseSSLRootCert string
}

// Connect initializes a new Service
func Connect(options Options) (*Service, error) {

	var connectString string

	if options.DatabaseName == "" {
		options.DatabaseName = "peque"
	}
	connectString += "dbname=" + options.DatabaseName

	if options.DatabaseUser == "" {
		options.DatabaseUser = "peque"
	}
	connectString += "user=" + options.DatabaseUser

	if options.DatabaseHost != "" {
		connectString += "host=" + options.DatabaseHost
	}

	if options.DatabasePort != 0 {
		connectString += "port=" + strconv.FormatInt(options.DatabasePort, 10)
	}

	if options.DatabaseSSLMode != "" {
		connectString += "sslmode=" + options.DatabaseSSLMode
	}

	if options.DatabaseConnectTimeout != 0 {
		connectString += "connect_timeout=" + strconv.FormatInt(options.DatabaseConnectTimeout, 10)
	}

	if options.DatabaseSSLCert != "" {
		connectString += "sslcert=" + options.DatabaseSSLCert
	}

	if options.DatabaseSSLKey != "" {
		connectString += "sslkey=" + options.DatabaseSSLKey
	}

	if options.DatabaseSSLRootCert != "" {
		connectString += "sslrootcert=" + options.DatabaseSSLRootCert
	}

	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		log.Fatalln(err)
	}

	service := Service{
		db: db,
	}

	roleName := "test"
	role, err := service.db.Exec("SELECT 1 FROM pg_roles WHERE rolname=?", roleName)
	fmt.Println(role, err)

	return &service, nil
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
