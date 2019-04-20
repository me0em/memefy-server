// Package api represents REST API and implements
// database abstraction, url routing and some jwt layer functions
package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mailru/go-clickhouse"
)

var db *sql.DB

// InitDB initialize the database
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("clickhouse", dataSourceName)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		log.Panic(err)
	}
	fmt.Println("Success!")
}

// GetAllUsers returns array with all users from database
func GetAllUsers(db *sql.DB) ([]User, error) {
	user := User{}
	users := []User{}

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		err = errors.New("failed to make a DB query")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// TODO: db binding: if err = rows.Scan(&user); err != nil {...}
		err = rows.Scan(&user.UserID, &user.IDType, &user.UserMetadata.Timestamp, &user.UserMetadata.Device, &user.UserMetadata.Model, &user.UserMetadata.DeviceLang, &user.UserMetadata.IPv4)
		if err != nil {
			err = errors.New("failed to load DB data to struct")
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		err = errors.New("something going wrong while try to read data from DB")
		return nil, err
	}

	return users, err
}

// GetUserByID returns user if he exists
func GetUserByID(db *sql.DB, userID string) (User, error) {
	user := User{}
	err := db.QueryRow("SELECT * FROM user WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.IDType, &user.UserMetadata.Timestamp, &user.UserMetadata.Device, &user.UserMetadata.Model, &user.UserMetadata.DeviceLang, &user.UserMetadata.IPv4)
	if err != nil {
		//err = errors.New("failed to make a DB query")
		return user, err
	}

	return user, err
}

// InsertUser register a user in the database
func InsertUser(db *sql.DB, user User) error {
	stmt, err := db.Prepare("INSERT INTO memefy.user " +
		"VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(user.UserID, user.IDType, user.UserMetadata.Timestamp, user.UserMetadata.Device, user.UserMetadata.Model, user.UserMetadata.DeviceLang, user.UserMetadata.IPv4)
	if err != nil {
		log.Fatal(err)
	}
	_ = stmt.Close()
	return err
}

// PatchUserData updates user data in the database
func PatchUserData(db *sql.DB, userID, column, value string) error {
	var err error
	var ex *sql.Stmt

	switch column {
	case "model":
		ex, err = db.Prepare("UPDATE memefy.user  SET `model`=? WHERE user_id=?")
	case "device":
		ex, err = db.Prepare("UPDATE memefy.user  SET `device`=? WHERE user_id=?")
	case "device_language":
		ex, err = db.Prepare("UPDATE memefy.user  SET `device_language`=? WHERE user_id=?")
	case "IPv4":
		ex, err = db.Prepare("UPDATE memefy.user  SET `IPv4`=? WHERE user_id=?")
	default:
		err = errors.New("there is not such column")
	}

	if err != nil {
		return err
	}

	if _, err = ex.Exec(value, userID); err != nil {
		err := errors.New("failed to make a DB query")
		return err
	}
	return nil
}

// DeleteUser delete a user from the database
func DeleteUser(db *sql.DB, userID string) error {
	var err error
	var ex *sql.Stmt

	ex, err = db.Prepare("DELETE FROM memefy.user WHERE user_id = ?")
	_, err = ex.Exec(userID)
	if err != nil {
		err := errors.New("failed to make a DB query")
		return err
	}
	return nil
}

func SaveReaction(db *sql.DB, reactions ReactionContext) error  {
	stmt, err := db.Prepare("INSERT INTO memefy.reactions " +
		"VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(reactions)
	if err != nil {
		log.Fatal(err)
	}
	_ = stmt.Close()
	return err
}