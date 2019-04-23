// Package api represents REST API and implements
// database abstraction, url routing and some jwt layer functions
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"memefy-server/config"
	"net/http"
	"strconv"
	"time"
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
	userData.UserMetadata.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	err := decoder.Decode(&userData)

	if err != nil || !userData.isValid() {
		ErrorsForTelegramBot(err, "<CreateUser>: Invalid user data")
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		// register user in database
		_, err := GetUserByID(db, userData.UserID)
		if err == nil {
			http.Error(w, "another payload was expected", http.StatusBadRequest)
			return
		}

		err = InsertUser(db, *userData)
		if err != nil {
			ErrorsForTelegramBot(err, "<CreateUser>: Can not register user in database")
			http.Error(w, "database error", http.StatusBadRequest)
			return
		}

	case http.MethodPatch:
		// TODO: дописать
		// update field :column: in database
		var column, value string
		err = PatchUserData(db, userData.UserID, column, value)
		if err != nil {
			ErrorsForTelegramBot(err, "<CreateUser>: Can not update user data")
			return
		}

	}

	// generate access token and response it
	respPayload := make(map[string]string)
	respPayload["access-token"] = GenerateToken(userData.UserID)
	respPayload["refresh-token"] = GenerateToken(userData.UserID + "-refresh")
	response, err := json.Marshal(respPayload)
	if err != nil {
		ErrorsForTelegramBot(err, "<CreateUser>: Critical error while generate JWT")
		http.Error(w, "something went wrong on the our side", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		ErrorsForTelegramBot(err, "[critical!] <CreateUser>: Can not sends response with token")
	}
}

// SendMemes requests a few memes from ML-model and response it to client
func SendMemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}

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
	if err != nil || amount < 1 || amount > 30 {
		ErrorsForTelegramBot(err, "<SendMemes>: Invalid limit number")
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	// TODO: model response with error msg, need to handle it:
	transportStorage := &MemesTransport{UserID: userID, Amount: amount}
	jsonForModel, err := json.Marshal(transportStorage)
	if err != nil {
		fmt.Println(err)
	}
	jsonForModelBit := bytes.NewReader(jsonForModel)
	resp, err := http.Post(config.MLModelHost, "application/json", jsonForModelBit) //отправляю в модель
	if err != nil {
		ErrorsForTelegramBot(err, "<SendMemes>: Error while sending data to ML model")
		fmt.Println(err)
	}

	decoder := json.NewDecoder(resp.Body)
	memes := &MemesFromModel{}
	err = decoder.Decode(&memes)
	if err != nil || len(memes.Memes) < amount {
		//TODO: что делать если модель прислала меньше мемов чем запрошено?
		ErrorsForTelegramBot(err, "SendMemes")
		fmt.Println(err)
	}
	var arrtext = make([]string, 0, amount)
	var arrhash = make([]string, 0, amount)
	for i := 0; i < len(memes.Memes); i++ {
		var text, hash string = GetMemeText(db, memes.Memes[i])
		arrtext = append(arrtext, text)
		arrhash = append(arrhash, hash)
	}
	memesForUser := &ResponseMemes{}
	memesForUser.Memes = arrhash
	memesForUser.Text = arrtext
	jsonForUser, err := json.Marshal(memesForUser)
	if err != nil {
		ErrorsForTelegramBot(err, "<SendMemes>: Error while reciving data from ML model")
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(jsonForUser); err != nil {
		ErrorsForTelegramBot(err, "[critical!] <SendMemes>: Error while response with memes")
	}
}

// GetReaction take meme_id and boolean feedbeak on it by user_id
// and make them love each other in one database table
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
		ErrorsForTelegramBot(err, "<GetReaction>: Invalid reaction dat	a")
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	userID, _, err := Authorization(w, r)
	content.UserID = userID
	content.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	err = SaveReaction(db, content)
	if err != nil {
		ErrorsForTelegramBot(err, "<GetReaction>: Can not save reaction into database")
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RefreshToken allows to refresh token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}
	userID, err := CheckRefreshToken(w, r)
	if err != nil {
		ErrorsForTelegramBot(err, "<RefreshToken>: Looks like we are under hackers attack")
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	respPayload := make(map[string]string)
	respPayload["access-token"] = GenerateToken(userID)
	response, err := json.Marshal(respPayload)
	if err != nil {
		ErrorsForTelegramBot(err, "<RefreshToken>: Refresh token is ok, but error while try to generate new access-token")
		http.Error(w, "something went wrong on the our side", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		ErrorsForTelegramBot(err, "<RefreshToken>: Can not send response")
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

// ErrorsForTelegramBot sends error msg to Error Telegram Bot
func ErrorsForTelegramBot(error error, where string) {

	errorr := &ErrorForTelegram{Error: error, Where: where}
	jsonForModel, err := json.Marshal(errorr)
	if err != nil {
		fmt.Println(err)
	}
	jsonForModelBit := bytes.NewReader(jsonForModel)
	re, err := http.Post(config.ErrorHost, "application/json", jsonForModelBit) //отправляю в модель
	fmt.Println(re)
	if err != nil {
		fmt.Println(err)
	}
}
