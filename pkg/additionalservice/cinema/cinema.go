package cinema

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Service struct {
	mongoDB *mongo.Database
}

func NewService(mongoDB *mongo.Database) *Service {
	return &Service{mongoDB: mongoDB}
}

func (s *Service) GetAllOrders(ctx context.Context) ([]byte, error) {
	cursor, err := s.mongoDB.Collection("orders").Find(ctx, bson.D{})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Print(cerr)
		}
	}()

	orders := make([]Order, 0)
	for cursor.Next(ctx) {
		var order Order
		err = cursor.Decode(&order)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = cursor.Err(); err != nil {
		log.Print(err)
		return nil, err
	}

	body, err := json.Marshal(orders)
	if err != nil {
		e := fmt.Errorf("Error when marshaling orders: %w", err)
		log.Println(e)
		return nil, e
	}

	return body, nil
}

func (s *Service) GetOrderByID(ctx context.Context, orderID string) ([]byte, error) {
	id, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, err
	}

	var order Order
	err = s.mongoDB.Collection("orders").FindOne(
		ctx,
		bson.D{{"_id", id}},
	).Decode(&order)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Service) SearchByMinRating(ctx context.Context, rating float64) ([]byte, error) {
	cursor, err := s.mongoDB.Collection("orders").Find(
		ctx,
		bson.D{
			{"film.rating", bson.D{
				{"$gt", rating},
			}},
		},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Print(cerr)
		}
	}()

	orders := make([]Order, 0)
	for cursor.Next(ctx) {
		var order Order
		err = cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	body, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Service) Save(ctx context.Context, order *Order) error {
	order.Created = time.Now().Unix()
	result, err := s.mongoDB.Collection("orders").InsertOne(ctx, order)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *Service) Update(ctx context.Context, order *Order) error {
	result, err := s.mongoDB.Collection("orders").UpdateOne(
		ctx,
		bson.D{{"_id", order.ID}},
		bson.D{
			{"$set", bson.D{
				{"start", order.Start},
				{"price	", order.Price},
			}},
		})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("Something going wrong and entity wasn`t updated", errors.New("Update failed"))
	}

	return nil
}
