package config

import "os"

// ServerPort is a port using for listening
var ServerPort = os.Getenv("SERVERPORT")

// DBInitReq is init database request
//var DBInitReq = "http://default@memefy.fun:8123/memefy"
//var DBInitReq = "http://default@127.0.0.1:8123/memefy"
var DBInitReq = os.Getenv("DBINITREQ")
// SecretKey is a JWT secret key for generating tokens
//var SecretKey = "wer6YTIFpojneEfe34fr4go8ukcyyjr45y8867"
var SecretKey = os.Getenv("SECRETKAY")
// MLModelHost is a host of machine learning server
//var MLModelHost = "http://127.0.0.1:8228/"
var MLModelHost = os.Getenv("MLMODELHOST")
//var Port = os.Getenv("SERVERPORT")