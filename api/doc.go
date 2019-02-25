/*
Package api provides a abstracts for building a REST api
with JWT auth layer and Clickhouse database.



- - - routing.go

func TestThings(w http.ResponseWriter, r *http.Request) {
	Handled function for test some things (wow)


func CreateUser(w http.ResponseWriter, r *http.Request) {
	Handled function for register users


- - - jwt.go

func GenerateToken(username string) string {
	Function for generate token
	Return: token

func Authorization(w http.ResponseWriter, r *http.Request) (string, float64, error) {
	Function for parse the token and check it for correctness
	Return: userID (token param), expired time, error

*/


package api