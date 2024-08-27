package model

import "time"

type Error struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func NewError(message string) Error {
	return Error{
		Message:   message,
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
	}
}
