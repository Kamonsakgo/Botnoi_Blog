package entities

import "time"

type WrongMessage struct {
	Language string    `json:"language" bson:"language"`
	Message  string    `json:"message" bson:"message"`
	Speaker  string    `json:"speaker" bson:"speaker"`
	Datetime time.Time `json:"datetime" bson:"datetime"`
}
