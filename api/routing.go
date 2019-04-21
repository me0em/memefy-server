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
	//fmt.Println(userData)
	err := decoder.Decode(&userData)

	if err != nil || !userData.isValid() {
		ErrorsForTelegramBot(err, "CreateUser1")
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		// register user in database

		err = InsertUser(db, *userData)

		if err != nil {
			ErrorsForTelegramBot(err, "CreateUser2")
			http.Error(w, "database error", http.StatusBadRequest)
			return
		}

	case http.MethodPatch:
		// update field :column: in database
		var column, value string
		err = PatchUserData(db, userData.UserID, column, value)
		// TODO: дописать
		if err != nil {
			ErrorsForTelegramBot(err, "CreateUser3")
			return
		}

	}

	// generate access token and response it
	respPayload := make(map[string]string)
	respPayload["access-token"] = GenerateToken(userData.UserID)
	respPayload["refresh-token"] = GenerateToken(userData.UserID + "-refresh")
	response, err := json.Marshal(respPayload)
	if err != nil {
		ErrorsForTelegramBot(err, "CreateUser4")
		http.Error(w, "something went wrong on the our side", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		ErrorsForTelegramBot(err, "CreateUser5")
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
		ErrorsForTelegramBot(err, "SendMemes")
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	// parse number of memes which will be returned
	limit := r.URL.Query().Get("limit")
	amount, err := strconv.Atoi(limit)
	if err != nil || amount < 1 || amount > 10 {
		ErrorsForTelegramBot(err, "SendMemes")
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	// TODO: send error to Telegram Bot
	// TODO: model response with error msg, need to handle it:
	// {"error_msg": "Username not converted to user_id", "user_id": "3"}
	transportStorage := &MemesTransport{UserID: userID, Amount: amount}
	jsonForModel, err := json.Marshal(transportStorage)
	if err != nil {
		ErrorsForTelegramBot(err, "SendMemes")
		fmt.Println(err)
	}
	jsonForModelBit := bytes.NewReader(jsonForModel)
	// TODO: config package with urls etc
	resp, err := http.Post("http://localhost:8228/model", "application/json", jsonForModelBit) //отправляю в модель
	if err != nil {
		ErrorsForTelegramBot(err, "SendMemes")
		fmt.Println(err)
	}

	memes, err := ioutil.ReadAll(resp.Body) //обрабатываю полученные даныне, получаю что нужно
	if err != nil {
		ErrorsForTelegramBot(err, "SendMemes")
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(memes); err != nil {
		ErrorsForTelegramBot(err, "SendMemes")
		panic(err)
	}

	// respPayload := make(map[string]string)
	// respPayload["limit"] = strconv.Itoa(limit)
	// response, err := json.Marshal(respPayload)
	// if err != nil {
	// 	http.Error(w, "something went wrong on the our side", http.StatusInternalServerError)
	// }

	//w.WriteHeader(http.StatusOK)
	//if wtf, err := w.Write(response); err != nil {
	//	// TODO: logging
	//	fmt.Printf("%v", wtf)
	//	panic(err)
	//}
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
		ErrorsForTelegramBot(err, "GetReaction")
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	userID, _, err := Authorization(w, r)
	content.UserID = userID
	//TODO закидывать время
	//content.Timestamp = time.Now()
	err = SaveReaction(db, content)
	if err != nil {
		ErrorsForTelegramBot(err, "GetReaction")
		//http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}
	//w.WriteHeader(http.StatusOK)
	//if wtf, err := w.Write(response); err != nil {
	//	// TODO: logging
	//	ErrorsForTelegramBot(err, "GetReaction")
	//	fmt.Printf("%v", wtf)
	//	panic(err)
	//}

}

// RefreshToken allows to refresh token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}
	userID, err := CheckRefreshToken(w, r)
	if err != nil {
		ErrorsForTelegramBot(err, "TestThings")
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	respPayload := make(map[string]string)
	respPayload["access-token"] = GenerateToken(userID)
	response, err := json.Marshal(respPayload)
	if err != nil {
		ErrorsForTelegramBot(err, "CreateUser4")
		http.Error(w, "something went wrong on the our side", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		// TODO: logging
		ErrorsForTelegramBot(err, "CreateUser5")
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
		ErrorsForTelegramBot(err, "TestThings")
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
		ErrorsForTelegramBot(err, "TestThings")
		panic(err)
	}
}

func ErrorsForTelegramBot(error error, where string) {
	//errorr := &ErrorForTelegram{Error:error, Where:where}
	//jsonForModel, err := json.Marshal(errorr)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//jsonForModelBit := bytes.NewReader(jsonForModel)
	//re, err := http.Post("http://127.0.0.1:5000/", "application/json", jsonForModelBit) //отправляю в модель
	//fmt.Println(re)
	//if err != nil {
	//	fmt.Println(err)
	//}
	fmt.Println(error, where)

}
