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

////////////////////
//	POST /api/user
//	PUT  /api/user
////////////////////
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}

	// get POST data
	decoder := json.NewDecoder(r.Body)
	userInfo := &User{}
	err := decoder.Decode(&userInfo)
	if err != nil || !userInfo.isValid() {
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		return
	}

	// TODO: think about PUT method and rewrite it
	// register or rewrite usera in our database
	switch r.Method {
	case http.MethodPost:
		// TODO: подрубить DB
	case http.MethodPut:
		// TODO: подрубить DB
	}

	// generate access token and response it
	respPayload := make(map[string]string)
	respPayload["access-token"] = GenerateToken(userInfo.UserID)
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


////////////////////
//	GET api/memes
////////////////////
func ThrowMemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}
	// parse number of memes which will be returned
	parsedLimit := r.URL.Query().Get("limit")
	userID, _, err := Authorization(w, r) //получаю айди(логин) юзера
	if err != nil{
		fmt.Println(err)
	}
	CountMeme, err := strconv.Atoi(parsedLimit)//количество мемов которое необходимо
	if err != nil || CountMeme <= 0 || CountMeme > 100 {
		http.Error(w, "another payload was expected", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	formodel := &GiveMeme{User: userID, Count: CountMeme}
	formodeljson, err := json.Marshal(formodel)//формирую джесон для модели
	if err != nil {
		fmt.Println(err)
	}

	formodeljsonbit := bytes.NewReader(formodeljson)
	resp, err := http.Post("http://localhost:8228/model", "application/json", formodeljsonbit) //отправляю в модель
	//TODO отправлять в бота ошибки
	//TODO модель присылает в ответ джесон такой с таким джесоном { "error_msg": "Username not converted to user_id", "user_id": "3" }, надо обрабатывать

	if err != nil {
		//fmt.Println(err)
		//var error = Error{Where:"ThrowMemes", What:"convert to json"}
		//b1, _ := json.Marshal(error)
		//error1 := bytes.NewReader(b1)
		//_, err := http.Get(" http://127.0.0.1:5000/error")
		//defer respe.Body.Close()
		fmt.Println(err)
	}
	meme, err := ioutil.ReadAll(resp.Body) //обрабатываю полученные даныне, получаю что нужно
	if err != nil {
		fmt.Println(err)

	}

	w.WriteHeader(http.StatusOK)

	if wtf, err := w.Write(meme); err != nil {
		// TODO: logging
		fmt.Printf("%v", wtf)
		panic(err)
	}
}

////////////////////
//	POST api/reaction
////////////////////
func HandleReaction(w http.ResponseWriter, r *http.Request) {
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

////////////////////
//	GET api/test
////////////////////
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