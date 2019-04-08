// Package api represents REST API and implements
// database abstraction, url routing and some jwt layer functions
package api

import (
	"time"
	"unicode"
)

type User struct {
	UserID string 			 `db:"user_id" json:"user_id"`
	IDType string 			 `db:"id_type" json:"id_type"`
	UserMetadata struct {
		Timestamp  time.Time `db:"timestamp" json:"timestamp"`
		Device     string    `db:"device" json:"device"`
		Model      string    `db:"model" json:"model"`
		DeviceLang string    `db:"device_language" json:"device_language"`
		IPv4       string    `db:"IPv6" json:"IPv6"`
	}
}


// TODO: UserID that does not exist yet
// This method allow you to check the fields of
// User structure for correctness.
func (user User) isValid() bool {
	for _, char := range user.UserID {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < 0 || char > 9) && char != '_' {
			return false
		}
	}
	switch {
	case len(user.UserID) > 15 || len(user.UserID) < 5:
		return false
	case unicode.IsDigit( []rune(user.UserID)[0] ):
		return false
	case user.UserID[0] == '_':
		return false
	default:
		return true
	}
}

type ReactionContext struct {
	UserID string
	MemeID string
	Reaction string
}