package config

//import "os"

// ServerPort is a port using for listening
var ServerPort = ":8085"
// DBInitReq is init database request
var DBInitReq = "http://default@memefy.fun:8123/memefy"
//SecretKey is a JWT secret key for generating tokens
var SecretKey = "wer6YTIFpojneEfe34fr4go8ukcyyjr45y8867"

//MLModelHost is a host of machine learning server
var MLModelHost = "http://127.0.0.1:8228/"

var ErrorHost = "http://127.0.0.1:8228/error"

//var MLModelHost = os.Getenv("MLMODELHOST")
//var SecretKey = os.Getenv("SECRETKAY")
//var DBInitReq = os.Getenv("DB")
//var ServerPort = os.Getenv("SERVERPORT")
//var ErrorHost = os.Getenv("ERRORHOST")