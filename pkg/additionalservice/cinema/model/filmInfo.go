package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilmInfo struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title  string             `json:"title" `
	Rating float64            `json:"rating"`
	Genres []string           `json:"genres"`
	Start  int64              `json:"start"`
}
