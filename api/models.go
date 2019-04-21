// Package api represents REST API and implements
// database abstraction, url routing and some jwt layer functions
package api

import (
	"unicode"
)

// User is an entity for json and database abstracts
// Represent user data
type User struct {
	UserID       string `db:"user_id" json:"user_id"`
	IDType       string `db:"id_type" json:"id_type"`
	UserMetadata struct {
		Timestamp  string `db:"timestamp" json:"timestamp"`
		Device     string `db:"device" json:"device"`
		Model      string `db:"model" json:"model"`
		DeviceLang string `db:"device_language" json:"device_language"`
		IPv4       string `db:"IPv4" json:"IPv4"`
	} `json:"user_metadata"`
}

// TODO: UserID that does not exist yet
// This method allow you to check the fields of
// User structure for correctness.
func (user User) isValid() bool {
	for _, char := range user.UserID {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < 0 || char > 9) && char != '_' && char != '@' && char != '.' && char != '-' && char != ':' && char != ' ' {
			return false
		}
	}
	switch {
	case len(user.UserID) > 30 || len(user.UserID) < 5:
		return false
	case unicode.IsDigit([]rune(user.UserID)[0]):
		return false
	case user.UserID[0] == '_':
		return false
	default:
		return true
	}
}

// ReactionContext and me doesn't give a fuck why it's necessary
type ReactionContext struct {
	UserID    string `json:"user_id"`
	MemeID    int    `json:"meme_id"`
	Reaction  int    `json:"reaction"`
	Timestamp string `json:"timestamp"`
}

// MemesTransport is a struct for request to ML model
type MemesTransport struct {
	UserID string `json:"user_id"`
	Amount int    `json:"count_meme"`
}

// ErrorForTelegram is a struct for storage errors sendings to error bot
type ErrorForTelegram struct {
	Error error  `json:"error"`
	Where string `json:"where"`
}
type MemesFromModel struct {
	UserID string `json:"user_id"`
	Memes  []int  `json:"meme_id"`
}
type ResponseMemes struct {
	Memes []string `json:"meme_id"`
	Text  []string `json:"meme_text"`
}
type MemeWithText struct {
	Hash string
	Text string
}
