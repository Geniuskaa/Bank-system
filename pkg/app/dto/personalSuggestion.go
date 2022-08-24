package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Suggestion struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      int64              `json:"user_id" bson:"user_id"`
	Description []string           `json:"description"`
}
