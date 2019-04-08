// Package db represents databases abstractions
package api

import (
	"database/sql"
	"errors"
	"log"
)


// Function returns
// array with all of users from database
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


// Function returns
// user if he exists
func GetUserByID(db *sql.DB, userID string) (User, error) {
	user := User{}
	err := db.QueryRow("SELECT * FROM user WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.IDType,  &user.UserMetadata.Timestamp, &user.UserMetadata.Device, &user.UserMetadata.Model, &user.UserMetadata.DeviceLang, &user.UserMetadata.IPv4)
	if err != nil {
		//err = errors.New("failed to make a DB query")
		return user, err
	}
	return user, err
}


// Function returns error
// Insert a user to the database
func InsertUser(db *sql.DB, user User) error {
	stmt, err := db.Prepare("INSERT INTO memefy.user " +
		"VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(user.UserID, user.IDType, user.UserMetadata.Timestamp, user.UserMetadata.Device, user.UserMetadata.Model, user.UserMetadata.DeviceLang, user.UserMetadata.IPv4)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	_ = stmt.Close()
	return err
}


// Function returns error
// Change some data in the user entity
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

	_, err = ex.Exec(value, userID)
	if err != nil {
		err := errors.New("failed to make a DB query")
		return err
	}
	return nil
}


// Function returns error
// Delete User from the database
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

