package config

// ServerPort is a port using for listening
var ServerPort = ":8085"

// DBInitReq is init database request
var DBInitReq = "http://default::tyz1214@127.0.0.1:8123/memefy?parseTime=true"

// SecretKey is a JWT secret key for generating tokens
var SecretKey = "wer6YTIFpojneEfe34fr4go8ukcyyjr45y8867"

// MLModelHost is a host of machine learning server
var MLModelHost = "http://127.0.0.1:5000/"
