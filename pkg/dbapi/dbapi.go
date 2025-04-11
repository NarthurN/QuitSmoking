package dbapi

import (
	"database/sql"
	"fmt"

	"github.com/NarthurN/QuitSmoking/internal/models"
)

func IsUserExists(db *sql.DB, username string) (bool, error) {
	var exists bool

	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM smokers WHERE username = $1)",
		username,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("error checking user existence: %v", err)
	}

	return exists, nil
}

func Add(db *sql.DB, smoker models.Smoker) (int, error) {
	result, err := db.Exec(
		`INSERT INTO smokers (name, username, password, stopped_smoking)
		VALUES ($1, $2, $3, $4)`,
		smoker.Name,
		smoker.Username,
		smoker.Password,
		smoker.StoppedSmoking,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func Get(db *sql.DB, username string) (models.Smoker, error) {
	row := db.QueryRow(`SELECT id, name, username, password stopped_smoking FROM smokers`)
	smoker := models.Smoker{}

	err := row.Scan(&smoker.ID, &smoker.Name, &smoker.Username, &smoker.Password, &smoker.StoppedSmoking)
	if err != nil {
		return smoker, err
	}

	return smoker, nil
}
