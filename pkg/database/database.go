package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	Login(*goth.User) (int64, error)

	GetUser(string) (*User, error)

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type User struct {
	UserID  int    `db:"user_id" json:"user_id"`
	Email   string `db:"email" json:"email"`
	Name    string `db:"user_name" json:"user_name"`
	Picture string `db:"picture" json:"picture"`
}

type service struct {
	db *sql.DB
}

var (
	dbname     = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Opening a driver typically will not attempt to connect to the database.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Login(user *goth.User) (int64, error) {
	var userID int64
	err := s.db.QueryRow("select user_id from users where auth_id = ?", user.UserID).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error fetching user")
		return -1, err
	}

	if err == sql.ErrNoRows {
		res, err := s.db.Exec("insert into users(auth_id, name, email, picture) values(?, ?, ?, ?)", user.UserID, user.Name, user.Email, user.AvatarURL)
		if err != nil {
			return -1, fmt.Errorf("Could not save user. Try again later")
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return -1, fmt.Errorf("could not retreive new user id. Try again later")
		}

		userID = lastID
	}

	log.Printf("user id %d", userID)
	return userID, nil

}

func (s *service) GetUser(auth_id string) (*User, error) {
	var user_id int
	var name, email, picture string

	err := s.db.QueryRow("select user_id, name, email, picture from blueprint.users where auth_id = ?", auth_id).Scan(&user_id, &name, &email, &picture)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error fetching user: %s", err)
		return nil, err
	}

	user := &User{
		UserID:  user_id,
		Name:    name,
		Email:   email,
		Picture: picture,
	}

	log.Printf("User data: %v", user)
	return user, nil

}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	return s.db.Close()
}
