package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func initProducer() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.ClientID = "kafka"
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{BROKER_URL}, config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return producer
}

func runConsumer() {
	config := sarama.NewConfig()
	config.ClientID = "kafka"
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{BROKER_URL}, config)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		partitionConsumer, err := consumer.ConsumePartition(TOPIC_NAME, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatal(err.Error())
		}

		for {
			select {
			case msg := <-partitionConsumer.Messages():
				newMessage := MessageItem{Text: string(msg.Value)}
				db.Create(&newMessage)
				for _, ch := range messageChannels {
					go func(ch chan MessageItem, value *MessageItem) {
						ch <- *value
					}(ch, &newMessage)
				}

			case err := <-partitionConsumer.Errors():
				fmt.Printf("Error: %s\n", err.Err)
			}
		}
	}

}

func sendToQueue(producer sarama.SyncProducer, text string) {
	message := &sarama.ProducerMessage{
		Topic: TOPIC_NAME,
		Value: sarama.StringEncoder(text),
	}
	_, _, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println(err.Error())
	}

}
