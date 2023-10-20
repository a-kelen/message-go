package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func createMessageHandler(w http.ResponseWriter, r *http.Request) {
	body := parseBody[Message](w, r)
	sendToQueue(producer, body.Text)
	fmt.Fprint(w, "Message added to queue: "+body.Text)
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	requestChannel := make(chan MessageItem)
	messageChannels[requestID] = requestChannel
	fmt.Printf("Client connected with ID = %s \n", requestID)
	var messages []MessageItem
	db.Find(&messages)

	w.Header().Set("Content-Type", "text/event-stream")
	messagesSuit := make([]string, len(messages))

	for i, message := range messages {
		messagesSuit[i] = fmt.Sprintf("[%d] - %s\n", message.ID, message.Text)
	}
	fmt.Fprint(w, strings.Join(messagesSuit, ""))

	w.(http.Flusher).Flush()

	for {
		select {
		case message := <-requestChannel:
			fmt.Fprintf(w, "[%d] - %s\n", message.ID, message.Text)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			delete(messageChannels, requestID)
			fmt.Printf("Client disconnected with ID = %s \n", requestID)
			return
		}
	}
}

var producer sarama.SyncProducer
var db *gorm.DB
var messageChannels = make(map[string]chan MessageItem)

func main() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	log.Println("Starting...")

	dbInstanse, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=root dbname=defaultdb sslmode=disable", DB_HOST, DB_PORT))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected...")
	dbInstanse.AutoMigrate(&MessageItem{})
	db = dbInstanse

	producer = initProducer()
	go runConsumer()

	router := mux.NewRouter()
	router.HandleFunc("/messages", createMessageHandler).Methods("POST")
	router.HandleFunc("/messages", getMessagesHandler).Methods("GET")
	http.Handle("/", router)

	log.Println("Server is running...")
	http.ListenAndServe(APP_PORT, nil)
}
