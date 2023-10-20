package main

import "os"

var BROKER_URL = os.Getenv("BROKER_URL")
var DB_HOST = os.Getenv("DB_HOST")
var DB_PORT = os.Getenv("DB_PORT")
var APP_PORT = os.Getenv("APP_PORT")
var TOPIC_NAME = os.Getenv("TOPIC_NAME")
