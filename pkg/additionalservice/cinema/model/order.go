package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Start   int                `json:"start"`
	Film    Film               `json:"film"`
	Seats   []Seat             `json:"seats"`
	Price   int64              `json:"price"`
	Created int64              `json:"created"`
}

type Film struct {
	Title    string   `json:"title"`
	Rating   float64  `json:"rating"`
	Cashback float64  `json:"cashback"`
	Genres   []string `json:"genres"`
}

type Seat struct {
	Row    int `json:"row"`
	Number int `json:"number"`
}
