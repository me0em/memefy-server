// Package api represents REST API and implements
// database abstraction, url routing and some jwt layer functions
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)



// URL schema
//
//	POST /user
//	PUT	 /user
//	GET	 /memes
//	POST /reaction
//	GET  /test



//	POST /api/user
//	PUT  /api/user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}

	// get POST data
	decoder := json.NewDecoder(r.Body)
	content := &User{}
	err := decoder.Decode(&content)
	if err != nil || !content.isValid() {
		http.Error(w, "expected another payload", http.StatusBadRequest)
		return
	}

	// TODO: think about PUT method and rewrite it
	// register or rewrite user in our database
	switch r.Method {
	case http.MethodPost:
		//content.save()
		TestDB()
	case http.MethodPut:
		//content.save()
		TestDB()
	}

	// generate access token and response it
	respBody := make(map[string]string)
	respBody["access-token"] = GenerateToken(content.UserID)
	response, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, "sorry, something went wrong", http.StatusMethodNotAllowed)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}


//	GET /memes
func ThrowMemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid API method", http.StatusMethodNotAllowed)
		return
	}

	// get number of memes which we will return
	limitstr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitstr)
	if err != nil || limit <= 0 || limit > 10 {
		http.Error(w, "invalid limit argument", http.StatusMethodNotAllowed)
		return
	}

	respBody := make(map[string]string)
	respBody["limit"] = strconv.Itoa(limit)

	response, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, "sorry, something went wrong", http.StatusMethodNotAllowed)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

//	POST /reaction
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
		http.Error(w, "expected another payload", http.StatusBadRequest)
		return
	}
	fmt.Println(content)


	respBody := make(map[string]string)
	respBody["reaction"] = content.Reaction

	response, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, "sorry, something went wrong", http.StatusMethodNotAllowed)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

//	GET /test
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

	respData := make(map[string]string)
	respData["userID"] = userID
	respData["Expired_time"] = fmt.Sprint(expiredTime)

	TestDB()

	respJson, _ := json.Marshal(respData)
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}