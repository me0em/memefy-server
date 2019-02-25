// Package api represents REST API and implements
// database abstraction, url routing and some jwt layer functions
package api

import (
	"unicode"
)

// TODO: что хранить на сервере?
// Чтобы максимально privacy, но максимально полезно
type User struct {
	UserID       string			// json: "UserID"
	IDType       string			// json: "IDType"
	UserMetadata struct {		// json: "UserMetadata"
		Device         string   // json: "Device"
		Model          string	// json: "Model"
		IPv4           string	// json: "IPv4"
		DeviceLanguage string	// json: "DeviceLanguage"
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