package main

type Message struct {
	Text string `json:"text"`
}

type MessageItem struct {
	ID   uint64 `gorm:"primaryKey"`
	Text string
}
