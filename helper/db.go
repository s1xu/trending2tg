package helper

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"trending2telbot/model"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase(filepath string) (*sql.DB, error) {
	db, err := initDB(filepath)
	if err != nil {
		return nil, err
	}

	err = createTable(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// initDB initializes the SQLite database connection.
func initDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("db nil")
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)
	return db, nil
}

// createTable creates a table in the database if it does not already exist.
func createTable(db *sql.DB) error {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS trends(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		language TEXT,
		url TEXT,
		create_at DATETIME DEFAULT (datetime('now', 'localtime')),
		UNIQUE(title,language)
	);`
	_, err := db.Exec(sqlTable)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	return nil
}

// createIndex creates an index on the title column of the trends table to improve query performance.
func createIndex(db *sql.DB) error {
	sqlIndex := `CREATE INDEX IF NOT EXISTS idx_language ON trends(title);`
	_, err := db.Exec(sqlIndex)
	if err != nil {
		log.Printf("Failed to create index: %v", err)
		return err
	}
	return nil
}

func InsertTrendIfNotExists(db *sql.DB, trend model.Trends) error {
	sqlAddItem := `INSERT OR IGNORE INTO trends(title, language, url, create_at) VALUES(?, ?, ?, ?)`
	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Convert the current time to Beijing time
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("Failed to load location: %v", err)
		return err
	}
	beijingTime := time.Now().In(location).Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(trend.Title, trend.Language, trend.Url, beijingTime)
	if err != nil {
		log.Printf("Failed to insert trend: %v", err)
	}
	return err
}

func InsertIfNotExists(db *sql.DB, trend model.Trends) (bool, error) {
	exists, err := CheckIfTitleExists(db, trend.Title)
	if err != nil {
		log.Printf("Error checking if title exists: %v", err)
		return false, err
	}
	if exists {
		log.Printf("Title already exists, skipping insert: %s", trend.Title)
		return false, nil
	}

	sqlAddItem := `INSERT INTO trends(title, language, url, create_at) VALUES(?, ?, ?, ?)`
	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return false, err
	}
	defer stmt.Close()

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("Failed to load location: %v", err)
		return false, err
	}
	beijingTime := time.Now().In(location).Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(trend.Title, trend.Language, trend.Url, beijingTime)
	if err != nil {
		log.Printf("Failed to insert trend: %v", err)
		return false, err
	}
	return true, nil
}

func CheckIfTitleExists(db *sql.DB, title string) (bool, error) {
	sqlRead := `SELECT title FROM trends WHERE title = ?`
	row := db.QueryRow(sqlRead, title)
	var titleRead string
	err := row.Scan(&titleRead)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("Failed to check if title exists: %v", err)
		return false, err
	}
	return true, nil
}
