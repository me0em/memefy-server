// Package api represents REST API and implements
// database abstraction, url routing and some jwt layer functions
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// CreateUser represents user registration and update user data
//	POST   /api/user
//	PATCH  /api/user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}

	// get POST data
	decoder := json.NewDecoder(r.Body)
	userData := &User{}
	err := decoder.Decode(&userData)
	if err != nil || !userData.isValid() {
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		// register user in database
		err = InsertUser(db, *userData)
		if err != nil {
			http.Error(w, "database error", http.StatusBadRequest)
			return
		}

	case http.MethodPatch:
		// update field :column: in database
		var column, value string
		err = PatchUserData(db, userData.UserID, column, value)
		// TODO: дописать
	}

	// generate access token and response it
	respPayload := make(map[string]string)
	respPayload["access-token"] = GenerateToken(userData.UserID)
	response, err := json.Marshal(respPayload)
	if err != nil {
		http.Error(w, "something went wrong on the our side", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if wtf, err := w.Write(response); err != nil {
		// TODO: logging
		fmt.Printf("%v", wtf)
		panic(err)
	}
}

// SendMemes requests a few memes from ML-model and response it to client
// GET /api/meme
func SendMemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}

	// TODO: check authorization function, may be as decorator (everywhere)
	// TODO: authorization vs authentification
	userID, _, err := Authorization(w, r)
	// TODO: usually http.response with custom error text if it possible,
	// if it not - just send http.response (everywhere)
	if err != nil {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	// parse number of memes which will be returned
	limit := r.URL.Query().Get("limit")
	amount, err := strconv.Atoi(limit)
	if err != nil || amount < 1 || amount > 10 {
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	// TODO: send error to Telegram Bot
	// TODO: model response with error msg, need to handle it:
	// {"error_msg": "Username not converted to user_id", "user_id": "3"}
	transportStorage := &MemesTransport{UserID: userID, Amount: amount}
	jsonForModel, err := json.Marshal(transportStorage)
	if err != nil {
		fmt.Println(err)
	}
	jsonForModelBit := bytes.NewReader(jsonForModel)
	// TODO: config package with urls etc
	resp, err := http.Post("http://localhost:8228/model", "application/json", jsonForModelBit) //отправляю в модель
	if err != nil {
		fmt.Println(err)
	}

	memes, err := ioutil.ReadAll(resp.Body) //обрабатываю полученные даныне, получаю что нужно
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(memes); err != nil {
		panic(err)
	}

	// respPayload := make(map[string]string)
	// respPayload["limit"] = strconv.Itoa(limit)
	// response, err := json.Marshal(respPayload)
	// if err != nil {
	// 	http.Error(w, "something went wrong on the our side", http.StatusInternalServerError)
	// }

	// w.WriteHeader(http.StatusOK)
	// if wtf, err := w.Write(response); err != nil {
	// 	// TODO: logging
	// 	fmt.Printf("%v", wtf)
	// 	panic(err)
	// }
}

// TODO: sql-query to save reaction in databasedatabase

// GetReaction take meme_id and boolean feedbeak on it by user_id
// and make them love each other in one database table
// POST /api/reaction
func GetReaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}

	// get POST data
	decoder := json.NewDecoder(r.Body)
	var content ReactionContext
	err := decoder.Decode(&content)
	if err != nil {
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	respPayload := make(map[string]string)
	respPayload["reaction"] = content.Reaction
	response, err := json.Marshal(respPayload)
	if err != nil {
		http.Error(w, "something went wrong on the our side", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if wtf, err := w.Write(response); err != nil {
		// TODO: logging
		fmt.Printf("%v", wtf)
		panic(err)
	}

}

// TestThings represents useless shit which is will be
// removed as soon as possible
func TestThings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}

	userID, expiredTime, err := Authorization(w, r)
	if err != nil {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	respPayload := make(map[string]string)
	respPayload["userID"] = userID
	respPayload["Expired_time"] = fmt.Sprint(expiredTime)

	response, _ := json.Marshal(respPayload)
	w.WriteHeader(http.StatusOK)
	if wtf, err := w.Write(response); err != nil {
		// TODO: logging
		fmt.Printf("%v", wtf)
		panic(err)
	}
}
