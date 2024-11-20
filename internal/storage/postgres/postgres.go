package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(username string, password string, host string, port int, dbName string) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"
	resp, err := s.db.Exec(`INSERT INTO users (email, pass_hash) VALUES ($1, $2)`, email, passHash)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := resp.RowsAffected()
	if rowsAffected == 0 {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) UserByEmail(email string) (models.User, error) {
	const op = "storage.postgres.UserByEmail"

	var user models.User

	err := s.db.Get(&user, `SELECT id, email, pass_hash FROM users WHERE email = $1`, email)

	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserByID(id int64) (models.User, error) {
	const op = "storage.postgres.UserByID"

	var user models.User
	err := s.db.Get(&user, `SELECT id, email FROM users WHERE id = $1`, id)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) DeleteUser(id int64) error {
	const op = "storage.postgres.DeleteUser"

	res, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) IsAdmin(userID int64) (bool, error) {
	const op = "storage.postgres.isAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRow(userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(appID int32) (models.App, error) {
	const op = "storage.postgres.App"
	var app models.App
	err := s.db.Get(&app, `SELECT * FROM apps WHERE id = $1`, appID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, storage.ErrAppNotFound
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
