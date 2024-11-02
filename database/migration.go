package database

import (
    "database/sql"
    "log"
)

func InitMigration(db *sql.DB) {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

    _, err := db.Exec(query)
    if err != nil {
        log.Fatalf("Could not initialize database: %v", err)
    }
    log.Println("Database migrated successfully!")
}
