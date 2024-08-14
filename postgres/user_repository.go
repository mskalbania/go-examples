package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/*
CRUD example with postgres pgx driver, using two tables: users and user_data.
Using transactions to ensure data consistency.
*/

var (
	insertUserQuery           = "INSERT INTO users (id, user_name) VALUES ($1, $2)"
	insertUserDataQuery       = "INSERT INTO user_data (user_id, name, surname, role) VALUES ($1, $2, $3, $4)"
	selectUserByUsernameQuery = "SELECT id FROM users WHERE user_name=$1"
	selectUsersQuery          = "SELECT id, user_name, name, surname, role FROM users LEFT JOIN user_data ud ON users.id = ud.user_id"
	deleteUserDataQuery       = "DELETE FROM user_data WHERE user_id=$1"
	deleteUserQuery           = "DELETE FROM users WHERE id=$1"
	updateUserDataQuery       = "UPDATE user_data SET name=$1, surname=$2, role=$3 WHERE user_id=$4"
)

var UserAlreadyExistsError = fmt.Errorf("user with given username already exists")

type User struct {
	ID       string
	Username string
	Name     string
	Surname  string
	Role     string
}

func SaveUser(user User) (string, error) {
	if exist, err := userExists(user.Username); err != nil {
		return "", fmt.Errorf("error saving user: %w", err)
	} else if exist {
		return "", fmt.Errorf("error saving user: %w", UserAlreadyExistsError)
	}
	id := uuid.New().String()
	tx, err := connection.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return "", fmt.Errorf("error opening transaction while saving user: %w", err)
	}
	if _, err := connection.Exec(context.Background(), insertUserQuery, id, user.Username); err != nil {
		return "", rollback(err, "saving user record", tx)
	}
	if _, err := connection.Exec(context.Background(), insertUserDataQuery, id, user.Name, user.Surname, user.Role); err != nil {
		return "", rollback(err, "saving user data record", tx)
	}
	if err := tx.Commit(context.Background()); err != nil {
		return "", fmt.Errorf("error commiting transaction while saving user: %w", err)
	}
	return id, nil
}

func GetAllUsers() ([]*User, error) {
	var users []*User
	rows, err := connection.Query(context.Background(), selectUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Role); err != nil {
			return nil, fmt.Errorf("error reading users: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}

func DeleteUser(id string) error {
	tx, err := connection.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return fmt.Errorf("error opening transaction while deleting user: %w", err)
	}
	if _, err := connection.Exec(context.Background(), deleteUserDataQuery, id); err != nil {
		return rollback(err, "deleting user data record", tx)
	}
	if _, err := connection.Exec(context.Background(), deleteUserQuery, id); err != nil {
		return rollback(err, "deleting user record", tx)
	}
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("error commiting transaction while deleting user: %w", err)
	}
	return nil
}

func UpdateUser(user User) error {
	if _, err := connection.Exec(context.Background(), updateUserDataQuery, user.Name, user.Surname, user.Role, user.ID); err != nil {
		return fmt.Errorf("error updating user data record %w", err)
	}
	return nil
}

func rollback(err error, message string, tx pgx.Tx) error {
	if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
		return fmt.Errorf("error rolling back transaction while %s: %w", message, rollbackErr)
	}
	return fmt.Errorf("error %s: %w", message, err)
}

func userExists(username string) (bool, error) {
	var found string
	err := connection.QueryRow(context.Background(), selectUserByUsernameQuery, username).Scan(&found)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		err = fmt.Errorf("error checking if user exists: %w", err)
	}
	return err == nil, err
}
